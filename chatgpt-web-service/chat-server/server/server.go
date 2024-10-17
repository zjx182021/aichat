package server

import (
	chatcontext "chatgpt-web-service/chat-server/chat-context"
	"chatgpt-web-service/chat-server/data"
	metricsbus "chatgpt-web-service/chat-server/metrics_bus"
	"chatgpt-web-service/chat-server/vector_data"
	"chatgpt-web-service/pkg/config"
	predis "chatgpt-web-service/pkg/db/redis"
	"chatgpt-web-service/pkg/log"
	"chatgpt-web-service/proto"
	"chatgpt-web-service/services/tokenizer"
	"context"
	"encoding/json"
	"io"
	"strconv"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
)

type ChatServer struct {
	proto.UnimplementedChatServer
	config     *config.Config
	log        log.ILogger
	data       data.IChatRecords
	vector     vector_data.IChatRecordsData
	busMetrics *metricsbus.BusMetrics
}

func NewChatService(db data.IChatRecords, config *config.Config, log log.ILogger, vector vector_data.IChatRecordsData) proto.ChatServer {
	return &ChatServer{
		data:   db,
		config: config,
		log:    log,
		vector: vector,
	}
}
func (c *ChatServer) ChatCompletion(ctx context.Context, p *proto.ChatCompletionRequest) (*proto.ChatCompletionResponse, error) {
	app := c.NewApp(p, chatcontext.Newrediscache())
	predis.InitRedis()
	keywords := app.Keywords(p)
	ok, msg, err := app.Sensitive(p)
	if err != nil {
		return nil, err
	}
	if !ok {
		res := app.buildChatCompletionResponse(msg)
		return res, nil
	}
	if len(keywords) > 0 {
		idstr, score, err := c.vector.QueryData(context.Background(), map[string][]string{"keywords": {strings.Join(keywords, ",")}})
		if err != nil {
			log.My_log.Error(err)
		} else if score > 0.99 {
			id, err := strconv.ParseUint(idstr, 10, 64)
			if err != nil {
				log.My_log.Error(err)
			} else {

				record := c.data.GetById(uint(id))
				if record.ID == 0 {
					log.My_log.Errorf("GetById查询无记录或者断言失败%s", err)
				} else {
					res := app.buildChatCompletionResponse(record.AIMsg)
					return res, nil
				}

			}
		}

	}
	client := app.GetOpenaiconn()
	req, tokens, currtokens, currmsg, err := app.buildChatCompletionRequest(p, false)
	if err != nil {
		log.My_log.Error(err)
		return nil, err
	}
	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		log.My_log.Error(err)
		return nil, err
	}
	res := &proto.ChatCompletionResponse{}
	bytes, err := json.Marshal(resp)
	if err != nil {
		log.My_log.Error(err)
		return nil, err
	}
	err = jsonpb.UnmarshalString(string(bytes), res)
	if err != nil {
		log.My_log.Error(err)
		return nil, err
	}
	go func() {
		reqctx := &chatcontext.Message{
			ID:      p.Id,
			PID:     p.Pid,
			Message: currmsg,
			Tokens:  currtokens,
		}
		err := app.SaveContext(reqctx)
		if err != nil {
			log.My_log.Error(err)
			return
		}
		resctx := &chatcontext.Message{
			ID:      resp.ID,
			PID:     reqctx.PID,
			Message: resp.Choices[0].Message,
			Tokens:  resp.Usage.CompletionTokens,
		}
		err = app.SaveContext(resctx)
		if err != nil {
			log.My_log.Error(err)
			return
		}
	}()

	go func() {
		keywordsJSON, err := json.Marshal(keywords)
		if err != nil {
			log.My_log.Error(err)
			return
		}
		records := &data.ChatRecords{
			UserMsg:         p.Message,
			UserMsgToken:    currtokens,
			UserMsgKeywords: keywordsJSON,
			AIMsg:           resp.Choices[0].Message.Content,
			AIMsgTokens:     resp.Usage.CompletionTokens,
			ReqTokens:       int64(tokens),
		}
		c.data.AddRecords(records)

		if len(keywords) > 0 {
			list := []*vector_data.ChatRecord{
				{
					ID: strconv.FormatInt(int64(records.ID), 10),
					KVs: map[string]string{
						"keywords": strings.Join(keywords, ","),
					},
				},
			}
			err := c.vector.UpsertData(context.Background(), list)
			if err != nil {
				log.My_log.Error(err)
				return
			}
		}
	}()
	return res, err
}
func (s *ChatServer) ChatCompletionStream(in *proto.ChatCompletionRequest, stream proto.Chat_ChatCompletionStreamServer) error {
	redisContextCache := chatcontext.Newrediscache()

	app := s.NewApp(in, redisContextCache)
	//敏感词过滤
	ok, msg, err := app.Sensitive(in)
	if err != nil {
		s.busMetrics.ErrQuestionsTotalCounter.Inc()
		s.log.Error(err)
		return err
	}
	if !ok {
		s.busMetrics.SensitiveQuestionsTotalCounter.Inc()
		resId := uuid.New().String()
		startRes := app.buildChatCompletionStreamResponse(resId, "", "")
		endRes := app.buildChatCompletionStreamResponse(resId, "", "stop")
		err = stream.Send(startRes)
		if err != nil {
			s.log.Error(err)
			return err
		}
		resList := app.buildChatCompletionStreamResponseList(resId, msg)
		for _, res := range resList {
			err = stream.Send(res)
			if err != nil {
				s.log.Error(err)
				return err
			}
		}
		err = stream.Send(endRes)
		if err != nil {
			s.log.Error(err)
			return err
		}
		return nil
	}

	//关键词提取
	keywords := app.Keywords(in)

	if len(keywords) > 0 {
		s.busMetrics.KeywordsQuestionsTotalCounter.Inc()
		idStr, score, err := s.vector.QueryData(context.Background(), map[string][]string{"keywords": {strings.Join(keywords, ",")}})
		if err != nil {
			s.log.Error(err)
		} else if score > 0.99 {
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				s.log.Error(err)
			} else {
				record := s.data.GetById(uint(id))

				resId := uuid.New().String()
				startRes := app.buildChatCompletionStreamResponse(resId, "", "")
				endRes := app.buildChatCompletionStreamResponse(resId, "", "stop")
				err = stream.Send(startRes)
				if err != nil {
					s.log.Error(err)
					return err
				}
				resList := app.buildChatCompletionStreamResponseList(resId, record.AIMsg)
				for _, res := range resList {
					err = stream.Send(res)
					if err != nil {
						s.log.Error(err)
						return err
					}
				}
				err = stream.Send(endRes)
				if err != nil {
					s.log.Error(err)
					return err
				}
				return nil

			}
		}
	}

	client := app.GetOpenaiconn()
	req, tokens, currTokens, currMessage, err := app.buildChatCompletionRequest(in, false)
	if err != nil {
		s.log.Error(err)
		return err
	}
	chatStream, err := client.CreateChatCompletionStream(stream.Context(), req)
	if err != nil {
		s.busMetrics.ErrQuestionsTotalCounter.Inc()
		s.log.Error(err)
		return err
	}
	defer chatStream.Close()
	completionContent := ""
	resultID := ""
	for {
		resp, err := chatStream.Recv()
		if err != nil && err != io.EOF {
			s.log.Error(err)
			return err
		}
		if err == io.EOF {
			break
		}
		if resultID == "" {
			resultID = resp.ID
		}
		completionContent += resp.Choices[0].Delta.Content
		res := &proto.ChatCompletionStreamResponse{}
		bytes, err := json.Marshal(resp)
		if err != nil {
			s.log.Error(err)
			return err
		}
		err = jsonpb.UnmarshalString(string(bytes), res)
		if err != nil {
			s.log.Error(err)
			return err
		}
		err = stream.Send(res)
		if err != nil {
			s.log.Error(err)
			return err
		}
	}
	resultMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: completionContent,
	}
	model := s.config.Chat.Model
	if in.ChatParam != nil && in.ChatParam.Model != "" {
		model = in.ChatParam.Model
	}
	resultTokens, err := tokenizer.Gettokens(&resultMessage, model)
	if err != nil {
		s.busMetrics.ErrQuestionsTotalCounter.Inc()
		s.log.Error(err)
		return err
	}

	go func() {
		reqContext := &chatcontext.Message{
			ID:      in.Id,
			PID:     in.Pid,
			Message: currMessage,
			Tokens:  currTokens,
		}
		err := app.SaveContext(reqContext)
		if err != nil {
			s.log.Error(err)
			return
		}
		resContext := &chatcontext.Message{
			ID:      resultID,
			PID:     reqContext.ID,
			Message: resultMessage,
			Tokens:  resultTokens,
		}
		err = app.SaveContext(resContext)
		if err != nil {
			s.log.Error(err)
			return
		}
	}()
	go func() {
		s.busMetrics.QuestionsTotalCounter.Inc()
		keywordsJSON, err := json.Marshal(keywords)
		if err != nil {
			s.log.Error(err)
			return
		}
		records := &data.ChatRecords{
			UserMsg:         in.Message,
			UserMsgToken:    currTokens,
			UserMsgKeywords: keywordsJSON,
			AIMsg:           completionContent,
			AIMsgTokens:     resultTokens,
			ReqTokens:       int64(tokens),
		}
		s.data.AddRecords(records)
		//保存到向量数据库
		if len(keywords) > 0 {
			list := []*vector_data.ChatRecord{
				{
					ID: strconv.FormatInt(int64(records.ID), 10),
					KVs: map[string]string{
						"keywords": strings.Join(keywords, ","),
					},
				},
			}
			err = s.vector.UpsertData(context.Background(), list)
			if err != nil {
				s.log.Error(err)
				return
			}
		}
	}()
	return nil
}

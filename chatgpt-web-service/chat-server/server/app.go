package server

import (
	chatcontext "chatgpt-web-service/chat-server/chat-context"
	"chatgpt-web-service/pkg/config"
	"chatgpt-web-service/pkg/log"
	"chatgpt-web-service/proto"
	"chatgpt-web-service/services"
	keywordsfilter "chatgpt-web-service/services/key_words_filter"
	words_proto "chatgpt-web-service/services/key_words_filter/proto"
	"chatgpt-web-service/services/tokenizer"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
)

const ChatPrimedTokens = 2

// chat:
//   api_key: "i0jey84SdkFdw5u43780yjr3h7se8nth0yi295nr94ksDngKprEh"
//   base_url: "http://192.168.88.135:8084/v1"
//   model: "gpt-3.5-turbo"
//   max_tokens: 4096
//   temperature: 0.8
//   top_p: 0.9
//   frequency_penalty: 0.8
//   presence_penalty: 0.5
//   bit_desc: "你是一个ai助手,我需要你模拟一名软件工程师来回答我的问题"
//   min_response_tokens: 2048
//   context_ttl: 1800
//   context_len: 4

type Openaiconf struct {
	Api_key             string  `mapstructure:"api_key"`
	Base_url            string  `mapstructure:"base_url"`
	Model               string  `mapstructure:"model"`
	Max_tokens          int     `mapstructure:"max_tokens"`
	Temperature         float32 `mapstructure:"temperature"`
	Top_p               float32 `mapstructure:"top_p"`
	Frequency_penalty   float32 `mapstructure:"frequency_penalty"`
	Presence_penalty    float32 `mapstructure:"presence_penalty"`
	Min_response_tokens int     `mapstructure:"min_response_tokens"`
	Context_ttl         int     `mapstructure:"context_ttl"`
	Context_len         int     `mapstructure:"context_len"`
	BotDescription      string  `mapstructure:"bot_desc"`
}

type App struct {
	openconf *Openaiconf
	log      log.ILogger
	cache    chatcontext.ContextCache
}

func (s *ChatServer) NewApp(in *proto.ChatCompletionRequest, cache chatcontext.ContextCache) *App {
	conf := &Openaiconf{
		Api_key:             s.config.Chat.APIKey,
		Base_url:            s.config.Chat.BaseURL,
		Model:               s.config.Chat.Model,
		Max_tokens:          s.config.Chat.MinResponseTokens,
		Temperature:         s.config.Chat.Temperature,
		Top_p:               s.config.Chat.TopP,
		Frequency_penalty:   s.config.Chat.FrequencyPenalty,
		Presence_penalty:    s.config.Chat.PresencePenalty,
		Min_response_tokens: s.config.Chat.MinResponseTokens,
		Context_ttl:         s.config.Chat.ContextTTL,
		Context_len:         s.config.Chat.ContextLen,
	}
	if in.ChatParam != nil {
		if in.ChatParam.Model != "" {
			conf.Model = in.ChatParam.Model
		}
		if in.ChatParam.TopP != 0 {
			conf.Top_p = in.ChatParam.TopP
		}
		if in.ChatParam.FrequencyPenalty != 0 {
			conf.Frequency_penalty = in.ChatParam.FrequencyPenalty
		}
		if in.ChatParam.PresencePenalty != 0 {
			conf.Presence_penalty = in.ChatParam.PresencePenalty
		}
		if in.ChatParam.Temperature != 0 {
			conf.Temperature = in.ChatParam.Temperature
		}
		if in.ChatParam.BotDesc != "" {
			conf.BotDescription = in.ChatParam.BotDesc
		}
		if in.ChatParam.MaxTokens != 0 {
			conf.Max_tokens = int(in.ChatParam.MaxTokens)
		}
		if in.ChatParam.ContextTTL != 0 {
			conf.Context_ttl = int(in.ChatParam.ContextTTL)
		}
		if in.ChatParam.ContextLen != 0 {
			conf.Context_len = int(in.ChatParam.ContextLen)
		}
		if in.ChatParam.MinResponseTokens != 0 {
			conf.Min_response_tokens = int(in.ChatParam.MinResponseTokens)
		}
	}
	return &App{
		openconf: conf,
		log:      s.log,
		cache:    cache,
	}
}

func (a *App) Keywords(in *proto.ChatCompletionRequest) []string {
	pool := keywordsfilter.GetKeywordsClientPool()
	conn := pool.Get()
	defer pool.Put(conn)
	accesstoken := config.GetConfig().DependOn.Keywords.AccessToken
	client := words_proto.NewFilterClient(conn)
	ctx := services.PusBearToToken(context.Background(), accesstoken)
	req := &words_proto.FilterReq{
		Text: in.Message,
	}
	res, err := client.FindAll(ctx, req)
	if err != nil {
		a.log.Error(err)
		return []string{}
	}
	return res.GetKeywords()
}

func (a *App) Sensitive(in *proto.ChatCompletionRequest) (bool, string, error) {
	// 获得连接
	keywordspool := keywordsfilter.GetKeywordsClientPool()
	conn := keywordspool.Get()
	defer keywordspool.Put(conn)
	client := words_proto.NewFilterClient(conn)
	// 获得token
	localtoken := config.GetConfig().DependOn.Sensitive.AccessToken
	ctx := services.PusBearToToken(context.Background(), localtoken)
	// findall
	req := &words_proto.FilterReq{
		Text: in.Message,
	}
	res, err := client.Validate(ctx, req)
	if err != nil {
		a.log.Error(err)
		return false, "", err
	}
	ok := res.Ok
	if !ok {
		return false, "触发知识盲区，换个问题再问", nil
	}
	return true, res.GetKeyword(), nil
}

func (a *App) GetOpenaiconn() *openai.Client {
	accesstoken := a.openconf.Api_key
	clientConfig := openai.DefaultConfig(accesstoken)
	clientConfig.BaseURL = a.openconf.Base_url
	client := openai.NewClientWithConfig(clientConfig)
	return client
}

func (a *App) buildChatCompletionResponse(msg string) *proto.ChatCompletionResponse {
	return &proto.ChatCompletionResponse{
		Id:      uuid.New().String(),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   a.openconf.Model,
		Choices: []*proto.ChatCompletionChoice{
			{
				Message: &proto.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleAssistant,
					Content: msg,
				},
				FinishReason: "stop",
			},
		},
		Usage: &proto.Usage{
			PromptTokens:     0,
			CompletionTokens: 0,
			TotalTokens:      0,
		},
		// Model   string                  `protobuf:"bytes,4,opt,name=model,proto3" json:"model,omitempty"`
		// Choices []*ChatCompletionChoice `protobuf:"bytes,5,rep,name=choices,proto3" json:"choices,omitempty"`
		// Usage   *Usage                  `protobuf:"bytes,6,opt,name=usage,proto3" json:"usage,omitempty"`
	}
}

func (a *App) buildChatCompletionRequest(in *proto.ChatCompletionRequest, stream bool) (req openai.ChatCompletionRequest,
	tokens, currToken int, currMessage openai.ChatCompletionMessage, err error) {
	currMessage = openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: in.Message,
	}
	req = openai.ChatCompletionRequest{
		Model: a.openconf.Model,
		Messages: []openai.ChatCompletionMessage{
			currMessage,
		},
		MaxTokens:        a.openconf.Max_tokens,
		Temperature:      a.openconf.Temperature,
		TopP:             a.openconf.Top_p,
		FrequencyPenalty: a.openconf.Frequency_penalty,
		PresencePenalty:  a.openconf.Presence_penalty,
		Stream:           stream,
	}
	contextList := make([]*chatcontext.Message, 0)
	if in.EnableContext {
		contextList = a.GetContext(in.Pid)
	}
	tokens, currToken, req.Messages, err = a.RebuildMessage(contextList, currMessage)
	if err != nil {
		log.My_log.Errorf("RebuildMessage failed %s", err)
		return
	}
	req.MaxTokens = a.openconf.Max_tokens - tokens
	return
}

// func (a *App)buildChatCompletionResponse(msg string) openai.ChatCompletionResponse {
// 	res := openai.ChatCompletionResponse{
// 		ID:      uuid.New().String(),
// 		Object:  "chat.completion",
// 		Created: time.Now().Unix(),
// 		Model:   "gpt-3.5-turbo-0301",
// 		Choices: []openai.ChatCompletionChoice{
// 			{
// 				Message: openai.ChatCompletionMessage{
// 					Role:    openai.ChatMessageRoleAssistant,
// 					Content: msg,
// 				},
// 				FinishReason: "stop",
// 			},
// 		},
// 		Usage: openai.Usage{
// 			PromptTokens:     10,
// 			CompletionTokens: 20,
// 			TotalTokens:      30,
// 		},
// 	}
// 	return res
// }

func (a *App) buildChatCompletionStreamResponseList(id, msg string) []*proto.ChatCompletionStreamResponse {
	list := make([]*proto.ChatCompletionStreamResponse, 0)
	for _, delta := range msg {
		list = append(list, a.buildChatCompletionStreamResponse(id, string(delta), ""))
	}

	return list
}

func (a *App) buildChatCompletionStreamResponse(id, delta, finishReason string) *proto.ChatCompletionStreamResponse {
	res := &proto.ChatCompletionStreamResponse{
		Id:      id,
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   "gpt-3.5-turbo-0301",
		Choices: []*proto.ChatCompletionStreamChoice{
			{
				Index: 0,
				Delta: &proto.ChatCompletionStreamChoiceDelta{
					Content: delta,
					Role:    openai.ChatMessageRoleAssistant,
				},
				FinishReason: finishReason,
			},
		},
	}
	return res
}

func (a *App) GetContext(id string) []*chatcontext.Message {
	maxlen := a.openconf.Context_len
	list := make([]*chatcontext.Message, 0, maxlen)
	for i := 0; i < maxlen; i++ {
		msg, err := a.cache.Get(id)
		if err != nil || msg == nil {
			log.My_log.Errorf("openai 获得上下文失败%s", err)
		}
		list = append(list, msg)
		id = msg.PID
	}
	return list
}

func (a *App) SaveContext(value *chatcontext.Message) error {
	err := a.cache.Set(value.ID, *value, a.openconf.Context_ttl)
	if err != nil {
		log.My_log.Errorf("openai 存储上下文失败%s", err)
		return err
	}
	return nil
}

func (a *App) RebuildMessage(contextlist []*chatcontext.Message, currmessage openai.ChatCompletionMessage) (tokens, currTokens int, messages []openai.ChatCompletionMessage, err error) {
	var sysmsg openai.ChatCompletionMessage
	var tokennum = 0
	var currtokennum = 0
	if a.openconf.BotDescription != "" {
		sysmsg = openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: a.openconf.BotDescription,
		}
		tokennum, err = tokenizer.Gettokens(&sysmsg, a.openconf.Model)
		if err != nil {
			return 0, 0, nil, err
		}
	}
	messages = []openai.ChatCompletionMessage{currmessage}
	currtokennum, err = tokenizer.Gettokens(&currmessage, a.openconf.Model)
	if err != nil {
		return 0, 0, nil, err
	}
	if currtokennum+tokennum+a.openconf.Min_response_tokens+ChatPrimedTokens > a.openconf.Max_tokens {
		log.My_log.Error("token 过长")
		return 0, 0, nil, errors.New("token 过长")
	}
	tokens = currtokennum + tokennum + ChatPrimedTokens
	for _, ctx := range contextlist {
		if tokens+ctx.Tokens+ChatPrimedTokens < a.openconf.Max_tokens {
			tokens += ctx.Tokens + ChatPrimedTokens
			messages = append(messages, ctx.Message)
		} else {
			break
		}
	}
	for i := 0; i < len(messages)/2; i++ {
		messages[i], messages[len(messages)-i-1] = messages[len(messages)-i-1], messages[i]
	}
	if tokennum > 0 {
		messages = append([]openai.ChatCompletionMessage{sysmsg}, messages...)
	}
	return tokens, currTokens, messages, nil
}

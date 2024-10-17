package tokenizer

import (
	"bytes"
	"chatgpt-web-service/pkg/config"
	"chatgpt-web-service/pkg/log"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sashabaranov/go-openai"
)

type tokensinfo struct {
	Code   int    `json:"code"`
	Tokens int    `json:"num_tokens"`
	Msg    string `json:"msg"`
}

func Gettokens(msg *openai.ChatCompletionMessage, model string) (int, error) {
	cnf := config.GetConfig()
	url := cnf.DependOn.Tokenizer.Address
	url = fmt.Sprintf("%s/tokenizer/%s", url, model)
	msgdata, err := json.Marshal(msg)
	if err != nil {
		log.My_log.ErrorF("msg marshal失败%s", err)
		return 0, err
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(msgdata))
	if err != nil {
		log.My_log.ErrorF("post request 到tokenizer 失败%s", err)
		return 0, err
	}
	info := &tokensinfo{}
	err = json.NewDecoder(resp.Body).Decode(info)
	if err != nil {
		log.My_log.ErrorF("json decode 失败%s", err)
		return 0, err
	}
	if info.Code != 200 {
		log.My_log.ErrorF("tokenizer 出错, code:%d, msg:%s", info.Code, info.Msg)
		return 0, fmt.Errorf("tokenizer 出错, code:%d, msg:%s", info.Code, info.Msg)
	}
	return info.Tokens, nil
}

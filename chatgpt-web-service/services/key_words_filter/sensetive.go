package keywordsfilter

import (
	"chatgpt-web-service/pkg/config"
	grpcclient "chatgpt-web-service/services/grpc-client"
	"sync"
)

var sensitivepool grpcclient.ClientPool
var sensitiveonce sync.Once

type Sensitive struct {
	grpcclient.DefaultClient
}

func GetsensitivePool() *grpcclient.ClientPool {
	sensitiveonce.Do(func() {
		cnf := config.GetConfig()
		c := Keyword{}
		sensitivepool = c.GetPool(cnf.DependOn.Sensitive.Address)
	})
	return &sensitivepool
}

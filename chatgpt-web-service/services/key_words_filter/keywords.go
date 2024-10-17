package keywordsfilter

import (
	"chatgpt-web-service/pkg/config"
	grpcclient "chatgpt-web-service/services/grpc-client"
	"sync"
)

var keywordspool grpcclient.ClientPool
var keywordsonce sync.Once

type Keyword struct {
	grpcclient.DefaultClient
}

func GetKeywordsClientPool() grpcclient.ClientPool {
	keywordsonce.Do(func() {
		cnf := config.GetConfig()
		c := Keyword{}
		keywordspool = c.GetPool(cnf.DependOn.Keywords.Address)
	})
	return keywordspool
}

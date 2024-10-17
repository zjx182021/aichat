package vector

import (
	"chatgpt-web-service/pkg/config"
	"chatgpt-web-service/pkg/log"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb"
	"time"
)

var vdb *tcvectordb.Client

func InitDB(config *config.Config) {
	var defaultOption = &tcvectordb.ClientOption{
		Timeout:            time.Second * time.Duration(config.VectorDB.Timeout),
		MaxIdldConnPerHost: config.VectorDB.MaxIdleConnPerHost,
		IdleConnTimeout:    time.Second * time.Duration(config.VectorDB.IdleConnTimeout),
		ReadConsistency:    tcvectordb.ReadConsistency(config.VectorDB.ReadConsistency),
	}
	var err error
	vdb, err = tcvectordb.NewClient(config.VectorDB.Url, config.VectorDB.Username, config.VectorDB.Pwd, defaultOption)
	if err != nil {
		log.My_log.Error(err)
		return
	}
}

func GetVdb() *tcvectordb.Client {
	return vdb
}

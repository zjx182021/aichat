package mysql

import (
	"chatgpt-web-service/pkg/config"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var MYSQL *gorm.DB
var mutex sync.Mutex

// host: "192.168.88.135"
// port: 3306
// username: "root"
// password: "!Zhang123456"
// dbname: "aichat"
// table: "chatget_web_service"
// max_open_conns: 100
// max_idle_conns: 50
// max_life_time: 1800
// max_conn_lifetime: 1800
// idle_timeout: 1800
func Initmysql() {
	if MYSQL == nil {
		mutex.Lock()
		defer mutex.Unlock()
		if MYSQL == nil {
			mylog := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.LogLevel(4),
				Colorful:      true,
			})
			//   dns: "root:!Zhang123456@tcp(127.0.0.1:3306)/ginchat?charset=utf8mb4&parseTime=True&loc=Local"
			mysqldns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				config.GetConfig().Mysql.Username,
				config.GetConfig().Mysql.Password,
				config.GetConfig().Mysql.Host,
				config.GetConfig().Mysql.Port,
				config.GetConfig().Mysql.DBname,
			)

			mysql, err := gorm.Open(mysql.Open(mysqldns), &gorm.Config{
				Logger: mylog,
			})

			if err != nil {
				log.Fatalf("failed to connect database: %v", err)
			}
			MYSQL = mysql

		}
	}
}

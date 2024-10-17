package data

import (
	"chatgpt-web-service/pkg/db/mysql"
	"encoding/json"

	"gorm.io/gorm"
)

type ChatRecords struct {
	gorm.Model
	UserMsg         string
	UserMsgToken    int
	UserMsgKeywords json.RawMessage `gorm:"serializer:json"`
	AIMsg           string
	AIMsgTokens     int
	ReqTokens       int64
}

func Inittable() {
	mysql.MYSQL.AutoMigrate(&ChatRecords{})
}

type IChatRecords interface {
	AddRecords(c *ChatRecords)
	GetById(id uint) *ChatRecords
}

func Newrecords() IChatRecords {
	return &ChatRecords{}
}
func (c *ChatRecords) AddRecords(newc *ChatRecords) {
	mysql.MYSQL.Create(newc)
}

func (c *ChatRecords) GetById(id uint) *ChatRecords {
	var record ChatRecords
	mysql.MYSQL.Where("id =?", id).First(&record)
	return &record
}

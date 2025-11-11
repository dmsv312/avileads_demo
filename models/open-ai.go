package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

const (
	ROLE_USER      = "user"
	ROLE_ASSISTANT = "assistant"
)

type OpenAIMessages struct {
	Id          int       `orm:"column(id);auto"`
	SessionID   string    `orm:"column(session_id)"`
	Role        string    `orm:"column(role)"`
	Message     string    `orm:"column(message)"`
	UserID      int       `orm:"column(user_id)"`
	Temperature float32   `orm:"column(temperature)"`
	CreateDate  time.Time `orm:"column(create_date)"`
}

type OpenAIRequest struct {
	SessionID   string  `json:"session_id"`
	Message     string  `json:"message"`
	Temperature float32 `json:"temperature"`
}

func (m *OpenAIMessages) TableName() string {
	return "openai_messages"
}

func SaveOpenAIMessage(request OpenAIRequest, message string, userId int, role string, o orm.Ormer) (bool, error) {
	newMessage := OpenAIMessages{
		SessionID:   request.SessionID,
		Role:        role,
		Message:     message,
		UserID:      userId,
		Temperature: request.Temperature,
		CreateDate:  time.Now().UTC(),
	}

	_, err := o.Insert(&newMessage)

	if err != nil {
		return false, err
	}

	return true, nil
}

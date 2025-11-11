package models

import "time"

type AuditLog struct {
	Id        int       `orm:"auto"`
	UserId    int       `orm:"column(user_id)"`
	EntityId  int       `orm:"column(entity_id)"`
	Changes   string    `orm:"type(text)"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
	ModelId   int       `orm:"column(model_id)"`
	ActionId  int       `orm:"column(action_id)"`
	RkId      int       `orm:"column(rk_id)"`
	//
	//AuditModel  *AuditModel  `orm:"rel(fk)"`
	//AuditAction *AuditAction `orm:"rel(fk)"`
}

type AuditAction struct {
	Id   int    `orm:"pk"`
	Name string `orm:"size(16)"`
}

type AuditModel struct {
	Id   int    `orm:"pk;auto"`
	Name string `orm:"size(128);unique"`
}

func (AuditLog) TableName() string { return "audit_logs" }

func (AuditAction) TableName() string { return "audit_actions" }

func (AuditModel) TableName() string { return "audit_models" }

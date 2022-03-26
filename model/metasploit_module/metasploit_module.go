package metasploit_module

import (
	"github-cve/db"
	"time"
)

type MetasploitModule struct {
	ID              uint       `json:"id" gorm:"type:int not null auto_increment comment '记录ID'"`
	Content         string     `json:"content" gorm:"type:text not null comment '内容'"`
	VulnerabilityID uint       `json:"vulnerability_id" gorm:"type:int not null comment '漏洞ID'"`
	CreatedAt       *time.Time `json:"created_at" gorm:"type:timestamp not null default current_timestamp"`
	UpdatedAt       *time.Time `json:"updated_at" gorm:"type:timestamp null on update current_timestamp"`
}

func (receiver *MetasploitModule) TableName() string {
	return "metasploit_module"
}

func CreateTable() error {
	return db.DB.Set("gorm:table_options", "ENGINE=InnoDB comment 'metasploit模块'").AutoMigrate(new(MetasploitModule))
}

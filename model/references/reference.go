package references

import (
	"github-cve/db"
	"time"
)

type Reference struct {
	ID              uint       `json:"id" gorm:"type:int not null auto_increment comment '记录ID'"`
	Link            string     `json:"link" gorm:"type:text null comment '链接'"`
	Content         string     `json:"content" gorm:"type:text null comment '链接内容'"`
	VulnerabilityID uint       `json:"vulnerability_id" gorm:"type:int not null comment '漏洞ID'"`
	CreatedAt       *time.Time `json:"created_at" gorm:"type:timestamp not null default current_timestamp"`
	UpdatedAt       *time.Time `json:"updated_at" gorm:"type:timestamp null on update current_timestamp"`
}

func (receiver *Reference) TableName() string {
	return "reference"
}

func CreateTable() error {
	return db.DB.Set("gorm:table_options", "ENGINE=InnoDB comment '参考链接'").AutoMigrate(new(Reference))
}

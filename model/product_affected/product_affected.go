package product_affected

import (
	"github-cve/db"
	"time"
)

type ProductAffected struct {
	ID              uint       `json:"id" gorm:"int not null auto_increment comment '记录ID'"`
	ProductType     string     `json:"product_type" gorm:"type:text null comment '产品类型'"`
	Vendor          string     `json:"vendor" gorm:"type:text null comment '产品销售公司'"`
	Product         string     `json:"product" gorm:"type:text null comment '产品名称'"`
	Version         string     `json:"version" gorm:"type:varchar(10) null comment '版本号'"`
	Update          string     `json:"update" gorm:"type:text null comment '更新时间'"`
	Edition         string     `json:"edition" gorm:"type:text null comment '发行版本'"`
	Language        string     `json:"language" gorm:"type:text null comment '语言'"`
	VersionDetail   string     `json:"version_detail" gorm:"type:text null comment '版本细节'"`
	Vulnerabilities string     `json:"vulnerabilities" gorm:"type:text null comment '产品漏洞'"`
	VulnerabilityID uint       `json:"vulnerability_id" gorm:"type:int not null comment '漏洞ID'"`
	CreatedAt       *time.Time `json:"created_at" gorm:"type:timestamp not null default current_timestamp"`
	UpdatedAt       *time.Time `json:"updated_at" gorm:"type:timestamp null on update current_timestamp"`
}

func (receiver ProductAffected) TableName() string {
	return "product_affected"
}

func CreateTable() error {
	return db.DB.Set("gorm:table_options", "ENGINE=InnoDB comment '受影响产品'").AutoMigrate(new(ProductAffected))
}

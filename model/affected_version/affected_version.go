package affected_version

import (
	"github-cve/db"
	"time"
)

type AffectedVersion struct {
	ID                 uint       `json:"id" gorm:"type:int;not null;auto_increment;comment:记录ID"`
	Vendor             string     `json:"vendor" gorm:"type:text;comment:供应商"`
	Product            string     `json:"product" gorm:"type:text;comment:产品"`
	VulnerableVersions string     `json:"vulnerable_versions" gorm:"type:text;comment:易受攻击的版本"`
	VulnerabilityID    uint       `json:"vulnerability_id" gorm:"type:int not null comment '漏洞ID'"`
	CreatedAt          *time.Time `json:"created_at" gorm:"type:timestamp not null default current_timestamp"`
	UpdatedAt          *time.Time `json:"updated_at" gorm:"type:timestamp null on update current_timestamp"`
}

func (receiver *AffectedVersion) TableName() string {
	return "affected_version"
}
func CreateTable() error {
	return db.DB.Set("gorm:table_options", "ENGINE=InnoDB comment '受影响版本'").AutoMigrate(new(AffectedVersion))
}

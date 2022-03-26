package db_affected_version

import (
	"fmt"
	"github-cve/db"
	"github-cve/model/affected_version"
)

func CreateItem(affectedVersion *affected_version.AffectedVersion) error {
	fmt.Println(affectedVersion)
	return db.DB.Create(affectedVersion).Error
}

func DeleteItemByVulnerabilityId(vulnerabilityId uint) error {
	fmt.Println(vulnerabilityId)
	return db.DB.Where("vulnerability_id = ?", vulnerabilityId).Delete(new(affected_version.AffectedVersion)).Error
}

func FirstOrCreate(affectedVersion *affected_version.AffectedVersion) error {
	fmt.Println(affectedVersion)
	return db.DB.FirstOrCreate(affectedVersion, affectedVersion).Error
}

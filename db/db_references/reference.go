package db_references

import (
	"fmt"
	"github-cve/db"
	"github-cve/model/references"
)

func CreateItem(reference *references.Reference) error {
	fmt.Println(reference)
	return db.DB.Create(reference).Error
}

func DeleteItemByVulnerabilityId(vulnerabilityId uint) error {
	fmt.Println(vulnerabilityId)
	return db.DB.Where("vulnerability_id = ?", vulnerabilityId).Delete(new(references.Reference)).Error
}

func FirstOrCreate(reference *references.Reference) error {
	fmt.Println(reference)
	return db.DB.FirstOrCreate(reference, reference).Error
}

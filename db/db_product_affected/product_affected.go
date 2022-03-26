package db_product_affected

import (
	"fmt"
	"github-cve/db"
	"github-cve/model/product_affected"
)

func CreateItem(productAffected *product_affected.ProductAffected) error {
	fmt.Println(productAffected)
	return db.DB.Create(productAffected).Error
}

func DeleteItemByVulnerabilityId(vulnerabilityId uint) error {
	fmt.Println(vulnerabilityId)
	return db.DB.Where("vulnerability_id = ?", vulnerabilityId).Delete(new(product_affected.ProductAffected)).Error
}

func FirstOrCreate(productAffected *product_affected.ProductAffected) error {
	fmt.Println(productAffected)
	return db.DB.FirstOrCreate(productAffected, productAffected).Error
}

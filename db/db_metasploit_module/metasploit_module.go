package db_metasploit_module

import (
	"fmt"
	"github-cve/db"
	"github-cve/model/metasploit_module"
)

func CreateItem(metasploitModule *metasploit_module.MetasploitModule) error {
	fmt.Println(metasploitModule)
	return db.DB.Create(metasploitModule).Error
}

func DeleteItemByVulnerabilityId(vulnerabilityId uint) error {
	fmt.Println(vulnerabilityId)
	return db.DB.Where("vulnerability_id = ?", vulnerabilityId).Delete(new(metasploit_module.MetasploitModule)).Error
}

func FirstOrCreate(metasploitModule *metasploit_module.MetasploitModule) error {
	fmt.Println(metasploitModule)
	return db.DB.FirstOrCreate(metasploitModule, metasploitModule).Error
}

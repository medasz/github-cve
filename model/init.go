package model

import (
	"log"

	"github-cve/model/affected_version"
	"github-cve/model/exploit"
	"github-cve/model/metasploit_module"
	"github-cve/model/product_affected"
	"github-cve/model/references"
	"github-cve/model/vulnerability"
)

func init() {
	// 创建数据库表
	var err error
	if err = affected_version.CreateTable(); err != nil {
		panic(err)
	}
	if err = metasploit_module.CreateTable(); err != nil {
		panic(err)
	}
	if err = product_affected.CreateTable(); err != nil {
		panic(err)
	}
	if err = references.CreateTable(); err != nil {
		panic(err)
	}
	if err = vulnerability.CreateTable(); err != nil {
		panic(err)
	}
	if err = exploit.CreateTable(); err != nil {
		panic(err)
	}
	log.Println("表创建成功")
}

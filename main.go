package main

import (
	"github-cve/crawlers"
	_ "github-cve/db"
	_ "github-cve/model"
)

func init() {
}

func main() {
	github := crawlers.Github{}
	go github.Run()
	exploitDB := crawlers.ExploitDB{}
	exploitDB.Run()
}

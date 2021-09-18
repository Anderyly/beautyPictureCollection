package main

import (
	"github.com/robfig/cron/v3"
	"log"
	"warm/task/niutu114"
	"warm/task/xiannvtu"
	"warm/task/yuacg"
)

func main() {
	spec := "30 1 * * *"
	c := cron.New()
	c.AddFunc(spec, run)
	c.Start()
	select {}
}

func run() {
	log.Println("已开启线程采集，author:anderyly")

	// 仙女图采集
	go xiannvtu.Start()
	// 牛图采集
	go niutu114.Start()
	// 雨溪萌域采集
	go yuacg.Start()

	select {}
}

package main

import (
	"log"
	"warm/task/niutu114"
	"warm/task/xiannvtu"
	"warm/task/yuacg"
)


func main() {

	start()
	return

	//spec := "45 10 * * *"
	//c := cron.New()
	//c.AddFunc(spec, start)
	//c.Start()
	//select {}
}

func start() {
	log.Println("已开启线程采集，author:anderyly")

	// 仙女图采集
	go xiannvtu.Start()
	// 牛图采集
	go niutu114.Start()
	// 雨溪萌域采集
	go yuacg.Start()

	select {

	}
}

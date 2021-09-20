package main

import (
	"github.com/robfig/cron/v3"
	"warm/task/mm29"
	"warm/task/niutu114"
	"warm/task/shejumm"
	"warm/task/xiannvtu"
	"warm/task/yuacg"
)

func main() {
	spec := "30 10 * * *"
	c := cron.New()
	c.AddFunc(spec, run)
	c.Start()
	select {}
}

func run() {

	go xiannvtu.Start() // 仙女图采集
	go niutu114.Start() // 牛图采集
	go yuacg.Start()    // 雨溪萌域采集
	go shejumm.Start()  // 射菊mm采集
	go mm29.Start()     // mm29采集

	select {}
}

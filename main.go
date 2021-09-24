package main

import (
	"github.com/robfig/cron/v3"
	"log"
	"strings"
	"warm/other"
	"warm/task/mm29"
	"warm/task/niutu114"
	"warm/task/shejumm"
	"warm/task/xiannvtu"
	"warm/task/yuacg"
)

func main() {
	log.Println("开始运行")

	isTask := other.GetConf("control.IsTask")

	if isTask == "1" {
		run()
	} else {
		log.Println("开启定时任务")
		spec := strings.Replace(other.GetConf("control.TaskTime"), "/", " ", -1)
		c := cron.New()
		c.AddFunc(spec, run)
		c.Start()
	}
	select {}

}

func run() {
	go xiannvtu.Start() // 仙女图采集
	go niutu114.Start() // 牛图采集
	go yuacg.Start()    // 雨溪萌域采集
	go shejumm.Start()  // 射菊mm采集
	go mm29.Start()     // mm29采集
}

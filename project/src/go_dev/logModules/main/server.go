package main

import(
"go_dev/oldBoy/logModules/tailf"
"go_dev/oldBoy/logModules/kafka"
"github.com/astaxie/beego/logs"
"time"
"fmt"
)
func serverRun()(err error){
	for {
		msg :=tailf.GetOneLine()
		err = mainSendToKafka(msg)
		if err !=nil {
			logs.Error("send to kafka failed, err:%v", err)
			time.Sleep(time.Second)
			continue
		}
	}	
	//return
}
func mainSendToKafka(msg *tailf.TextMsg)(err error){
	//logs.Debug("read msg:%s, topic:%s", msg.Msg, msg.Topic)
	fmt.Printf("read msg:%s, topic:%s\n", msg.Msg, msg.Topic)

	err = kafka.SendToKafka(msg.Msg, msg.Topic)
/* 	if err !=nil {
		logs.Error("send msg to kafka failed, err:%v", err)
		return
	} */
	return
}
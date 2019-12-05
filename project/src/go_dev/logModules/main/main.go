package main

import(
	"fmt"
	"github.com/astaxie/beego/logs"
	"go_dev/oldBoy/logModules/etcd"
	"go_dev/oldBoy/logModules/tailf"
	"go_dev/oldBoy/logModules/kafka"
	//"time"
)
func main(){
	conftype := "ini"
	filename :="./conf/logagent.conf"

	err :=loadConf(conftype, filename)
	if err !=nil {
		fmt.Printf("load conf failed, err:%v\n", err)
		panic("load conf failed")
	}
	logs.Debug("load conf succ, config:%v", appConfig)

	err = initLogger()
	if err !=nil {
		fmt.Printf("load logger failed, err:%v\n",err)
		panic("load logger failed")
	}
	//logs.Debug("initialize logger succ")
	confArray, err := etcd.InitGetEtcd(appConfig.etcdAddr, appConfig.etcdKey)
	if err !=nil {
		logs.Error("init etcd failed, err:%v", err)
		return
	}
	logs.Debug("initialize etcd succ")
//etcd.InitGetEtcd() doit placer-exécuter avant tailf.InitTail(), 
//parce que le résultat de InitGetEtcd() confArray doit alimenter InitTail()!
	err = tailf.InitTail(confArray, appConfig.chanSize)
	//err = tailf.InitTail(appConfig.collectConf, appConfig.chanSize)
	if err !=nil {
		logs.Error("initialize tail failed, err:%v", err)
		return
	}
	logs.Debug("initialize tailf succ")
	
	err = kafka.InitKafka(appConfig.kafkaAddr)
	if err !=nil {
		logs.Error("initialize kafka failed, err:%v", err)
		return
	}
	
	logs.Debug("initialize all succ!")
	/* go func() {
		var count int
		for{
			count++
			logs.Debug("test for logger %d", count)
			time.Sleep (time.Millisecond*100)

		}
	}() */

	err = serverRun()
	if err !=nil {
		logs.Error("serverRun failed, err:%v", err)
		return
	}
	logs.Info("program exited")

}
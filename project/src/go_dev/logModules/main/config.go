package main

import(
	"fmt"
	"github.com/astaxie/beego/config"
	"go_dev/oldBoy/logModules/tailf"
	"errors"
)

var (
	appConfig *Config //ce variable (globale)va être utilisé par tout programme
)

func loadCollectConf(conf config.Configer)(err error){
	//la partie pour obtenir les contenus de [collect] 
//ainsi que "collectConf []tailf.CollectConf", doit déplacer dans ETCD, puis
	var cf tailf.CollectConf

	cf.LogPath = conf.String("collect::log_path")
	if len(cf.LogPath)==0 {
		err = errors.New("invalid collect::log_path")
		return
	}
	
	cf.Topic= conf.String("collect::topic")
	if len(cf.Topic)==0{
		err = errors.New("invalid collect::topic")
		return
	}
	appConfig.collectConf = append(appConfig.collectConf, cf)
	return
}
type Config struct{
	logLevel string
	logPath  string
	
	collectConf []tailf.CollectConf

	chanSize int
	kafkaAddr string

	etcdAddr string
	etcdKey string
}


func loadConf(confType, filename string)(err error){
	conf, err := config.NewConfig(confType, filename)
	if err !=nil {
		fmt.Println("new config failed, err:", err)
		return
	}
	appConfig = &Config{}
	appConfig.logPath = conf.String("logs::log_path")
	if len(appConfig.logPath)==0{
		appConfig.logPath ="./logs/logagent.log"
	}

	appConfig.logLevel = conf.String("logs::log_level")
	if len(appConfig.logPath)==0 {
		appConfig.logLevel = "debug"
	}

	appConfig.chanSize, err = conf.Int("collect::chan_size")
	if err !=nil {
		appConfig.chanSize = 100
	}

	appConfig.kafkaAddr = conf.String("kafka::server_addr")
	if len(appConfig.kafkaAddr)==0{
		err = fmt.Errorf("invalid kafka addresse")
	}

	appConfig.etcdAddr = conf.String("etcd::etcdAddr")
	if len(appConfig.etcdAddr) ==0 {
		err = fmt.Errorf("invalid etcdAddr")
		return
	}
	appConfig.etcdKey = conf.String("etcd::configKey")
	if len(appConfig.etcdKey) == 0 {
		err = fmt.Errorf("invalid etcdKey")
		return
	}

	err = loadCollectConf(conf)
	if err !=nil {
		fmt.Printf("load collect conf failes, err:%v\n", err)
		return
	}
	return
}
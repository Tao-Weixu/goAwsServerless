package etcd

import(
	"go_dev/oldBoy/logModules/tailf"
	"github.com/astaxie/beego/logs"
	"encoding/json"
	"fmt"
	"strings"
	"context"
	"time"
)
func InitGetEtcd(etcdAddr, etcdKey string)(confArray []tailf.CollectConf, err error){

	connOpenEtcd()

	if strings.HasSuffix(etcdKey, "/") == false{
		etcdKey = etcdKey +"/"
	}

	var count int
	for _, ip := range localIPArray {
		count ++
		etcdKeyConfig := fmt.Sprintf("%s%s", etcdKey, ip)
		etcdClient.logKeys= append(etcdClient.logKeys, etcdKeyConfig)
		ctx, cancel :=context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		resp, err :=etcdClient.client.Get(ctx, etcdKeyConfig)
		if err !=nil {
			logs.Error("client get from etcd failed for n°%d config, err:%v\n", count, err)
			continue
		}
		//cancel()
		logs.Debug("resp get key-values from etcd: %v", resp.Kvs)

		for _, v := range resp.Kvs{
			if string(v.Key) == etcdKeyConfig {
				err = json.Unmarshal(v.Value, &confArray)
				if err !=nil {
					logs.Error("unmarchal failed, err:%v", err)
					continue
				}
				logs.Debug("log config from etcd: %v", confArray)
			}
		}
	}
	initWatchEtcd()
	return
}
func initWatchEtcd(){
	for _, logKey :=range etcdClient.logKeys {
		go WatchEtcd(logKey)
	}
}
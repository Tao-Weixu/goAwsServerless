package etcd

import(
	"context"
//"fmt"
	"github.com/astaxie/beego/logs"
	"go_dev/oldBoy/logModules/tailf"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"encoding/json"

)
func WatchEtcd(logKey string){
	
	//connOpenEtcd()
	logs.Debug("begin watch every logKey: %s", logKey)
	for {
		logKeyValues := etcdClient.client.Watch(context.Background(), logKey)
		var confArry []tailf.CollectConf
		var getConfSucc = true
		for wResp := range logKeyValues {
			for _, ev := range wResp.Events {
				if ev.Type == mvccpb.DELETE{
					logs.Warn("logKey[%s]'s config deleted", logKey)
					continue
				}
				if ev.Type == mvccpb.PUT && string(ev.Kv.Key) == logKey {
					err := json.Unmarshal(ev.Kv.Value, &confArry)
					if err !=nil {
						logs.Error("logKey [%s], Unmarchal[%s], err: %v", err)
						getConfSucc = false
						continue
					}
				}

				//fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				logs.Debug("get config from etcd, %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
			if getConfSucc {
				logs.Debug("get config from etcd SUCC, %v", confArry)
				tailf.UpdateConfigEtcd(confArry)
			}
		}
	}
}
package main

import(
	"go_dev/oldBoy/logModules/tailf"
	"github.com/astaxie/beego/logs"
	//"go_dev/oldBoy/logModules/etcd"
	etcd_client "go.etcd.io/etcd/clientv3"
	"encoding/json"
	"fmt"
	"context"
	"time"
)
const (
	EtcdKey = "/oldBoy/logModules/config/192.***.*.**"
)
func main(){
	InitSetConfigToEtcd()
}

func InitSetConfigToEtcd(){
	cli, err :=etcd_client.New(etcd_client.Config{
		Endpoints: []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	defer cli.Close()
	if err !=nil {
		logs.Error("connect failed, err:", err)
		return
	}
	fmt.Println("etcd connect succ +close()")
	//etcd.ConnOpenCloseEtcd()

	var logConfArr []tailf.CollectConf
	logConfArr = append(
		logConfArr, tailf.CollectConf{
			LogPath: "./logs/logagent.log",
			Topic: "nginx_log",
		},
	)
	logConfArr = append(
		logConfArr, tailf.CollectConf{
			LogPath: "./nginx/logs/error.log",
			Topic: "nginx_log_err",
		},
	)
	dataConf, err := json.Marshal(logConfArr)
	if err !=nil {
		fmt.Println("marchal json failed:", err)
		return
	}
	ctx, cancel :=context.WithTimeout(context.Background(), 2*time.Second)
	/* cli.Delete(ctx, EtcdKey)
	defer cancel() 
	return */
	_, err = cli.Put(ctx, EtcdKey, string(dataConf))
	defer cancel()

	if err !=nil {
		fmt.Println("Put congLog failed, err:", err)
		return
	}

/* ctx, cancel :=context.WithTimeout(context.Background(), time.Second)
	resp, err :=etcdClient.client.Get(ctx, EtcdKey)
	defer cancel()
	if err !=nil {
		fmt.Println("Get congLog failed, err:", err)
		return
	}
	for _, ev :=range resp.Kvs{
		fmt.Printf("%s: %s\n", ev.Key, ev.Value)
	}
 */
}

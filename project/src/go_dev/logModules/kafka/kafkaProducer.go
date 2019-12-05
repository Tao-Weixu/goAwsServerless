package kafka

import(
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
	"fmt"
)
var (
	client sarama.SyncProducer
)
func InitKafka(addr string)(err error){
	config :=sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	client, err = sarama.NewSyncProducer([]string{addr}, config)
	//defer client.Close() //surtout il ne faut pas ajouter ce defer!!!!!!!!!!!! 
	if err !=nil {
		//fmt.Println("producer close, err:", err)
		logs.Error("init kafka producer failed, err:", err)
		return
	}
	logs.Debug("initialize kafka succ")
	return
}
func SendToKafka (dataLine, topic string)(err error){
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(dataLine)

	//_, _, err :=client.SendMessage(msg)
	pid, offset, err :=client.SendMessage(msg)
	if err !=nil {
		logs.Error("send message failed, err:%v, dataLine:%v, topic:%v", err, dataLine, topic)
		return
	}
	fmt.Printf("pid:%v, offset:%v, topic:%v\n", pid, offset, topic)
	//logs.Debug("pid:%v, offset:%v, topic:%v\n", pid, offset, topic)
	return
}
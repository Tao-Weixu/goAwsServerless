package main

import(
	"github.com/Shopify/sarama"
	"strings"
	"fmt"
	//"time"
	"sync"
)
var(
	wg sync.WaitGroup
)
func main(){
	kafkaConsumer()
}
func kafkaConsumer(){
	consumer, err :=sarama.NewConsumer(strings.Split("192.***.*.**:9092", ","), nil)
	if err !=nil{
		fmt.Printf("failed to start consummer: %s", err)
	}
	partitionList, err :=consumer.Partitions("nginx_log")
	if err !=nil{
		fmt.Println("failed to get the partitions list of partitions: ", err)
		return
	}
	fmt.Println(partitionList)
	for partition :=range partitionList{
		pc, err :=consumer.ConsumePartition("nginx_log", int32(partition), sarama.OffsetNewest)
		defer pc.AsyncClose()
		if err !=nil{
			fmt.Printf("failed to start consumer for partition %d: %s\n", partition, err)
			return
		}
		go func(pc sarama.PartitionConsumer){
			wg.Add(1)
			for msg := range pc.Messages(){
				fmt.Printf("Patition:%d, Offset:%d, key:%s, Value:%s", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
				fmt.Println()
			}
			wg.Done()
		}(pc)
	}
	//time.Sleep(time.Hour)
	wg.Wait()
	consumer.Close()
}
package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"sync"
)

var (
	wg sync.WaitGroup
)

func main() {
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, nil)
	if err != nil {
		fmt.Println("new consumer error = ", err)
		return
	}
	//Partitions(topic): 该方法返回了该topic的所有分区id
	partitionList, err := consumer.Partitions("mykafka")
	if err != nil {
		panic(err)
	}
	for partition := range partitionList {
		//consumePartition方法根据主题，区分和给定的偏移量创建了相应的分区消费者
		//如果该分区消费者已经消费了该消息将返回error
		//sarama。offsetNewest： 表明了为最新消息
		pc, err := consumer.ConsumePartition("mykafka", int32(partition), sarama.OffsetNewest)
		if err != nil {
			panic(err)
		}
		wg.Add(1)
		go func(partitionConsumer sarama.PartitionConsumer) {
			defer partitionConsumer.AsyncClose()
			defer wg.Done()
			//Messages()方法返回一个消费消息类型的只读通道，由代理产生
			for msg := range partitionConsumer.Messages() {
				fmt.Printf("%s -- Partition: %d, Offset: %d, Key: %s, Value: %s\n",
					msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
			}
		}(pc)
	}
	wg.Wait()
	consumer.Close()
}

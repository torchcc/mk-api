package kafka

import "C"

type Config struct {
	Host        []string `json:"host"`
	Topic       string   `json:"topic"`
	Description string   `json:"description"`
}

/** it's a must to call pc.AsyncClose() after getting pc out of the list to use
using example
defer pc.AsyncClose()
go func(partitionConsumer sarama.PartitionConsumer) {
	for msg := range pc.Messages() {
		fmt.Printf("Partition: %d Offset: %d Key: %v Value: %v\n", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
	}
}(pc)
**/
// func NewKafkaConsumers(conf *Config) (pcs []*sarama.PartitionConsumer) {
// 	consumer, err := sarama.NewConsumer(conf.Host, nil)
// 	if err != nil {
// 		fmt.Printf("fail to start consumer, err: %v\n", err)
// 		panic(err)
// 		return
// 	}
// 	partitionList, err := consumer.Partitions(conf.Topic)
// 	if err != nil {
// 		fmt.Printf("fail to get list of partition. err: %v\n", err)
// 		panic(err)
// 		return
// 	}
// 	fmt.Printf("the partition list is: %v\n", partitionList)
//
// 	for partition := range partitionList { // 遍历所有partition
// 		// 针对每个分区创建一个对应的分区消费者
// 		pc, err := consumer.ConsumePartition(conf.Topic, int32(partition), sarama.OffsetNewest)
// 		if err != nil {
// 			fmt.Printf("fail to start consumer for partition %d, err: %v\n", partition, err)
// 			panic(err)
// 			return
// 		}
// 		pcs = append(pcs, &pc)
// 	}
// 	return
// }

// func NewKafkaConsumers(conf *Config, partitions ...int32) ([]chan *sarama.ConsumerMessage, error) {
// 	msgChans := make([]chan *sarama.ConsumerMessage, 0)
// 	consumer, err := sarama.NewConsumer(conf.Host, nil)
//
// 	if err != nil {
// 		fmt.Printf("fail to start consumer, err: %v\n", err)
// 		return msgChans, err
// 	}
//
// 	var partitionList []int32
// 	if len(partitions) == 0 {
// 		partitionList, err = consumer.Partitions(conf.Topic)
// 		if err != nil {
// 			fmt.Printf("fail to get list of partition. err: %v\n", err)
// 			return msgChans, err
// 		}
// 	} else {
// 		partitionList = partitions
// 	}
//
//
// 	fmt.Printf("the partition list is: %v\n", partitionList)
//
// 	for _, partition := range partitionList { // 遍历所有partition
// 		// 针对每个分区创建一个对应的分区消费者
// 		pc, err := consumer.ConsumePartition(conf.Topic, partition, sarama.OffsetOldest)
// 		if err != nil {
// 			fmt.Printf("fail to start consumer for partition %d, err: %v\n", partition, err)
// 			return msgChans, err
// 		}
// 		defer pc.AsyncClose()
//
// 		msgs := make(chan *sarama.ConsumerMessage)
// 		msgChans = append(msgChans, pc.Messages())
// 		// go func(partitionConsumer sarama.PartitionConsumer, msgs chan *sarama.ConsumerMessage) {
// 		// 	curCh := pc.Messages()
// 		// 	for {
// 		// 		select {
// 		// 		case msg :=  <-curCh:
// 		//
// 		// 			fmt.Println("hello world ")
// 		// 			fmt.Printf("Partition: %d Offset: %d Key: %v Value: %v\n", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
// 		// 			msgs <- msg
// 		// 		}
// 		// 	}
// 		//
// 		// 	// for msg := range pc.Messages() {
// 		// 	// 	fmt.Printf("Partition: %d Offset: %d Key: %v Value: %v\n", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
// 		// 	// 	msgs <- msg
// 		// 	// }
// 		// }(pc, msgs)
//
// 	}
// 	return msgChans, err
// }

package kafka

// func TestKafka(t *testing.T) {
// 	cnf := Config{
// 		Host:        []string{"192.168.10.139:31092", "192.168.10.139:31192", "192.168.10.139:31292"},
// 		Topic:       "monitor_order",
// 		Description: "",
// 	}
//
// 	intKey := 0
// 	pcs := NewKafkaConsumers(&cnf)
// 	for _, pc := range pcs {
// 		defer (*pc).AsyncClose()
//
// 		go func(partitionConsumer sarama.PartitionConsumer) {
// 			for msg := range (*pc).Messages() {
// 				t.Logf("Partition: %d Offset: %d Key: %v Value: %v\n", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
// 				intKey, _ = strconv.Atoi(string(msg.Value))
// 				fmt.Println(intKey)
// 			}
// 		}(*pc)
//
// 		select {
//
// 		}
// 	}
//
// }
//

//
// func TestKafka(t *testing.T) {
// 	cnf := Config{
// 		Host:        []string{"192.168.10.139:31092", "192.168.10.139:31192", "192.168.10.139:31292"},
// 		Topic:       "monitor_order",
// 		Description: "",
// 	}
//
// 	chs, err := NewKafkaConsumers(&cnf)
// 	if err != nil {
// 		t.Errorf("failed to create kafka consumer: %v\n", err)
// 		return
// 	}
// 	fmt.Println("heer")
//
// 	for _, ch := range chs {
//
// 			for {
// 				select {
// 				case msg := <-ch:
// 					fmt.Printf("Partition: %d Offset: %d Key: %v Value: %v\n", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
// 				}
// 				xtime.Sleep(1 * xtime.Second)
//
// 			}
// 	}
//
// 	select {
//
// 	}
//
//
// }

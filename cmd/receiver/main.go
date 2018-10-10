package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
	"github.com/matzew/kafka-receiver/pkg/config"
)

func main() {

	config := config.GetConfig()

	consumer, err := sarama.NewConsumer([]string{config.BootStrapServers}, nil)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// OffsetNewest
	initialOffset := sarama.OffsetOldest //get offset for the oldest message on the topic
	partitionConsumer, err := consumer.ConsumePartition(config.KafkaTopic, 0, initialOffset)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	consumed := 0
ConsumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Consumed message offset %s", msg.Value)
			consumed++
		case <-signals:
			break ConsumerLoop
		}
	}

	log.Printf("Consumed: %d\n", consumed)
}

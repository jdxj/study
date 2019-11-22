package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

func main() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewAsyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		panic(err)
	}

	// Trap SIGINT to trigger a graceful shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	var (
		wg                          sync.WaitGroup
		enqueued, successes, errors int
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range producer.Successes() {
			successes++
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range producer.Errors() {
			log.Println(err)
			errors++
		}
	}()

	i := 0

ProducerLoop:
	for {
		data := fmt.Sprintf("testing %d", i)
		message := &sarama.ProducerMessage{Topic: "my_topic", Value: sarama.StringEncoder(data)}
		select {
		case producer.Input() <- message:
			enqueued++
			fmt.Println(data)
		case <-signals:
			producer.AsyncClose() // Trigger a shutdown of the producer.
			break ProducerLoop
		}
		i++

		time.Sleep(time.Second)
	}
	wg.Wait()

	log.Printf("Successfully produced: %d; errors: %d\n", successes, errors)
}

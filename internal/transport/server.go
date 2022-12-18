package transport

import (
	"fmt"
	"github.com/mrasnake/messageQueue/cmd/run_server/service"
	"github.com/streadway/amqp"
	"github.com/urfave/cli/v2"
	"log"
)

type server struct {
	QueueConn string
	QueueName string
	LogFileName string

}

type Message struct{
	Request string `json:"request"`
	Value string `json:"value"`
}

// defineSettings gathers the configurations and creates the server object.
func defineSettings(ctx *cli.Context) *server {

	out := &server{
		QueueConn: ctx.String("connection"),
		QueueName: ctx.String("queue"),
		LogFileName: ctx.String("logs"),
	}
	return out
}

// Run sets up the message queue and runs the server waiting for requests.
func (s *server) Run() error{
	conn, err := amqp.Dial(s.QueueConn)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		s.QueueName, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return fmt.Errorf("unable to declare consumer: %w", err)
	}

	svr, err := service.NewService(s.LogFileName)
	if err != nil{
		return fmt.Errorf("service not created: %w", err)
	}

	forever := make(chan bool)

	go func() {
		for m := range msgs {
			svr.ProcessMessage(m.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

	return nil
}

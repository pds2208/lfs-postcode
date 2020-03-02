package messaging

import (
    "encoding/json"
    "github.com/rs/zerolog/log"
    "github.com/streadway/amqp"
)

type Messaging struct {
    Connection amqp.Connection
}

func NewConnection() *Messaging {
    a := &Messaging{}
    a.Connect()
    return a
}

func (m Messaging) Connect() {
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    failOnError(err, "Failed to connect to RabbitMQ")
    defer conn.Close()

    ch, err := conn.Channel()
    failOnError(err, "Failed to open a channel")
    defer ch.Close()

    msgs, err := ch.Consume(
        "lfs-postcode", // queue
        "",             // consumer
        true,           // auto-ack
        false,          // exclusive
        false,          // no-local
        false,          // no-wait
        nil,            // args
    )

    failOnError(err, "Failed to register a consumer")

    forever := make(chan bool)

    go func() {
        for d := range msgs {
            handleMessage(d.Body)
        }
    }()

    log.Info().Msg("Waiting for messages...")
    <-forever
}

func handleMessage(message []byte) {
    res := Message{}
    err := json.Unmarshal([]byte(message), &res)
    if err != nil {
        log.Error().Err(err)
        return
    }
    log.Printf("Received request for year: %d, month %d", res.Year, res.Month)
}

func failOnError(err error, msg string) {
    if err != nil {
        log.Fatal().
            Err(err).
            Str("Rabbit MQ", "Error")
    }
}

type Message struct {
    Month int `json:"month"`
    Year  int `json:"year"`
}

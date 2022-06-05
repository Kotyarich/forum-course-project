package server

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
	"os"
	"os/signal"
	"statistic-service/statistic"
	statisticHttp "statistic-service/statistic/delivery/http"
	"statistic-service/statistic/queue-handler"
	"statistic-service/statistic/repository/postgres"
	statisticUceCase "statistic-service/statistic/usecase"
	"time"
)

type App struct {
	httpServer *http.Server

	statisticUC  statistic.UseCase
	queueHandler *queue_handler.QueueHandler
}

func NewApp() *App {
	statisticRepo := postgres.NewStatisticRepository()

	return &App{
		statisticUC:  statisticUceCase.NewStatisticUseCase(statisticRepo),
		queueHandler: queue_handler.NewQueueHandler(statisticRepo),
	}
}

func (a *App) Run(port string) error {
	router := echo.New()

	statisticHttp.RegisterHTTPEndpoints(router, a.statisticUC)
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5002"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))

	router.Use(middleware.Logger())

	router.Static("/", "swaggerui")

	a.httpServer = &http.Server{
		Addr:           port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return fmt.Errorf("%v: failed to connect to RabbitMQ", err)
	}
	defer func() { _ = conn.Close() }()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("%v: failed to open a channel", err)
	}
	defer func () { _ = ch.Close() }()

	go func() {
		err := a.RunQueuePostConsumers(ch)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s", err.Error())
		}
	}()
	go func() {
		err := a.RunQueueUserConsumers(ch)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s", err.Error())
		}
	}()
	go func() {
		err := a.RunQueueVoteConsumers(ch)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s", err.Error())
		}
	}()
	go func() {
		err := a.RunQueueThreadConsumers(ch)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s", err.Error())
		}
	}()
	go func() {
		err := a.RunQueueForumConsumers(ch)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}

func (a *App) RunQueuePostConsumers(ch *amqp.Channel) error {
	q, err := ch.QueueDeclare(
		"post-created", // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return fmt.Errorf("%v: failed to declare a queue", err)
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
		return fmt.Errorf("%v: failed to register a consumer", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			err = a.queueHandler.HandlePostCreation(context.Background(), string(d.Body))
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "%s", err.Error())
			}
		}
	}()

	<-forever
	return nil
}

func (a *App) RunQueueUserConsumers(ch *amqp.Channel) error {
	q, err := ch.QueueDeclare(
		"user-created", // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return fmt.Errorf("%v: failed to declare a queue", err)
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
		return fmt.Errorf("%v: failed to register a consumer", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			err = a.queueHandler.HandleUserCreation(context.Background(), string(d.Body))
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "%s", err.Error())
			}
		}
	}()

	<-forever
	return nil
}

func (a *App) RunQueueVoteConsumers(ch *amqp.Channel) error {
	q, err := ch.QueueDeclare(
		"vote-created", // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return fmt.Errorf("%v: failed to declare a queue", err)
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
		return fmt.Errorf("%v: failed to register a consumer", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			err = a.queueHandler.HandleVoteCreation(context.Background(), string(d.Body))
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "%s", err.Error())
			}
		}
	}()

	<-forever
	return nil
}

func (a *App) RunQueueThreadConsumers(ch *amqp.Channel) error {
	q, err := ch.QueueDeclare(
		"thread-created", // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return fmt.Errorf("%v: failed to declare a queue", err)
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
		return fmt.Errorf("%v: failed to register a consumer", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			err = a.queueHandler.HandleThreadCreation(context.Background(), string(d.Body))
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "%s", err.Error())
			}
		}
	}()

	<-forever
	return nil
}

func (a *App) RunQueueForumConsumers(ch *amqp.Channel) error {
	q, err := ch.QueueDeclare(
		"forum-created", // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return fmt.Errorf("%v: failed to declare a queue", err)
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
		return fmt.Errorf("%v: failed to register a consumer", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			err = a.queueHandler.HandleForumCreation(context.Background(), string(d.Body))
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "%s", err.Error())
			}
		}
	}()

	<-forever
	return nil
}

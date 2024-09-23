package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/stanislavCasciuc/common"
	"github.com/stanislavCasciuc/common/broker"
	"github.com/stanislavCasciuc/common/discovery"
	"github.com/stanislavCasciuc/common/discovery/consul"
	stripeProcessor "github.com/stanislavCasciuc/payments/processor/stripe"
	"github.com/stripe/stripe-go/v79"
	"google.golang.org/grpc"
)

var (
	grpcAddr    = common.EnvString("GRPC_ADDR", "localhost:2001")
	serviceName = "payments"
	consulAddr  = common.EnvString("CONSUL_ADDR", "localhost:8500")
	amqpUser    = common.EnvString("AMQP_USER", "guest")
	amqpHost    = common.EnvString("AMQP_USER", "localhost")
	amqpPort    = common.EnvString("AMQP_USER", "5672")
	amqpPass    = common.EnvString("AMQP_USER", "guest")
	stripeKey   = common.EnvString("STRIPE_KEY", "")
)

func main() {
	// Register consul
	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, grpcAddr); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				log.Fatal("failed to heath check")
			}
			time.Sleep(time.Second * 1)
		}
	}()

	defer registry.Deregister(ctx, instanceID, serviceName)

	stripe.Key = stripeKey

	// Broker connection
	ch, close := broker.Connect(amqpUser, amqpPass, amqpHost, amqpPort)
	defer func() {
		close()
		ch.Close()
	}()

	stripeProcessor := stripeProcessor.NewProcessor()
	svc := NewService(stripeProcessor)

	amqpConsumer := NewConsumer(svc)
	go amqpConsumer.Listen(ch)

	// gRPC server
	grpcServer := grpc.NewServer()

	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer l.Close()

	log.Println("GRPC started at grpcAddr: ", grpcAddr)
	if err := grpcServer.Serve(l); err != nil {
		log.Fatal(err.Error())
	}

	// Graceful Shut Down
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-term
		if err := l.Close(); err != nil {
			log.Fatalf("Error closing listener: %v", err)
		}
		grpcServer.Stop()
	}()
}

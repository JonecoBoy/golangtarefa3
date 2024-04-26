package main

import (
	"database/sql"
	"fmt"
	"github.com/jonecoboy/golangtarefa3/internal/event"
	"github.com/jonecoboy/golangtarefa3/internal/infra/database"
	"github.com/jonecoboy/golangtarefa3/internal/infra/grpc/pb"
	"github.com/jonecoboy/golangtarefa3/internal/infra/grpc/service"
	"github.com/jonecoboy/golangtarefa3/internal/infra/web"
	"github.com/jonecoboy/golangtarefa3/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jonecoboy/golangtarefa3/configs"
	"github.com/jonecoboy/golangtarefa3/graph"
	"github.com/jonecoboy/golangtarefa3/internal/event/handler"
	"github.com/jonecoboy/golangtarefa3/internal/infra/web/webserver"
	"github.com/jonecoboy/golangtarefa3/pkg/events"
	"github.com/streadway/amqp"
	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	loadConfig, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(loadConfig.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", loadConfig.DBUser, loadConfig.DBPassword, loadConfig.DBHost, loadConfig.DBPort, loadConfig.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel(loadConfig.RabbitMQConnection)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	orderRepository := database.NewOrderRepository(db)
	orderCreated := event.NewOrderCreated()
	createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepository, orderCreated, eventDispatcher)

	//createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)

	newWebServer := webserver.NewWebServer(loadConfig.WebServerPort)
	webOrderHandler := web.NewWebOrderHandler(eventDispatcher, orderRepository, orderCreated)

	newWebServer.AddHandler("/order", webOrderHandler.Create)

	fmt.Println("Starting web server on port", loadConfig.WebServerPort)
	go newWebServer.Start()

	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", loadConfig.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", loadConfig.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", loadConfig.GraphQLServerPort)
	http.ListenAndServe(":"+loadConfig.GraphQLServerPort, nil)
}

func getRabbitMQChannel(connection string) *amqp.Channel {
	conn, err := amqp.Dial(connection)
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}

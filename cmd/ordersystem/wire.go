//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"

	"github.com/google/wire"
	"github.com/jonecoboy/golangtarefa3/internal/entity"
	"github.com/jonecoboy/golangtarefa3/internal/event"
	"github.com/jonecoboy/golangtarefa3/internal/infra/database"
	"github.com/jonecoboy/golangtarefa3/internal/infra/web"
	"github.com/jonecoboy/golangtarefa3/internal/usecase"
	"github.com/jonecoboy/golangtarefa3/pkg/events"
)

var setOrderRepositoryDependency = wire.NewSet(
	wire.NewSet(database.NewOrderRepository),
	wire.Bind(new(entity.OrderRepositoryInterface), new(*database.OrderRepository)),
)

var setEventDispatcherDependency = wire.NewSet(
	wire.NewSet(events.NewEventDispatcher, event.NewOrderCreated),
	// toda vez que voce ver a interface EventInterface voce via chamar o orderCreated!. TODO LUGAR QUE TIVER INTERFACE VAI TROCAR PELO OBJETO CONCRETO ORDER CREATED

	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
	// toda vez que você ver no código EventDispatcherInterface instancie o EventDispatcher
	wire.Bind(new(events.EventDispatcherInterface), new(*events.EventDispatcher)),
)

var setOrderCreatedEvent = wire.NewSet(
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
)

func NewCreateOrderUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.CreateOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		usecase.NewCreateOrderUseCase,
	)
	return &usecase.CreateOrderUseCase{}
}

func NewWebOrderHandler(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *web.WebOrderHandler {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		// essas struct aqui tem algumas interfaces!! OrderRepositoryInterface
		web.NewWebOrderHandler,
	)
	return &web.WebOrderHandler{}
}

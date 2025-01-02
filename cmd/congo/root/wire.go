//go:build wireinject
// +build wireinject

package root

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/wire"
	"github.com/henriquemarlon/congo/configs"
	"github.com/henriquemarlon/congo/internal/domain/entity"
	event_handler "github.com/henriquemarlon/congo/internal/domain/event/handler"
	"github.com/henriquemarlon/congo/internal/infra/repository"
	web_handler "github.com/henriquemarlon/congo/internal/infra/web"
	"github.com/henriquemarlon/congo/internal/usecase"
	"github.com/henriquemarlon/congo/pkg/events"
)

var setMongoCollectionProvider = wire.NewSet(
	configs.SetupMongoCollection,
	repository.NewSensorRepository,
	wire.Bind(new(entity.SensorRepository), new(*repository.SensorRepository)),
)

var setBrokerProvider = wire.NewSet(
	configs.SetupBrokerClient,
)

var setEventDispatcherProvider = wire.NewSet(
	events.NewEventDispatcher,
)

func NewBrokerClient(clientId string) (MQTT.Client, error) {
	wire.Build(setBrokerProvider)
	return nil, nil
}

func NewDataEmittedEventHandler(client MQTT.Client) (*event_handler.DataEmittedHandler, error) {
	wire.Build(event_handler.NewDataEmittedHandler)
	return nil, nil
}

func NewEmitDataUseCase(event events.EventInterface, eventDispatcher events.EventDispatcherInterface) (*usecase.EmitDataUseCase, error) {
	wire.Build(
		setMongoCollectionProvider,
		usecase.NewEmitDataUseCase,
	)
	return nil, nil
}

func NewFindAllSensorsUseCase() (*usecase.FindAllSensorsUseCase, error) {
	wire.Build(
		setMongoCollectionProvider,
		usecase.NewFindAllSensorsUseCase,
	)
	return nil, nil
}

func NewSensorWebHandlers(sensorChannel chan *entity.Sensor) (*web_handler.SensorHandlers, error) {
	wire.Build(
		setMongoCollectionProvider,
		web_handler.NewSensorHandlers,
	)
	return nil, nil
}

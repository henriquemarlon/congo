// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package root

import (
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/google/wire"
	"github.com/henriquemarlon/congo/configs"
	"github.com/henriquemarlon/congo/internal/domain/entity"
	"github.com/henriquemarlon/congo/internal/domain/event/handler"
	"github.com/henriquemarlon/congo/internal/infra/repository"
	"github.com/henriquemarlon/congo/internal/infra/web"
	"github.com/henriquemarlon/congo/internal/usecase"
	"github.com/henriquemarlon/congo/pkg/events"
)

// Injectors from wire.go:

func NewBrokerClient(clientId string) (mqtt.Client, error) {
	client, err := configs.SetupBrokerClient(clientId)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewDataEmittedEventHandler(client mqtt.Client) (*handler.DataEmittedHandler, error) {
	dataEmittedHandler := handler.NewDataEmittedHandler(client)
	return dataEmittedHandler, nil
}

func NewEmitDataUseCase(event events.EventInterface, eventDispatcher events.EventDispatcherInterface) (*usecase.EmitDataUseCase, error) {
	collection, err := configs.SetupMongoCollection()
	if err != nil {
		return nil, err
	}
	sensorRepository := repository.NewSensorRepository(collection)
	emitDataUseCase := usecase.NewEmitDataUseCase(event, sensorRepository, eventDispatcher)
	return emitDataUseCase, nil
}

func NewFindAllSensorsUseCase() (*usecase.FindAllSensorsUseCase, error) {
	collection, err := configs.SetupMongoCollection()
	if err != nil {
		return nil, err
	}
	sensorRepository := repository.NewSensorRepository(collection)
	findAllSensorsUseCase := usecase.NewFindAllSensorsUseCase(sensorRepository)
	return findAllSensorsUseCase, nil
}

func NewSensorWebHandlers(sensorChannel chan *entity.Sensor) (*web.SensorHandlers, error) {
	collection, err := configs.SetupMongoCollection()
	if err != nil {
		return nil, err
	}
	sensorRepository := repository.NewSensorRepository(collection)
	sensorHandlers := web.NewSensorHandlers(sensorChannel, sensorRepository)
	return sensorHandlers, nil
}

// wire.go:

var setMongoCollectionProvider = wire.NewSet(configs.SetupMongoCollection, repository.NewSensorRepository, wire.Bind(new(entity.SensorRepository), new(*repository.SensorRepository)))

var setBrokerProvider = wire.NewSet(configs.SetupBrokerClient)

var setEventDispatcherProvider = wire.NewSet(events.NewEventDispatcher)

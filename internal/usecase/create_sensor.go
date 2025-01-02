package usecase

import (
	"context"

	"github.com/henriquemarlon/congo/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateSensorUseCase struct {
	SensorRepository entity.SensorRepository
}

type CreateSensorInputDTO struct {
	Name      string                  `json:"name"`
	Latitude  float64                 `json:"latitude"`
	Longitude float64                 `json:"longitude"`
	Params    map[string]entity.Param `json:"params"`
}

type CreateSensorOutputDTO struct {
	Id        primitive.ObjectID      `json:"id"`
	Name      string                  `json:"name"`
	Latitude  float64                 `json:"latitude"`
	Longitude float64                 `json:"longitude"`
	Params    map[string]entity.Param `json:"params"`
}

func NewCreateSensorUseCase(sensorRepository entity.SensorRepository) *CreateSensorUseCase {
	return &CreateSensorUseCase{SensorRepository: sensorRepository}
}

func (c *CreateSensorUseCase) Execute(ctx context.Context, input CreateSensorInputDTO) (*CreateSensorOutputDTO, error) {
	sensor := entity.NewSensor(input.Name, input.Latitude, input.Longitude, input.Params)
	res, err := c.SensorRepository.CreateSensor(ctx, sensor)
	if err != nil {
		return nil, err
	}
	return &CreateSensorOutputDTO{
		Id:        res.Id,
		Name:      res.Name,
		Latitude:  res.Latitude,
		Longitude: res.Longitude,
		Params:    res.Params,
	}, nil
}

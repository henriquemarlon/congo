package usecase

import (
	"context"
	"math"
	"time"

	"github.com/henriquemarlon/congo/internal/domain/entity"
	"github.com/henriquemarlon/congo/pkg/events"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat"
)

type EmitDataUseCase struct {
	EmitData         events.EventInterface
	SensorRepository entity.SensorRepository
	EventDispatcher  events.EventDispatcherInterface
}

type EmitDataInputDTO struct {
	Id primitive.ObjectID `json:"id"`
}

type EmitDataOutputDTO struct {
	Id        primitive.ObjectID     `json:"id"`
	Name      string                 `json:"name"`
	Latitude  float64                `json:"latitude"`
	Longitude float64                `json:"longitude"`
	Data      map[string]interface{} `json:"data"`
}

func NewEmitDataUseCase(
	emitData events.EventInterface,
	sensorRepository entity.SensorRepository,
	eventDispatcher events.EventDispatcherInterface,
) *EmitDataUseCase {
	return &EmitDataUseCase{
		EmitData:         emitData,
		SensorRepository: sensorRepository,
		EventDispatcher:  eventDispatcher,
	}
}

func (e *EmitDataUseCase) Execute(ctx context.Context, input EmitDataInputDTO) error {
	res, err := e.SensorRepository.FindSensorById(ctx, input.Id)
	if err != nil {
		return err
	}

	data := make(map[string]interface{})
	for key, interval := range res.Params {
		intervalValues := make([]float64, int(interval.Max-interval.Min)+1)
		for i := range intervalValues {
			intervalValues[i] = float64(interval.Min) + float64(i)
		}
		mean, stdDev := stat.MeanStdDev(intervalValues, nil)
		factor := stdDev / math.Sqrt(float64(len(intervalValues)))
		confidenceInterval := []float64{mean - interval.Factor*factor, mean + interval.Factor*factor}

		rand.NewSource(uint64(time.Now().UnixNano()))
		value := math.Round(rand.Float64()*(confidenceInterval[0]-confidenceInterval[1]) + confidenceInterval[1])
		data[key] = value
	}

	dto := &EmitDataOutputDTO{
		Id:        res.Id,
		Name:      res.Name,
		Latitude:  res.Latitude,
		Longitude: res.Longitude,
		Data:      data,
	}

	e.EmitData.SetPayload(dto)
	if err := e.EventDispatcher.Dispatch(e.EmitData); err != nil {
		return err
	}
	return nil
}

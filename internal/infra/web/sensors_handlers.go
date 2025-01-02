package web

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/henriquemarlon/congo/internal/domain/entity"
	"github.com/henriquemarlon/congo/internal/infra/repository"
	"github.com/henriquemarlon/congo/internal/usecase"
)

type SensorHandlers struct {
	SensorChannel    chan *entity.Sensor
	SensorRepository *repository.SensorRepository
}

func NewSensorHandlers(
	sensorChannel chan *entity.Sensor,
	sensorRepository *repository.SensorRepository,
) *SensorHandlers {
	return &SensorHandlers{
		SensorChannel:    sensorChannel,
		SensorRepository: sensorRepository,
	}
}

func (s *SensorHandlers) CreateSensorHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var input usecase.CreateSensorInputDTO
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	createSensor := usecase.NewCreateSensorUseCase(s.SensorRepository)
	output, err := createSensor.Execute(ctx, input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	s.SensorChannel <- &entity.Sensor{
		Id:        output.Id,
		Name:      output.Name,
		Latitude:  output.Latitude,
		Longitude: output.Longitude,
		Params:    output.Params,
	}
	close(s.SensorChannel)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}

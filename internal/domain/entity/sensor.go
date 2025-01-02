package entity

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrSensorNotFound = errors.New("sensor not found")
)

type SensorRepository interface {
	CreateSensor(ctx context.Context, sensor *Sensor) (*Sensor, error)
	FindSensorById(ctx context.Context, id primitive.ObjectID) (*Sensor, error)
	FindAllSensors(ctx context.Context) ([]*Sensor, error)
}

type Sensor struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `json:"name"`
	Latitude  float64            `json:"latitude"`
	Longitude float64            `json:"longitude"`
	Params    map[string]Param   `json:"params"`
}

type Param struct {
	Min    int     `json:"min"`
	Max    int     `json:"max"`
	Factor float64 `json:"z"`
}

func NewSensor(name string, latitude float64, longitude float64, params map[string]Param) *Sensor {
	sensor := &Sensor{
		Name:      name,
		Latitude:  latitude,
		Longitude: longitude,
		Params:    params,
	}
	if err := sensor.Validate(); err != nil {
		return nil
	}
	return sensor
}

func (s *Sensor) Validate() error {
	if s.Name == "" {
		return errors.New("name is required")
	}
	if s.Latitude == 0 {
		return errors.New("latitude is required")
	}
	if s.Longitude == 0 {
		return errors.New("longitude is required")
	}
	if len(s.Params) == 0 {
		return errors.New("params is required")
	}
	return nil
}

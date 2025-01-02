package event

import (
	"fmt"
	"time"
)

type DataEmitted struct {
	Name    string
	Payload interface{}
}

func NewDataEmitted(sensorId string) *DataEmitted {
	return &DataEmitted{
		Name: fmt.Sprintf("DataEmitted_%s", sensorId),
	}
}

func (e *DataEmitted) GetName() string {
	return e.Name
}

func (e *DataEmitted) GetPayload() interface{} {
	return e.Payload
}

func (e *DataEmitted) SetPayload(payload interface{}) {
	e.Payload = payload
}

func (e *DataEmitted) GetDateTime() time.Time {
	return time.Now()
}

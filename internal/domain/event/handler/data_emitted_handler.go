package handler

import (
	"encoding/json"
	"log/slog"
	"sync"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/henriquemarlon/congo/pkg/events"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DataEmittedHandler struct {
	Client MQTT.Client
}

func NewDataEmittedHandler(client MQTT.Client) *DataEmittedHandler {
	return &DataEmittedHandler{
		Client: client,
	}
}

func (h *DataEmittedHandler) Handle(event events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	payload := event.GetPayload()

	bytesPayload, err := json.Marshal(payload)
	if err != nil {
		slog.Error("Error serializing the payload", "error", err)
	}

	token := h.Client.Publish("sensors", 1, false, bytesPayload)
	if token.Error() != nil {
		slog.Error("Failed to publish the message", "error", token.Error())
	}
	token.Wait()

	var p struct {
		Id        primitive.ObjectID `json:"id"`
		Name      string             `json:"name"`
		Latitude  float64            `json:"latitude"`
		Longitude float64            `json:"longitude"`
		Data      map[string]interface{}
	}
	err = json.Unmarshal(bytesPayload, &p)
	if err != nil {
		slog.Error("Error deserializing the payload", "error", err)
	}

	data, err := json.Marshal(p.Data)
	if err != nil {
		slog.Error("Error serializing the data", "error", err)
	}
	slog.Debug("Data emitted", "sensor", p.Id.Hex(), "name", p.Name, "latitude", p.Latitude, "longitude", p.Longitude, "data", string(data))
}

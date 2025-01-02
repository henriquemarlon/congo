package configs

import (
	"fmt"
	"os"

	"log/slog"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func SetupBrokerClient(clientId string) (MQTT.Client, error) {
	serverURL, ok := os.LookupEnv("HIVEMQ_SERVER_URL")
	if !ok {
		slog.Error("HIVEMQ_SERVER_URL not set")
		os.Exit(1)
	}

	username, ok := os.LookupEnv("HIVEMQ_USERNAME")
	if !ok {
		slog.Error("HIVEMQ_USERNAME not set")
		os.Exit(1)
	}

	password, ok := os.LookupEnv("HIVEMQ_PASSWORD")
	if !ok {
		slog.Error("HIVEMQ_PASSWORD not set")
		os.Exit(1)
	}

	opts := MQTT.NewClientOptions().
		AddBroker(serverURL).
		SetUsername(username).
		SetPassword(password).
		SetClientID(clientId)

	client := MQTT.NewClient(opts)
	if session := client.Connect(); session.Wait() && session.Error() != nil {
		return nil, fmt.Errorf("failed to connect to MQTT broker: %v", session.Error())
	}
	return client, nil
}

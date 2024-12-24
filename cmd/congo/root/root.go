package root

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/henriquemarlon/congo/configs"
	"github.com/henriquemarlon/congo/internal/domain/entity"
	"github.com/henriquemarlon/congo/internal/domain/event"
	"github.com/henriquemarlon/congo/internal/usecase"
	"github.com/henriquemarlon/congo/pkg/events"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
)

const CMD_NAME = "congo"

var (
	configPath string
	verbose    bool
	Cmd        = &cobra.Command{
		Use:   CMD_NAME,
		Short: "Congo is sensor simulator for smart city system development",
		Long:  "This CLI provides a scalable simulation of sensors, supporting any type and handling thousands of concurrent units",
		Run:   run,
	}
)

var startupMessage = `
Congo is ready with NUMBER sensors, and the following configuration:
MongoDB url MONGO_DB_URL - [ Database: MONGO_DB_NAME, Collection: MONGO_COLLECTION_NAME ]
HiveMQ url HIVEMQ_SERVER_URL - [ Username: HIVEMQ_USERNAME, Password: HIVEMQ_PASSWORD ]

Press Ctrl+C to stop the application.
`

func init() {
	Cmd.Flags().BoolVar(&verbose, "verbose", false, "Show detailed output, including sensitive information")

	Cmd.Flags().StringVar(&configPath, "config", "", "Path to the configuration file (required)")
	if err := Cmd.MarkFlagRequired("config"); err != nil {
		os.Exit(1)
	}

	Cmd.PreRun = func(cmd *cobra.Command, args []string) {
		if verbose {
			configs.ConfigureLogger(slog.LevelDebug)
		} else {
			configs.ConfigureLogger(slog.LevelInfo)
		}
	}
}

func run(cmd *cobra.Command, args []string) {
	if err := configs.LoadConfig(configPath); err != nil {
		slog.Error("Failed to load configuration file", "error", err)
		os.Exit(1)
	}

	ctx := cmd.Context()

	eventDispatcher := events.NewEventDispatcher()

	findAllSensors, err := NewFindAllSensorsUseCase()
	if err != nil {
		slog.Error("Failed to instantiate find all sensors usecase", "error", err)
		os.Exit(1)
	}

	sensors, err := findAllSensors.Execute(ctx)
	if err != nil {
		slog.Error("Failed to find all sensors", "error", err)
		os.Exit(1)
	}

	password := "********"
	if verbose {
		password = os.Getenv("HIVEMQ_PASSWORD")
	}

	message := strings.NewReplacer(
		"NUMBER", strconv.Itoa(len(sensors)),
		"MONGO_DB_URL", os.Getenv("MONGO_DB_URL"),
		"MONGO_DB_NAME", os.Getenv("MONGO_DB_NAME"),
		"MONGO_COLLECTION_NAME", os.Getenv("MONGO_COLLECTION_NAME"),
		"HIVEMQ_SERVER_URL", os.Getenv("HIVEMQ_SERVER_URL"),
		"HIVEMQ_USERNAME", os.Getenv("HIVEMQ_USERNAME"),
		"HIVEMQ_PASSWORD", password,
	).Replace(startupMessage)

	fmt.Println(message)

	sensorChannel := make(chan *entity.Sensor, len(sensors))
	for _, sensor := range sensors {
		sensorChannel <- &entity.Sensor{
			Id:        sensor.Id,
			Name:      sensor.Name,
			Latitude:  sensor.Latitude,
			Longitude: sensor.Longitude,
			Params:    sensor.Params,
		}
	}

	var wg sync.WaitGroup
	go func() {
		for sensor := range sensorChannel {
			wg.Add(1)
			go func(sensor *entity.Sensor) {
				defer wg.Done()
				client, err := NewBrokerClient(sensor.Id.Hex())
				if err != nil {
					slog.Error("Failed to setup broker client", "id", sensor.Id.Hex(), "error", err)
					os.Exit(1)
				}

				dataEmittedEvent := event.NewDataEmitted(sensor.Id.Hex())
				dataEmittedHandler, err := NewDataEmittedEventHandler(client)
				if err != nil {
					slog.Error("Failed to instantiate data emitted event handler", "id", sensor.Id.Hex(), "error", err)
					os.Exit(1)
				}
				eventDispatcher.Register(fmt.Sprintf("DataEmitted_%s", sensor.Id.Hex()), dataEmittedHandler)

				emitData, err := NewEmitDataUseCase(dataEmittedEvent, eventDispatcher)
				if err != nil {
					slog.Error("Failed to create emit data use case", "id", sensor.Id, "error", err)
					os.Exit(1)
				}

				for {
					if err := emitData.Execute(ctx, usecase.EmitDataInputDTO{
						Id: sensor.Id,
					}); err != nil {
						slog.Error("Failed to emit data", "id", sensor.Id.Hex(), "error", err)
						os.Exit(1)
					}
					time.Sleep(10 * time.Second)
				}
			}(sensor)
		}
	}()
	wg.Wait()

	sh, err := NewSensorWebHandlers(sensorChannel)
	if err != nil {
		slog.Error("Failed to instantiate sensor web handlers", "error", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/sensors", sh.CreateSensorHandler)
	r := cors.Default().Handler(mux)

	if err := http.ListenAndServe(":8081", r); err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}

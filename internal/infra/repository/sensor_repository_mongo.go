package repository

import (
	"context"

	"github.com/henriquemarlon/congo/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SensorRepository struct {
	Collection *mongo.Collection
}

func NewSensorRepository(collection *mongo.Collection) *SensorRepository {
	return &SensorRepository{
		Collection: collection,
	}
}

func (s *SensorRepository) CreateSensor(ctx context.Context, input *entity.Sensor) (*entity.Sensor, error) {
	res, err := s.Collection.InsertOne(ctx, input)
	if err != nil {
		return nil, err
	}

	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, mongo.ErrNilValue
	}

	return s.FindSensorById(ctx, id)
}

func (s *SensorRepository) FindSensorById(ctx context.Context, id primitive.ObjectID) (*entity.Sensor, error) {
	var sensor entity.Sensor
	err := s.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&sensor)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, entity.ErrSensorNotFound
		}
		return nil, err
	}
	return &sensor, nil
}

func (s *SensorRepository) FindAllSensors(ctx context.Context) ([]*entity.Sensor, error) {
	cursor, err := s.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var sensors []*entity.Sensor
	for cursor.Next(context.TODO()) {
		var sensor entity.Sensor
		if err := cursor.Decode(&sensor); err != nil {
			return nil, err
		}
		sensors = append(sensors, &sensor)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return sensors, nil
}

package sequencer

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

type Option struct {
	IsMemory bool
}

type Sequencer interface {
	Current() (int64, error)
	Save(seq int64) error
	Reset() error
}

type MongoSequence struct {
	Name  string `bson:"name,omitempty"`
	Value int64  `bson:"value,omitempty"`
}

const (
	DB_NAME       = "sample_sequence"
	COLL_NAME     = "sequences"
	SEQUENCE_NAME = "universal_sequence"
)

type MongoSequencer struct {
	client *mongo.Client
}

func NewSequencer(option Option) (Sequencer, error) {
	if option.IsMemory {
		return NewMemorySequencer()
	}
	return NewMongoSequencer()
}

func NewMongoSequencer() (*MongoSequencer, error) {
	uri := os.Getenv("MONGODB_URI")
	var err error
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return &MongoSequencer{
		client: client,
	}, nil
}

func (s *MongoSequencer) Current() (int64, error) {
	coll := s.client.Database(DB_NAME).Collection(COLL_NAME)
	// Retrieves the first matching document
	var r MongoSequence
	filter := bson.D{{"name", SEQUENCE_NAME}}
	err := coll.FindOne(context.TODO(), filter).Decode(&r)
	if errors.Is(err, mongo.ErrNoDocuments) {
		newSequence := MongoSequence{Name: SEQUENCE_NAME, Value: 0}
		_, err = coll.InsertOne(context.TODO(), newSequence)
		return 0, err
	}
	if err != nil {
		return 0, err
	}
	return r.Value, err
}

func (s *MongoSequencer) Save(seq int64) error {
	filter := bson.D{{"name", SEQUENCE_NAME}}
	update := bson.D{{"$set", bson.D{{"value", seq}}}}
	coll := s.client.Database(DB_NAME).Collection(COLL_NAME)
	_, err := coll.UpdateOne(context.TODO(), filter, update)
	return err
}

func (s *MongoSequencer) Reset() error {
	coll := s.client.Database(DB_NAME).Collection(COLL_NAME)
	filter := bson.D{{"name", SEQUENCE_NAME}}
	_, err := coll.DeleteOne(context.TODO(), filter)
	return err
}

type MemorySequencer struct {
	seq  int64
	used map[string]bool
}

func NewMemorySequencer() (*MemorySequencer, error) {
	return &MemorySequencer{}, nil
}

func (s *MemorySequencer) Current() (int64, error) {
	return s.seq, nil
}

func (s *MemorySequencer) Save(seq int64) error {
	s.seq = seq
	return nil
}

func (s *MemorySequencer) Reset() error {
	s.seq = 0
	return nil
}


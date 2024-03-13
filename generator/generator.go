package generator

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"rhzx3519/unique-key-generator/utils"
)

const (
	TotalLength = 8 // the length of key is 8
)

const (
	DbName   = "universal_unique_key"
	CollName = "universality"
)

// Generator interface
type Generator interface {
	Existed(key string) bool
	Generate() (string, error)
}

func NewGenerator(isMemory bool) (Generator, error) {
	if isMemory {
		return NewMemoryGenerator()
	}
	return NewMongoGenerator()
}

// MongoGenerator sector
type MongoGenerator struct {
	client *mongo.Client
}

type scheme struct {
	Key string `bson:"key,omitempty"`
}

func NewMongoGenerator() (*MongoGenerator, error) {
	uri := os.Getenv("MONGODB_URI")
	var err error
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return &MongoGenerator{
		client: client,
	}, nil
}

func (g *MongoGenerator) Existed(key string) bool {
	coll := g.client.Database(DbName).Collection(CollName)
	filter := bson.D{{"key", key}}
	var s scheme
	err := coll.FindOne(context.TODO(), filter).Decode(&s)
	return !errors.Is(err, mongo.ErrNoDocuments)
}

func (g *MongoGenerator) Generate() (string, error) {
	var key = utils.RandomBase64Str(TotalLength)
	for g.Existed(key) {
		key = utils.RandomBase64Str(TotalLength)
	}
	err := g.save(key)
	return key, err
}

func (g *MongoGenerator) save(key string) error {
	coll := g.client.Database(DbName).Collection(CollName)
	_, err := coll.InsertOne(context.TODO(), scheme{Key: key})
	return err
}

// MemoryGenerator sector
type MemoryGenerator struct {
	used map[string]bool
}

func NewMemoryGenerator() (*MemoryGenerator, error) {
	return &MemoryGenerator{
		used: make(map[string]bool),
	}, nil
}

func (g *MemoryGenerator) Existed(key string) bool {
	_, ok := g.used[key]
	return ok
}

func (g *MemoryGenerator) Generate() (string, error) {
	var key = utils.RandomBase64Str(TotalLength)
	for g.Existed(key) {
		key = utils.RandomBase64Str(TotalLength)
	}
	g.used[key] = true
	return key, nil
}

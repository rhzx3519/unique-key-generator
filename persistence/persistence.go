package persistence

import (
    "context"
    "errors"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "log"
    "os"
)

type Sequence struct {
    Name  string `bson:"Name,omitempty"`
    Value int64  `bson:"value,omitempty"`
}

const (
    DB_NAME       = "sample_sequence"
    COLL_NAME     = "sequences"
    SEQUENCE_NAME = "universal_sequence"
)

var client *mongo.Client

func Init() {
    uri := os.Getenv("MONGODB_URI")
    var err error
    client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
    if err != nil {
        log.Fatalln(err)
    }
}

func Destruct() {
    if err := client.Disconnect(context.TODO()); err != nil {
        log.Fatalln(err)
    }
}

func Current() (int64, error) {
    coll := client.Database(DB_NAME).Collection(COLL_NAME)
    // Retrieves the first matching document
    var r Sequence
    filter := bson.D{{"Name", SEQUENCE_NAME}}
    err := coll.FindOne(context.TODO(), filter).Decode(&r)
    if errors.Is(err, mongo.ErrNoDocuments) {
        newSequence := Sequence{Name: SEQUENCE_NAME, Value: 0}
        _, err = coll.InsertOne(context.TODO(), newSequence)
        return 0, err
    }
    if err != nil {
        return 0, err
    }
    return r.Value, err
}

func Save(seq int64) error {
    filter := bson.D{{"Name", SEQUENCE_NAME}}
    update := bson.D{{"$set", bson.D{{"value", seq}}}}
    coll := client.Database(DB_NAME).Collection(COLL_NAME)
    _, err := coll.UpdateOne(context.TODO(), filter, update)
    return err
}

func Reset() error {
    coll := client.Database(DB_NAME).Collection(COLL_NAME)
    filter := bson.D{{"Name", SEQUENCE_NAME}}
    _, err := coll.DeleteOne(context.TODO(), filter)
    return err
}

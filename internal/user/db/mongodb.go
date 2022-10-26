package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"rest_api/internal/user"
	"rest_api/pkg/logging"
)

type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func (d *db) Create(ctx context.Context, user user.User) (string, error) {
	d.logger.Debug("create user")
	result, err := d.collection.InsertOne(ctx, user)
	if err != nil {
		d.logger.Errorf("failed to create user due to error: %v.", err)
		return "", fmt.Errorf("failed to create user due to error: %v.", err)
	}

	d.logger.Debug("convert InsertedID to ObjectID")
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		d.logger.Trace(user)
		d.logger.Errorf("failed to convert objectID to hex%s", oid)
		return "", fmt.Errorf("failed to convert objectID to hex: %s", oid)
	}

	return oid.Hex(), nil
}

func (d *db) FindOne(ctx context.Context, id string) (u user.User, e error) {
	d.logger.Debug("findOne user")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		d.logger.Errorf("failed to convert hex to objectedId hex: %s due to error: %v", id, err)
		return u, fmt.Errorf("failed to convert hex to objectedId hex: %s due to error: %v", id, err)
	}

	filter := bson.M{"_id": oid}
	result := d.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		d.logger.Errorf("failed to convert hex to objectedId hex: %s due to error: %v", id, err)
		return u, fmt.Errorf("failed to convert hex to objectedId hex: %s due to error: %v", id, err)
	}

	if err = result.Decode(&u); err != nil {
		d.logger.Errorf("failed to decode user from DB error: %v", err)
		return u, fmt.Errorf("failed to decode user from DB error: %v", err)
	}

	return u, nil
}

func (d *db) Update(ctx context.Context, user user.User) error {
	panic("implement me")
}

func (d *db) Delete(ctx context.Context, id string) error {
	panic("implement me")
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) user.Storage {
	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}

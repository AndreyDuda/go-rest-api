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

func (d *db) FindAll(ctx context.Context) (users []user.User, err error) {
	errorMessage := ""
	d.logger.Debug("find all users")
	cursor, err := d.collection.Find(ctx, bson.M{})
	if cursor.Err() != nil {
		errorMessage = "failed to select users due to error: %v"
		d.logger.Errorf(errorMessage, err)
		return users, fmt.Errorf(errorMessage, err)
	}

	if err = cursor.All(ctx, &users); err != nil {
		errorMessage = "failed to read all documents from cursor. error: %v"
		d.logger.Errorf(errorMessage, err)
		return users, fmt.Errorf(errorMessage, err)
	}

	return users, err
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
		d.logger.Errorf("failed to select usre: %s due to error: %v", id, err)
		return u, fmt.Errorf("failed to select usre: %s due to error: %v", id, err)
	}

	if err = result.Decode(&u); err != nil {
		d.logger.Errorf("failed to decode user from DB error: %v", err)
		return u, fmt.Errorf("failed to decode user from DB error: %v", err)
	}

	return u, nil
}

func (d *db) Update(ctx context.Context, user user.User) error {
	var errorMessage = ""
	objectID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		errorMessage = "failed to convert user ID to objectID. ID=%s"
		d.logger.Errorf(errorMessage, user.ID)
		return fmt.Errorf(errorMessage, user.ID)
	}

	filter := bson.M{"_id": objectID}
	userBytes, err := bson.Marshal(user)
	if err != nil {
		errorMessage = "failed to marshal user. Error: %v"
		d.logger.Errorf(errorMessage, err)
		return fmt.Errorf(errorMessage, err)
	}

	var updateUserOdj bson.M
	if err = bson.Unmarshal(userBytes, &updateUserOdj); err != nil {
		errorMessage = "failed to unmarshal user bytes. Error: %v"
		d.logger.Errorf(errorMessage, err)
		return fmt.Errorf(errorMessage, err)
	}

	delete(updateUserOdj, "_id")

	update := bson.M{
		"$set": updateUserOdj,
	}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		errorMessage = "failed to executu update user query. Error: %v"
		d.logger.Errorf(errorMessage, err)
		return fmt.Errorf(errorMessage, err)
	}

	if result.MatchedCount == 0 {
		errorMessage = "not found user for update. ID=%s"
		d.logger.Errorf(errorMessage, user.ID)
		return fmt.Errorf(errorMessage, user.ID)
	}

	d.logger.Trace("Matched %d, documents and modified %d documents", result.MatchedCount, result.ModifiedCount)

	return nil
}

func (d *db) Delete(ctx context.Context, id string) error {
	var errorMessage = ""
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		errorMessage = "failed to convert user ID to objectID. ID=%s"
		d.logger.Errorf(errorMessage, id)
		return fmt.Errorf(errorMessage, id)
	}

	filter := bson.M{"_id": objectID}
	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		errorMessage = "failed to execute query. error: %v"
		d.logger.Errorf(errorMessage, err)
		return fmt.Errorf(errorMessage, err)
	}

	if result.DeletedCount == 0 {
		errorMessage = "not found user for update. ID=%s"
		d.logger.Errorf(errorMessage, id)
		return fmt.Errorf(errorMessage, id)
	}

	d.logger.Trace("Deleted %d, documents", result.DeletedCount)

	return nil
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) user.Storage {
	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}

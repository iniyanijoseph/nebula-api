package dao

import (
	"context"
	"fmt"

	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Filter interface {
	ToDocument() ([]byte, error)
	GetOffset() int64
}

type collectionHelper[F Filter, M any] struct {
	coll     *mongo.Collection
	pageSize int64
}

func newCollectionHelper[F Filter, M any](coll *mongo.Collection, pageSize int64) *collectionHelper[F, M] {
	return &collectionHelper[F, M]{coll: coll}
}

func IsUnknownFieldErr(err error) bool {
	multiErr, ok := errors.Cause(err).(schema.MultiError)
	if !ok {
		return false
	}

	for _, v := range multiErr {
		if _, ok := v.(schema.UnknownKeyError); ok {
			return true
		}
	}

	return false
}

func (coll *collectionHelper[F, M]) Filter(ctx context.Context, filter F) ([]M, error) {
	query, err := filter.ToDocument()
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse queries into filter")
	}

	limitOpt := options.Find().SetSkip(filter.GetOffset()).SetLimit(coll.pageSize)
	cursor, err := coll.coll.Find(ctx, query, limitOpt)
	if err != nil {
		return nil, fmt.Errorf("failed to query courses: %v", err)
	}

	var models []M
	err = cursor.All(ctx, &models)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal courses from query: %v", err)
	}

	return models, nil
}

func (coll *collectionHelper[F, M]) FindById(ctx context.Context, id string) (*M, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("failed to parse object id: %v", id)
	}

	var model M
	err = coll.coll.FindOne(ctx, bson.M{"_id": objId}).Decode(&model)
	if err != nil {
		return nil, fmt.Errorf("failed to find and decode course by id=%v : %v", objId, err)
	}

	return &model, nil
}

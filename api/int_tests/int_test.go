package inttest_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/ory/dockertest"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	resource, err := pool.Run("mongo", "latest")
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	client, err := mongo.NewClient(
		options.Client().ApplyURI(
			fmt.Sprintf("mongodb://localhost:%s/", resource.GetPort("27017/tcp"))))
	if err != nil {
		return err
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = client.Connect(ctx)
		if err != nil {
			return err
		}

		return client.Ping(ctx, nil)
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}
}

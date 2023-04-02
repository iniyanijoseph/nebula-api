package inttest_test

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMain(m *testing.M) {
	pool := createDockerPool()
	resource := createMongoResource(pool)
	defer purgeResourceFromPool(pool, resource)
	waitForMongoResource(pool, resource)
	importData(resource)
}

func createDockerPool() *dockertest.Pool {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("could not connect to Docker: %s", err)
	}

	return pool
}

func createMongoResource(pool *dockertest.Pool) *dockertest.Resource {
	resource, err := pool.Run("mongo", "latest", nil)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	return resource
}

func purgeResourceFromPool(pool *dockertest.Pool, resource *dockertest.Resource) {
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}

func waitForMongoResource(pool *dockertest.Pool, resource *dockertest.Resource) {
	log.Println("waiting for mongo database to come up")
	err := pool.Retry(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		client, err := mongo.Connect(ctx,
			options.Client().ApplyURI(fmt.Sprintf("mongodb://localhost:%s", resource.GetPort("27017/tcp"))))
		if err != nil {
			return err
		}

		return client.Ping(ctx, nil)
	})
	if err != nil {
		log.Fatalf("could not connect to local mongo database: %s", err)
	}

	log.Println("successfully connected to mongodb")
}

func importData(resource *dockertest.Resource) {
	colls := [4]string{"courses", "exams", "professors", "sections"}
	for _, coll := range colls {

		f, err := os.Open("./data/" + coll + ".json")
		if err != nil {
			log.Fatalf("failed to open collection file for %s: %s", coll, err)
		}
		defer f.Close()

		var stdout bytes.Buffer
		var stderr bytes.Buffer

		log.Printf("beginning import of %s collection", coll)

		code, err := resource.Exec(
			[]string{
				"mongoimport",
				"--uri",
				"mongodb://localhost:27017/combinedDB",
				"--collection",
				coll,
				"--type",
				"json",
				"--jsonArray",
			},
			dockertest.ExecOptions{StdIn: f, StdOut: &stdout, StdErr: &stderr},
		)

		if err != nil || code != 0 {
			log.Printf("failed to execute mongoimport on collection %s: %s", coll, err)
			log.Printf("exitcode = %d", code)
			log.Println("stderr:")
			log.Println(stderr.String())
			log.Println("stdout:")
			log.Fatalln(stdout.String())
		}

		log.Printf("successfully imported the %s collection", coll)
	}
}

package inttest_test

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMain(m *testing.M) {
	pool, resource := setup()
	defer teardown(pool, resource)
}

func setup() (*dockertest.Pool, *dockertest.Resource) {
	pool := createDockerPool()
	resource := createMongoResource(pool)
	waitForMongoResource(pool, resource)
	importData(resource)

	return pool, resource
}

func teardown(pool *dockertest.Pool, resource *dockertest.Resource) {
	purgeResourceFromPool(pool, resource)
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

		stdOutOut, stdOutIn := io.Pipe()
		stdErrOut, stdErrIn := io.Pipe()

		var wg sync.WaitGroup
		wg.Add(2)

		defer func() {
			f.Close()
			stdOutIn.Close()
			stdErrIn.Close()
			wg.Wait()
		}()

		go monitorPipe("mongoimport", stdOutOut, &wg)
		go monitorPipe("mongoimport", stdErrOut, &wg)

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
			dockertest.ExecOptions{StdIn: f, StdOut: stdOutIn, StdErr: stdErrIn},
		)

		if err != nil {
			log.Fatalf("failed to execute mongoimport on collection %s: %s", coll, err)
		}

		if code != 0 {
			log.Fatalf("mongoimport returned non-zero exit code on collection %s", coll)
		}

		log.Printf("successfully imported the %s collection", coll)
	}
}

func monitorPipe(prefix string, input io.Reader, wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		fmt.Printf("[%s] %s\n", prefix, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("scanner error: %s", err)
	}
}

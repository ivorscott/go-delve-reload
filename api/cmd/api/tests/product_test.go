package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/ivorscott/go-delve-reload/cmd/api/internal/handlers"
	"github.com/ivorscott/go-delve-reload/internal/schema"
	_ "github.com/lib/pq"

	"github.com/google/go-cmp/cmp"
	"github.com/ivorscott/go-delve-reload/internal/platform/conf"
	"github.com/ivorscott/go-delve-reload/internal/platform/database"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var cfg struct {
	Web struct {
		Address         string        `conf:"default:localhost:4000"`
		Production      bool          `conf:"default:false"`
		ReadTimeout     time.Duration `conf:"default:5s"`
		WriteTimeout    time.Duration `conf:"default:5s"`
		ShutdownTimeout time.Duration `conf:"default:5s"`
		FrontendAddress string        `conf:"default:https://localhost:3000"`
	}
	DB struct {
		User       string `conf:"default:postgres"`
		Password   string `conf:"default:postgres,noprint"`
		Host       string `conf:"default:localhost"`
		Name       string `conf:"default:postgres"`
		DisableTLS bool   `conf:"default:true"`
	}
}

func TestProducts(t *testing.T) {
	log := log.New(os.Stderr, "TEST : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	log.Println("Starting postgres container...")

	if err := conf.Parse(os.Args[1:], "API", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("API", &cfg)
			if err != nil {
				log.Fatal(err, "generating usage")
			}
			fmt.Println(usage)
		}
		log.Fatal(err, "error: parsing config")
	}

	// Background returns a non-nil, empty Context. It is never canceled, has no values, and has no deadline.
	// It is typically used by the main function, initialization, and tests, and as the top-level Context for incoming requests.
	ctx := context.Background()

	// Port is a string containing port number and protocol in the format "80/tcp"
	postgresPort := nat.Port("5432/tcp")

	// ContainerRequest represents the parameters used to get a running container
	req := tc.ContainerRequest{
		Image:        "postgres",
		ExposedPorts: []string{postgresPort.Port()},
		Env: map[string]string{
			"POSTGRES_PASSWORD": cfg.DB.Password,
			"POSTGRES_USER":     cfg.DB.User,
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections"),
			wait.ForListeningPort(postgresPort),
		),
	}

	// GenericContainer creates a generic container with parameters
	postgres, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true, // auto-start the container
	})
	if err != nil {
		t.Fatal("start:", err)
	}

	// MappedPort gets the externally mapped port for the container
	hostPort, err := postgres.MappedPort(ctx, postgresPort)
	if err != nil {
		log.Fatal("map:", err)
	}

	repo, err := database.NewRepository(database.Config{
		User:       cfg.DB.User,
		Host:       cfg.DB.Host + ":" + hostPort.Port(),
		Name:       cfg.DB.Name,
		Password:   cfg.DB.Password,
		DisableTLS: cfg.DB.DisableTLS,
	})
	if err != nil {
		t.Fatal(err, "connecting to db")
	}
	defer repo.Close()

	log.Printf("Postgres container started, running at:  %s\n", repo.URL.String())

	if err := schema.Migrate("postgres", repo.URL.String()); err != nil {
		t.Fatal(err)
	}

	if err := schema.Seed(repo.DB, "products"); err != nil {
		t.Fatal(err)
	}

	shutdown := make(chan os.Signal, 1)

	// create application handler
	tests := ProductTests{app: handlers.API(shutdown, repo, log, cfg.Web.FrontendAddress)}

	// test handlers
	t.Run("List", tests.List)
	t.Run("ProductCRUD", tests.ProductCRUD)
}

// ProductTests holds methods for each product subtest. This type allows
// passing dependencies for tests while still providing a convenient syntax
// when subtests are registered.
type ProductTests struct {
	app http.Handler
}

func (p *ProductTests) List(t *testing.T) {
	req := httptest.NewRequest("GET", "/v1/products", nil)
	resp := httptest.NewRecorder()

	p.app.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("getting: expected status code %v, got %v", http.StatusOK, resp.Code)
	}

	var list []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		t.Fatalf("decoding: %s", err)
	}

	want := []map[string]interface{}{
		{
			"id":          "cbef5139-323f-48b8-b911-dc9be7d0bc07",
			"name":        "Xbox One X",
			"price":       float64(499),
			"description": "Eighth-generation home video game console developed by Microsoft.",
			"created":     "2019-01-01T00:00:01.000001Z",
			"tags":        nil,
		},
		{
			"id":          "ce93a886-3a0e-456b-b7f5-8652d2de1e8f",
			"name":        "Playstation 4",
			"price":       float64(299),
			"description": "Eighth-generation home video game console developed by Sony Interactive Entertainment.",
			"created":     "2019-01-01T00:00:01.000001Z",
			"tags":        nil,
		},
		{
			"id":          "faa25b57-7031-4b37-8a89-de013418deb0",
			"name":        "Nintendo Switch",
			"price":       float64(299),
			"description": "Hybrid console that can be used as a stationary and portable device developed by Nintendo.",
			"created":     "2019-01-01T00:00:01.000001Z",
			"tags":        nil,
		},
	}

	if diff := cmp.Diff(want, list); diff != "" {
		t.Fatalf("Response did not match expected. Diff:\n%s", diff)
	}
}

func (p *ProductTests) ProductCRUD(t *testing.T) {
	var created map[string]interface{}

	{ // CREATE
		body := strings.NewReader(`{"name":"Some product","price":750,"description":"some description"}`)

		req := httptest.NewRequest("POST", "/v1/products", body)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		p.app.ServeHTTP(resp, req)

		if http.StatusCreated != resp.Code {
			t.Fatalf("posting: expected status code %v, got %v", http.StatusCreated, resp.Code)
		}

		if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
			t.Fatalf("decoding: %s", err)
		}

		if created["id"] == "" || created["id"] == nil {
			t.Fatal("expected non-empty product id")
		}
		if created["created"] == "" || created["created"] == nil {
			t.Fatal("expected non-empty product created")
		}

		want := map[string]interface{}{
			"id":          created["id"],
			"created":     created["created"],
			"name":        "Some product",
			"price":       float64(750),
			"description": "some description",
			"tags":        nil,
		}

		if diff := cmp.Diff(want, created); diff != "" {
			t.Fatalf("Response did not match expected. Diff:\n%s", diff)
		}
	}

	{ // READ
		url := fmt.Sprintf("/v1/products/%s", created["id"])
		req := httptest.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		p.app.ServeHTTP(resp, req)

		if http.StatusOK != resp.Code {
			t.Fatalf("retrieving: expected status code %v, got %v", http.StatusOK, resp.Code)
		}

		var fetched map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&fetched); err != nil {
			t.Fatalf("decoding: %s", err)
		}

		// Fetched product should match the one we created.
		if diff := cmp.Diff(created, fetched); diff != "" {
			t.Fatalf("Retrieved product should match created. Diff:\n%s", diff)
		}
	}
}

package api

import (
	"encoding/json"
	"github.com/ivorscott/go-delve-reload/internal/models"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestProductsHandler(t *testing.T) {

	products := []models.Product{}

	resp, err := http.Get("http://localhost:4000/products")
	if err != nil {
		t.Fatalf("GET /products err= %s; want nil", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("ioutils.ReadAll() err =%s; want nil", err)
	}

	err = json.Unmarshal(body, &products)
	if err != nil {
		t.Fatalf("json.Unmarshal() err =%s; want nil", err)
	}

	got := len(products)
	want := 3

	assert.Equal(t, got, want, "Response body length differs")
}

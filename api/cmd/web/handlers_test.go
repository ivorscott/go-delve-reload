package main

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductsHandler(t *testing.T) {
	resp, err := http.Get("http://localhost:4000/products")
	if err != nil {
		t.Fatalf("GET /products err= %s; want nil", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutils.ReadAll() err =%s; want nil", err)
	}

	got := string(body)
	want := `
	[{
		"ID":1,
		"Name":"Xbox One X",
		"Price":499,
		"Description":"Eighth-generation home video game console developed by Microsoft.",
		"Created":"2020-01-02T14:02:26.58977Z"
	},
	{
		"ID":2,
		"Name":"Playsation 4",
		"Price":299,
		"Description":"Eighth-generation home video game console developed by Sony Interactive Entertainment.",
		"Created":"2020-01-02T14:02:26.58977Z"
	},
	{
		"ID":3,
		"Name":"Nintendo Switch",
		"Price":299,
		"Description":"Hybrid console that can be used as a stationary and portable device developed by Nintendo.",
		"Created":"2020-01-02T14:02:26.58977Z"
	}]`

	assert.JSONEq(t, got, want, "Response body differs")
}

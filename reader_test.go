package cherry

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

type SampleBody struct {
	Title string `json:"title"`
}

type SampleRequest struct {
	ID         int64       `param:"id"`
	Name       int         `param:"name"`
	Price      int         `param:"price"`
	SampleBody *SampleBody `body:"json"`
}

func TestReader_Read(t *testing.T) {
	ctx := context.Background()
	bodyBytes, _ := json.Marshal(SampleBody{Title: "Chamith Udayange"})
	req, _ := http.NewRequest("GET", "http://localhost:8080/user/123?name=100", bytes.NewBuffer(bodyBytes))
	req.Header.Add("content-type", "application/json")

	s := SampleRequest{
		ID:         1,
		Name:       0,
		Price:      0,
		SampleBody: new(SampleBody),
	}
	if err := NewReader().Read(ctx, req, &s); err != nil {
		t.Error(err.Error())
	}
	t.Log(s)
	t.Log(s.SampleBody.Title)
}

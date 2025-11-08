package opensearch

import (
	"os"
	"testing"
)

func initTestClient() *Clients {
	clients, err := NewClients(&Config{
		Host:     "127.0.0.1",
		Port:     9200,
		UserName: os.Getenv("OPENSEARCH_USERNAME"),
		Password: os.Getenv("OPENSEARCH_PASSWORD"),
	})
	if err != nil {
		panic(err)
	}
	return clients
}

func Test_CatIndices(t *testing.T) {
	clients := initTestClient()
	oc := clients.GetOpenSearchClient()
	indices, err := oc.Cat.Indices(oc.Cat.Indices.WithFormat("json"))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(indices)
}

package opensearch

import (
	"context"
	"testing"
)

func TestClients_CatIndices(t *testing.T) {
	clients := initTestClient()
	result, _, err := clients.CatIndices(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(result)
}

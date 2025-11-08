package opensearch

import "sync"

var (
	logClusters     = map[string]*Clients{}
	logClustersLock = &sync.RWMutex{}
)

func GetCluster(cluster string) *Clients {
	logClustersLock.RLock()
	defer logClustersLock.RUnlock()
	if client, ok := logClusters[cluster]; ok {
		return client
	}
	return nil
}

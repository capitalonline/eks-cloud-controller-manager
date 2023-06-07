package provider

import "context"

type Clusters struct {
}

func (c *Clusters) ListClusters(ctx context.Context) ([]string, error) {
	return nil, nil
}

func (c *Clusters) Master(ctx context.Context, clusterName string) (string, error) {
	return "", nil
}

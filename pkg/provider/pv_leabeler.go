package provider

import (
	"context"
	v1 "k8s.io/api/core/v1"
)

type PVLabeler struct {
}

func (p *PVLabeler) GetLabelsForVolume(ctx context.Context, pv *v1.PersistentVolume) (map[string]string, error) {
	return nil, nil
}

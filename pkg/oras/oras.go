package oras

import (
	"fmt"
	"strings"

	"oras.land/oras-go/v2"
)

type Interface interface {
	Attach(reference string, artifact string, artifactType string, annotations ...string) ([]string, error)
}

type client struct{}

var _ Interface = &client{}

func NewClient() Interface {
	return &client{}
}

func (c *client) Attach(reference string, artifact string, artifactType string, annotations ...string) ([]string, error) {
	manifestAnnotations := make(map[string]string)
	for _, anno := range annotations {
		key, val, success := strings.Cut(anno, "=")
		if !success {
			return nil, fmt.Errorf("%w: %s", "invalid annotation format", anno)
		}
		if _, ok := manifestAnnotations[key]; ok {
			return nil, fmt.Errorf("%w: %v, ", "duplicate annotation", key)
		}
		manifestAnnotations[key] = val
	}

	packOpts := oras.PackOptions{
		Subject:             &subject,
		ManifestAnnotations: manifestAnnotations[option.AnnotationManifest],
	}
	return nil, nil
}

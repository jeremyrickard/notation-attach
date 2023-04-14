package notation

import "github.com/jeremyrickard/notation-attach/pkg/oras"

type Interface interface {
	AttachSignature(artifact string, file string) (string, error)
}

type client struct {
	oras oras.Interface
}

var _ Interface = &client{}

func NewClient(oras oras.Interface) Interface {
	return &client{
		oras: oras,
	}
}

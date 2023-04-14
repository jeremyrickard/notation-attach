package notation

import (
	"encoding/json"
	"fmt"
)

func (c *client) AttachSignature(artifact string, signaturePath string) (string, error) {

	var annotations []string
	thumbprints, err := getThumbprints(signaturePath)
	if err != nil {
		return "", fmt.Errorf("unable to get thumbprints: %w", err)
	}
	thumbprintS256, err := json.Marshal(thumbprints)
	if err != nil {
		return "", fmt.Errorf("unable to marshall thumbprints: %w", err)
	}
	annotations = []string{
		fmt.Sprintf(thumbprintAnnotationFormatString, string(thumbprintS256)),
	}
	c.oras.Attach(artifact, signaturePath, "", annotations...)
	return "", nil
}

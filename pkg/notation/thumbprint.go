package notation

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/notaryproject/notation-core-go/signature/cose"
)

// getThumbprints loads the specified notary compatible signature and
// extracts the thumbprints for the specified certificate chain,
// which are then returned as a slice. If unable to read the signature
// or the signature is malformed, an error is returned.
func getThumbprints(file string) ([]string, error) {
	sigBytes, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("unable to read %s: %w", file, err)
	}
	sig, err := cose.ParseEnvelope(sigBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse signature envelope for %s: %s", file, err)
	}
	content, err := sig.Content()
	if err != nil {
		return nil, fmt.Errorf("unable to get signature content for %s: %s", file, err)
	}

	var thumbprints []string
	for _, cert := range content.SignerInfo.CertificateChain {
		thumbprint := sha256.Sum256(cert.Raw)
		thumbprints = append(thumbprints, hex.EncodeToString(thumbprint[:]))
	}

	return thumbprints, nil
}

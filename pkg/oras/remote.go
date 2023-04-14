package oras

import (
	"context"
	"fmt"
	"net"
	"net/http"

	dockerAuth "oras.land/oras-go/pkg/auth/docker"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
)

type Remote struct {
	CACertFilePath    string
	PlainHTTP         bool
	Insecure          bool
	Configs           []string
	Username          string
	PasswordFromStdin bool
	Password          string

	resolveFlag           []string
	applyDistributionSpec bool
	distributionSpec      distributionSpec
	headerFlags           []string
	headers               http.Header
}

func (opts *Remote) NewRepository(reference string, debug bool) (repo *remote.Repository, err error) {
	repo, err = remote.NewRepository(reference)
	if err != nil {
		return nil, err
	}
	hostname := repo.Reference.Registry
	repo.PlainHTTP = opts.isPlainHttp(hostname)
	client := auth.DefaultClient
	client.Credential = getDockerClientAuth
	repo.Client = client
	if opts.distributionSpec.referrersAPI != nil {
		repo.SetReferrersCapability(*opts.distributionSpec.referrersAPI)
	}
	return
}

// distributionSpec option struct.
type distributionSpec struct {
	// referrersAPI indicates the preference of the implementation of the Referrers API.
	// Set to true for referrers API, false for referrers tag scheme, and nil for auto fallback.
	referrersAPI *bool

	// specFlag should be provided in form of`<version>-<api>-<option>`
	specFlag string
}

// Parse parses flags into the option.
func (opts *distributionSpec) Parse() error {
	switch opts.specFlag {
	case "":
		opts.referrersAPI = nil
	case "v1.1-referrers-tag":
		isApi := false
		opts.referrersAPI = &isApi
	case "v1.1-referrers-api":
		isApi := true
		opts.referrersAPI = &isApi
	default:
		return fmt.Errorf("unknown image specification flag: %q", opts.specFlag)
	}
	return nil
}

// isPlainHttp returns the plain http flag for a given registry.
func (opts *Remote) isPlainHttp(registry string) bool {
	host, _, _ := net.SplitHostPort(registry)
	if host == "localhost" || registry == "localhost" {
		return true
	}
	return opts.PlainHTTP
}

// getDockerClientAuth implements the interface required by v2 of the oras-go Client for interacting
// with registries and uses Docker auth to resolve the credential.
func getDockerClientAuth(_ context.Context, regString string) (auth.Credential, error) {
	client, err := dockerAuth.NewClient()
	if err != nil {
		return auth.Credential{}, err
	}

	dockerAuthClient, ok := client.(*dockerAuth.Client)
	if !ok {
		return auth.Credential{}, fmt.Errorf("could not get docker auth client")
	}
	user, pass, err := dockerAuthClient.Credential(regString)
	if err != nil {
		return auth.Credential{}, err
	}
	// Credential will either return a username/password OR an identity token.
	// if user is empty, set the user/pass in the auth.Credential. Otherwise
	// set the identity token
	if user != "" {
		ac := auth.Credential{
			Username: user,
			Password: pass,
		}
		return ac, nil
	}
	ac := auth.Credential{
		RefreshToken: pass,
	}
	return ac, nil
}

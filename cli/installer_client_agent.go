//go:build client || agent
// +build client agent

package cli

import (
	gofig "github.com/akutz/gofig/types"
	apitypes "github.com/nooperpudd/rexray/libstorage/api/types"
)

func installSelfCert(ctx apitypes.Context, config gofig.Config) error {
	return nil
}

//go:build !client && !agent && !controller
// +build !client,!agent,!controller

package util

import (
	gofig "github.com/akutz/gofig/types"

	apitypes "github.com/nooperpudd/rexray/libstorage/api/types"
	apiclient "github.com/nooperpudd/rexray/libstorage/client"
)

func newClient(
	ctx apitypes.Context, config gofig.Config) (apitypes.Client, error) {
	return apiclient.New(ctx, config)
}

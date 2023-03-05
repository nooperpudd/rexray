//go:build go1.7 && !agent && !controller
// +build go1.7,!agent,!controller

package scripts

import (
	"net/http"

	apitypes "github.com/nooperpudd/rexray/libstorage/api/types"
)

func doRequest(
	ctx apitypes.Context,
	req *http.Request) (*http.Response, error) {
	return doRequestWithClient(ctx, http.DefaultClient, req)
}

func doRequestWithClient(
	ctx apitypes.Context,
	client *http.Client,
	req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	return client.Do(req)
}

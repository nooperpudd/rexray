//go:build !client && !controller
// +build !client,!controller

package cli

import "github.com/nooperpudd/rexray/agent"

func init() {
	startFuncs = append(startFuncs, agent.Start)
}

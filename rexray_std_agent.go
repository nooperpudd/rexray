//go:build !client && !controller
// +build !client,!controller

package main

import (
	// Load the agent's modules
	_ "github.com/nooperpudd/rexray/agent/csi"

	// Load the in-tree CSI plug-ins
	_ "github.com/nooperpudd/rexray/agent/csi/libstorage"

	// Load vendored CSI plug-ins
	_ "github.com/rexray/csi-blockdevices/provider"
	//_ "github.com/rexray/csi-nfs/provider"
	//_ "github.com/rexray/csi-vfs/provider"
)

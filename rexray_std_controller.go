//go:build !client && !agent
// +build !client,!agent

package main

import (
	// load the libstorage packages
	_ "github.com/nooperpudd/rexray/libstorage/imports/storage"
)

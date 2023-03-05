//go:build !agent && !controller
// +build !agent,!controller

/*
Package scripts includes scripts that are embedded in REX-Ray during the
build process that is driven by the make file.
*/
package scripts

import (
	// depend upon this tool with a nil import in order to preserve it
	// in the dependency list
	_ "github.com/go-bindata/go-bindata"
)

package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"time"

	// load the golf package
	_ "github.com/akutz/golf"

	"github.com/nooperpudd/rexray/libstorage/api/context"
	"github.com/nooperpudd/rexray/libstorage/api/types"
)

// GetTypePkgPathAndName gets ths type and package path of the provided
// instance.
func GetTypePkgPathAndName(i interface{}) string {
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface {
		t = t.Elem()
	}
	pkgPath := t.PkgPath()
	typeName := t.Name()
	if pkgPath == "" {
		return typeName
	}
	return fmt.Sprintf("%s.%s", pkgPath, typeName)
}

// GetTempSockFile returns a new sock file in a temp space.
func GetTempSockFile(ctx types.Context) string {

	f, err := ioutil.TempFile(context.MustPathConfig(ctx).Run, "")
	if err != nil {
		panic(err)
	}
	name := f.Name()
	os.RemoveAll(name)
	return fmt.Sprintf("%s.sock", name)
}

// DeviceAttachTimeout gets the configured device attach timeout.
func DeviceAttachTimeout(val string) time.Duration {
	dur, err := time.ParseDuration(val)
	if err != nil {
		return time.Duration(30) * time.Second
	}
	return dur
}

//go:build darwin
// +build darwin

package darwin

import (
	"runtime"

	gofig "github.com/akutz/gofig/types"
	"github.com/akutz/goof"

	"github.com/nooperpudd/rexray/libstorage/api/registry"
	"github.com/nooperpudd/rexray/libstorage/api/types"
)

const driverName = "darwin"

var (
	errUnknownOS             = goof.New("unknown OS")
	errUnknownFileSystem     = goof.New("unknown file system")
	errUnsupportedFileSystem = goof.New("unsupported file system")
)

func init() {
	registry.RegisterOSDriver(driverName, newDriver)
}

type driver struct {
	config gofig.Config
}

func newDriver() types.OSDriver {
	return &driver{}
}

func (d *driver) Init(ctx types.Context, config gofig.Config) error {
	if runtime.GOOS != "darwin" {
		return errUnknownOS
	}
	d.config = config
	return nil
}

func (d *driver) Name() string {
	return driverName
}

func (d *driver) Mounts(
	ctx types.Context,
	deviceName, mountPoint string,
	opts types.Store) ([]*types.MountInfo, error) {

	return nil, nil
}

func (d *driver) Mount(
	ctx types.Context,
	deviceName, mountPoint string,
	opts *types.DeviceMountOpts) error {

	return nil
}

func (d *driver) Unmount(
	ctx types.Context,
	mountPoint string,
	opts types.Store) error {

	return nil
}

func (d *driver) IsMounted(
	ctx types.Context,
	mountPoint string,
	opts types.Store) (bool, error) {

	return false, nil
}

func (d *driver) Format(
	ctx types.Context,
	deviceName string,
	opts *types.DeviceFormatOpts) error {

	return nil
}

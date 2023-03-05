package types

import (
	"fmt"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/stretchr/testify/assert"
)

func newLocalDevicesObj() *LocalDevices {
	return &LocalDevices{
		Driver: "vfs",
		DeviceMap: map[string]string{
			"vfs-000": "/dev/xvda",
			"vfs-001": "/dev/xvdb",
			"vfs-002": "/dev/xvdc",
		},
	}
}

var expectedLD1String = "vfs=vfs-000::/dev/xvda,vfs-001::/dev/xvdb,vfs-002::/dev/xvdc"

func TestLocalDevicesMarshalText(t *testing.T) {

	ld1 := newLocalDevicesObj()
	assert.Equal(t, expectedLD1String, ld1.String())
	t.Logf("localDevices=%s", ld1)

	ld2 := &LocalDevices{}
	assert.NoError(t, ld2.UnmarshalText([]byte(ld1.String())))
	assert.EqualValues(t, ld1, ld2)
}

func TestLocalDevicesUnmarshalText(t *testing.T) {

	ld1 := &LocalDevices{}
	err := ld1.UnmarshalText([]byte("scaleio="))
	assert.NoError(t, err)
}

func TestLocalDevicesMarshalJSON(t *testing.T) {

	ld1 := newLocalDevicesObj()

	buf, err := ld1.MarshalJSON()
	assert.NoError(t, err)
	t.Logf("localDevices=%s", string(buf))

	ld2 := &LocalDevices{}
	assert.NoError(t, ld2.UnmarshalJSON(buf))

	assert.EqualValues(t, ld1, ld2)
}

func TestLocalDevicesMarshalToYAML(t *testing.T) {
	out, err := yaml.Marshal(newLocalDevicesObj())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(out))
}

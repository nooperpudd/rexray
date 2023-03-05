package libstorage

import (
	"github.com/nooperpudd/rexray/libstorage/api/types"
	"github.com/nooperpudd/rexray/libstorage/api/utils"
)

func (c *client) Instances(
	ctx types.Context) (map[string]*types.Instance, error) {

	if c.isController() {
		return nil, utils.NewUnsupportedForClientTypeError(
			c.clientType, "Instances")
	}

	ctx = c.withAllInstanceIDs(c.requireCtx(ctx))
	return c.APIClient.Instances(ctx)
}

func (c *client) InstanceInspect(
	ctx types.Context, service string) (*types.Instance, error) {

	if c.isController() {
		return nil, utils.NewUnsupportedForClientTypeError(
			c.clientType, "InstanceInspect")
	}

	ctx = c.withInstanceID(c.requireCtx(ctx), service)
	i, err := c.APIClient.InstanceInspect(ctx, service)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (c *client) Root(
	ctx types.Context) ([]string, error) {

	return c.APIClient.Root(c.requireCtx(ctx))
}

func (c *client) Services(
	ctx types.Context) (map[string]*types.ServiceInfo, error) {

	ctx = c.withAllInstanceIDs(c.requireCtx(ctx))

	svcInfo, err := c.APIClient.Services(ctx)
	if err != nil {
		return nil, err
	}
	for k, v := range svcInfo {
		c.serviceCache.Set(k, v)
	}
	return svcInfo, err
}

func (c *client) ServiceInspect(
	ctx types.Context, service string) (*types.ServiceInfo, error) {

	ctx = c.withInstanceID(c.requireCtx(ctx), service)
	return c.APIClient.ServiceInspect(ctx, service)
}

func (c *client) Volumes(
	ctx types.Context,
	attachments types.VolumeAttachmentsTypes) (types.ServiceVolumeMap, error) {

	ctx = c.requireCtx(ctx)

	ctxA, err := c.withAllLocalDevices(ctx)
	if err != nil {
		return nil, err
	}
	ctx = c.withAllInstanceIDs(ctxA)

	return c.APIClient.Volumes(ctx, attachments)
}

func (c *client) VolumesByService(
	ctx types.Context,
	service string,
	attachments types.VolumeAttachmentsTypes) (types.VolumeMap, error) {

	ctx = c.withInstanceID(c.requireCtx(ctx), service)
	ctxA, err := c.withAllLocalDevices(ctx)
	if err != nil {
		return nil, err
	}
	ctx = ctxA

	return c.APIClient.VolumesByService(ctx, service, attachments)
}

func (c *client) VolumeInspect(
	ctx types.Context,
	service, volumeID string,
	attachments types.VolumeAttachmentsTypes) (*types.Volume, error) {

	ctx = c.withInstanceID(c.requireCtx(ctx), service)
	ctxA, err := c.withAllLocalDevices(ctx)
	if err != nil {
		return nil, err
	}
	ctx = ctxA

	return c.APIClient.VolumeInspect(ctx, service, volumeID, attachments)
}

func (c *client) VolumeInspectByName(
	ctx types.Context,
	service, volumeName string,
	attachments types.VolumeAttachmentsTypes) (*types.Volume, error) {

	ctx = c.withInstanceID(c.requireCtx(ctx), service)
	ctxA, err := c.withAllLocalDevices(ctx)
	if err != nil {
		return nil, err
	}
	ctx = ctxA

	return c.APIClient.VolumeInspectByName(
		ctx, service, volumeName, attachments)
}

func (c *client) VolumeCreate(
	ctx types.Context,
	service string,
	request *types.VolumeCreateRequest) (*types.Volume, error) {

	ctx = c.withInstanceID(c.requireCtx(ctx), service)
	ctxA, err := c.withAllLocalDevices(ctx)
	if err != nil {
		return nil, err
	}
	ctx = ctxA

	vol, err := c.APIClient.VolumeCreate(ctx, service, request)
	if err != nil {
		return nil, err
	}

	return vol, nil
}

func (c *client) VolumeCreateFromSnapshot(
	ctx types.Context,
	service, snapshotID string,
	request *types.VolumeCreateRequest) (*types.Volume, error) {

	ctx = c.withInstanceID(c.requireCtx(ctx), service)

	vol, err := c.APIClient.VolumeCreateFromSnapshot(
		ctx, service, snapshotID, request)
	if err != nil {
		return nil, err
	}

	return vol, nil
}

func (c *client) VolumeCopy(
	ctx types.Context,
	service, volumeID string,
	request *types.VolumeCopyRequest) (*types.Volume, error) {

	ctx = c.withInstanceID(c.requireCtx(ctx), service)

	vol, err := c.APIClient.VolumeCopy(ctx, service, volumeID, request)
	if err != nil {
		return nil, err
	}

	return vol, nil
}

func (c *client) VolumeRemove(
	ctx types.Context,
	service, volumeID string,
	force bool) error {

	ctx = c.withInstanceID(c.requireCtx(ctx), service)

	err := c.APIClient.VolumeRemove(ctx, service, volumeID, force)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) VolumeAttach(
	ctx types.Context,
	service string,
	volumeID string,
	request *types.VolumeAttachRequest) (*types.Volume, string, error) {

	if c.isController() {
		return nil, "", utils.NewUnsupportedForClientTypeError(
			c.clientType, "VolumeAttach")
	}

	ctx = c.withInstanceID(c.requireCtx(ctx), service)
	ctxA, err := c.withAllLocalDevices(ctx)
	if err != nil {
		return nil, "", err
	}
	ctx = ctxA

	return c.APIClient.VolumeAttach(ctx, service, volumeID, request)
}

func (c *client) VolumeDetach(
	ctx types.Context,
	service string,
	volumeID string,
	request *types.VolumeDetachRequest) (*types.Volume, error) {

	if c.isController() {
		return nil, utils.NewUnsupportedForClientTypeError(
			c.clientType, "VolumeDetach")
	}

	ctx = c.withInstanceID(c.requireCtx(ctx), service)
	ctxA, err := c.withAllLocalDevices(ctx)
	if err != nil {
		return nil, err
	}
	ctx = ctxA

	return c.APIClient.VolumeDetach(ctx, service, volumeID, request)
}

func (c *client) VolumeDetachAll(
	ctx types.Context,
	request *types.VolumeDetachRequest) (types.ServiceVolumeMap, error) {

	if c.isController() {
		return nil, utils.NewUnsupportedForClientTypeError(
			c.clientType, "VolumeDetachAll")
	}

	ctx = c.withAllInstanceIDs(c.requireCtx(ctx))
	ctxA, err := c.withAllLocalDevices(ctx)
	if err != nil {
		return nil, err
	}
	ctx = ctxA

	return c.APIClient.VolumeDetachAll(ctx, request)
}

func (c *client) VolumeDetachAllForService(
	ctx types.Context,
	service string,
	request *types.VolumeDetachRequest) (types.VolumeMap, error) {

	if c.isController() {
		return nil, utils.NewUnsupportedForClientTypeError(
			c.clientType, "VolumeDetachAllForService")
	}

	ctx = c.withInstanceID(c.requireCtx(ctx), service)
	ctxA, err := c.withAllLocalDevices(ctx)
	if err != nil {
		return nil, err
	}
	ctx = ctxA

	return c.APIClient.VolumeDetachAllForService(ctx, service, request)
}

func (c *client) VolumeSnapshot(
	ctx types.Context,
	service string,
	volumeID string,
	request *types.VolumeSnapshotRequest) (*types.Snapshot, error) {

	ctx = c.withInstanceID(c.requireCtx(ctx), service)
	return c.APIClient.VolumeSnapshot(ctx, service, volumeID, request)
}

func (c *client) Snapshots(
	ctx types.Context) (types.ServiceSnapshotMap, error) {

	ctx = c.withAllInstanceIDs(c.requireCtx(ctx))
	return c.APIClient.Snapshots(ctx)
}

func (c *client) SnapshotsByService(
	ctx types.Context, service string) (types.SnapshotMap, error) {

	ctx = c.withInstanceID(c.requireCtx(ctx), service)
	return c.APIClient.SnapshotsByService(ctx, service)
}

func (c *client) SnapshotInspect(
	ctx types.Context,
	service, snapshotID string) (*types.Snapshot, error) {

	ctx = c.withInstanceID(c.requireCtx(ctx), service)
	return c.APIClient.SnapshotInspect(ctx, service, snapshotID)
}

func (c *client) SnapshotRemove(
	ctx types.Context,
	service, snapshotID string) error {

	ctx = c.withInstanceID(c.requireCtx(ctx), service)
	return c.APIClient.SnapshotRemove(ctx, service, snapshotID)
}

func (c *client) SnapshotCopy(
	ctx types.Context,
	service, snapshotID string,
	request *types.SnapshotCopyRequest) (*types.Snapshot, error) {

	ctx = c.withInstanceID(c.requireCtx(ctx), service)
	return c.APIClient.SnapshotCopy(ctx, service, snapshotID, request)
}

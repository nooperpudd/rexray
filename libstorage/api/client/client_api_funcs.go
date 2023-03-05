package client

import (
	"bytes"
	"fmt"

	"github.com/nooperpudd/rexray/libstorage/api/types"
)

func (c *client) Root(ctx types.Context) ([]string, error) {

	reply := []string{}
	if _, err := c.httpGet(ctx, "/", &reply); err != nil {
		return nil, err
	}
	return reply, nil
}

const (
	ctxInstanceForSvc = 1000 + iota
)

type ctxInstanceForSvcT struct{}

func (c *client) Instances(
	ctx types.Context) (map[string]*types.Instance, error) {

	sis, err := c.Services(
		ctx.WithValue(ctxInstanceForSvc, &ctxInstanceForSvcT{}))
	if err != nil {
		return nil, err
	}
	instances := map[string]*types.Instance{}
	for service, si := range sis {
		instances[service] = si.Instance
	}
	return instances, nil
}

func (c *client) InstanceInspect(
	ctx types.Context, service string) (*types.Instance, error) {

	si, err := c.ServiceInspect(
		ctx.WithValue(ctxInstanceForSvc, &ctxInstanceForSvcT{}), service)
	if err != nil {
		return nil, err
	}
	return si.Instance, nil
}

func (c *client) Services(
	ctx types.Context) (map[string]*types.ServiceInfo, error) {

	reply := map[string]*types.ServiceInfo{}

	url := "/services"
	if ctx.Value(ctxInstanceForSvc) != nil {
		url = "/services?instance"
	}

	if _, err := c.httpGet(ctx, url, &reply); err != nil {
		return nil, err
	}
	return reply, nil
}

func (c *client) ServiceInspect(
	ctx types.Context, name string) (*types.ServiceInfo, error) {

	reply := &types.ServiceInfo{}

	url := fmt.Sprintf("/services/%s", name)
	if ctx.Value(ctxInstanceForSvc) != nil {
		url = fmt.Sprintf("/services/%s?instance", name)
	}

	if _, err := c.httpGet(ctx, url, &reply); err != nil {
		return nil, err
	}
	return reply, nil
}

func (c *client) Volumes(
	ctx types.Context,
	attachments types.VolumeAttachmentsTypes) (types.ServiceVolumeMap, error) {

	reply := types.ServiceVolumeMap{}
	url := fmt.Sprintf("/volumes?attachments=%v", attachments)
	if _, err := c.httpGet(ctx, url, &reply); err != nil {
		return nil, err
	}
	return reply, nil
}

func (c *client) VolumesByService(
	ctx types.Context,
	service string,
	attachments types.VolumeAttachmentsTypes) (types.VolumeMap, error) {

	reply := types.VolumeMap{}
	url := fmt.Sprintf("/volumes/%s?attachments=%v", service, attachments)
	if _, err := c.httpGet(ctx, url, &reply); err != nil {
		return nil, err
	}
	return reply, nil
}

func (c *client) VolumeInspect(
	ctx types.Context,
	service, volumeID string,
	attachments types.VolumeAttachmentsTypes) (*types.Volume, error) {

	reply := types.Volume{}
	url := fmt.Sprintf(
		"/volumes/%s/%s?attachments=%v", service, volumeID, attachments)
	if _, err := c.httpGet(ctx, url, &reply); err != nil {
		return nil, err
	}
	return &reply, nil
}

func (c *client) VolumeInspectByName(
	ctx types.Context,
	service, volumeName string,
	attachments types.VolumeAttachmentsTypes) (*types.Volume, error) {

	reply := types.Volume{}
	url := fmt.Sprintf(
		"/volumes/%s/%s?attachments=%v&byName", service, volumeName,
		attachments)
	if _, err := c.httpGet(ctx, url, &reply); err != nil {
		return nil, err
	}
	return &reply, nil
}

func (c *client) VolumeCreate(
	ctx types.Context,
	service string,
	request *types.VolumeCreateRequest) (*types.Volume, error) {

	reply := types.Volume{}
	if _, err := c.httpPost(ctx,
		fmt.Sprintf("/volumes/%s", service), request, &reply); err != nil {
		return nil, err
	}
	return &reply, nil
}

func (c *client) VolumeCreateFromSnapshot(
	ctx types.Context,
	service, snapshotID string,
	request *types.VolumeCreateRequest) (*types.Volume, error) {

	reply := types.Volume{}
	if _, err := c.httpPost(ctx,
		fmt.Sprintf("/snapshots/%s/%s?create",
			service, snapshotID), request, &reply); err != nil {
		return nil, err
	}
	return &reply, nil
}

func (c *client) VolumeCopy(
	ctx types.Context,
	service, volumeID string,
	request *types.VolumeCopyRequest) (*types.Volume, error) {

	reply := types.Volume{}
	if _, err := c.httpPost(ctx,
		fmt.Sprintf("/volumes/%s/%s?copy", service, volumeID),
		request, &reply); err != nil {
		return nil, err
	}
	return &reply, nil
}

func (c *client) VolumeRemove(
	ctx types.Context,
	service, volumeID string,
	force bool) error {

	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "/volumes/%s/%s", service, volumeID)
	if force {
		fmt.Fprintf(buf, "?force")
	}
	if _, err := c.httpDelete(ctx, buf.String(), nil); err != nil {
		return err
	}
	return nil
}

func (c *client) VolumeAttach(
	ctx types.Context,
	service string,
	volumeID string,
	request *types.VolumeAttachRequest) (*types.Volume, string, error) {

	reply := types.VolumeAttachResponse{}
	if _, err := c.httpPost(ctx,
		fmt.Sprintf("/volumes/%s/%s?attach",
			service, volumeID), request, &reply); err != nil {
		return nil, "", err
	}
	return reply.Volume, reply.AttachToken, nil
}

func (c *client) VolumeDetach(
	ctx types.Context,
	service string,
	volumeID string,
	request *types.VolumeDetachRequest) (*types.Volume, error) {

	reply := types.Volume{}
	if _, err := c.httpPost(ctx,
		fmt.Sprintf("/volumes/%s/%s?detach",
			service, volumeID), request, &reply); err != nil {
		return nil, err
	}
	return &reply, nil
}

func (c *client) VolumeDetachAll(
	ctx types.Context,
	request *types.VolumeDetachRequest) (types.ServiceVolumeMap, error) {

	reply := types.ServiceVolumeMap{}
	if _, err := c.httpPost(ctx,
		fmt.Sprintf("/volumes?detach"), request, &reply); err != nil {
		return nil, err
	}
	return reply, nil
}

func (c *client) VolumeDetachAllForService(
	ctx types.Context,
	service string,
	request *types.VolumeDetachRequest) (types.VolumeMap, error) {

	reply := types.VolumeMap{}
	if _, err := c.httpPost(ctx,
		fmt.Sprintf(
			"/volumes/%s?detach", service), request, &reply); err != nil {
		return nil, err
	}
	return reply, nil
}

func (c *client) VolumeSnapshot(
	ctx types.Context,
	service string,
	volumeID string,
	request *types.VolumeSnapshotRequest) (*types.Snapshot, error) {

	reply := types.Snapshot{}
	if _, err := c.httpPost(ctx,
		fmt.Sprintf("/volumes/%s/%s?snapshot",
			service, volumeID), request, &reply); err != nil {
		return nil, err
	}
	return &reply, nil
}

func (c *client) Snapshots(
	ctx types.Context) (types.ServiceSnapshotMap, error) {

	reply := types.ServiceSnapshotMap{}
	if _, err := c.httpGet(ctx, "/snapshots", &reply); err != nil {
		return nil, err
	}
	return reply, nil
}

func (c *client) SnapshotsByService(
	ctx types.Context, service string) (types.SnapshotMap, error) {

	reply := types.SnapshotMap{}
	if _, err := c.httpGet(ctx,
		fmt.Sprintf("/snapshots/%s", service), &reply); err != nil {
		return nil, err
	}
	return reply, nil
}

func (c *client) SnapshotInspect(
	ctx types.Context,
	service, snapshotID string) (*types.Snapshot, error) {

	reply := types.Snapshot{}
	if _, err := c.httpGet(ctx,
		fmt.Sprintf(
			"/snapshots/%s/%s", service, snapshotID), &reply); err != nil {
		return nil, err
	}
	return &reply, nil
}

func (c *client) SnapshotRemove(
	ctx types.Context,
	service, snapshotID string) error {

	if _, err := c.httpDelete(ctx,
		fmt.Sprintf("/snapshots/%s/%s", service, snapshotID), nil); err != nil {
		return err
	}
	return nil
}

func (c *client) SnapshotCopy(
	ctx types.Context,
	service, snapshotID string,
	request *types.SnapshotCopyRequest) (*types.Snapshot, error) {

	reply := types.Snapshot{}
	if _, err := c.httpPost(ctx,
		fmt.Sprintf("/snapshots/%s/%s?copy",
			service, snapshotID), request, &reply); err != nil {
		return nil, err
	}
	return &reply, nil
}

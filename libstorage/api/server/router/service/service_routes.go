package service

import (
	"net/http"

	"github.com/nooperpudd/rexray/libstorage/api/context"
	"github.com/nooperpudd/rexray/libstorage/api/server/httputils"
	"github.com/nooperpudd/rexray/libstorage/api/server/services"
	"github.com/nooperpudd/rexray/libstorage/api/types"
	"github.com/nooperpudd/rexray/libstorage/api/utils"
)

func (r *router) servicesList(
	ctx types.Context,
	w http.ResponseWriter,
	req *http.Request,
	store types.Store) error {

	reply := map[string]*types.ServiceInfo{}
	for service := range services.StorageServices(ctx) {
		ctx := context.WithStorageService(ctx, service)
		si, err := toServiceInfo(ctx, service, store)
		if err != nil {
			return err
		}
		reply[si.Name] = si
	}

	httputils.WriteJSON(w, http.StatusOK, reply)
	return nil
}

func (r *router) serviceInspect(
	ctx types.Context,
	w http.ResponseWriter,
	req *http.Request,
	store types.Store) error {

	service := context.MustService(ctx)
	si, err := toServiceInfo(ctx, service, store)
	if err != nil {
		return err
	}
	httputils.WriteJSON(w, http.StatusOK, si)
	return nil
}

func toServiceInfo(
	ctx types.Context,
	service types.StorageService,
	store types.Store) (*types.ServiceInfo, error) {

	d := service.Driver()

	var instance *types.Instance
	if store.GetBool("instance") {

		if _, ok := context.InstanceID(ctx); !ok {
			return nil, utils.NewMissingInstanceIDError(service.Name())
		}

		var err error
		instance, err = d.InstanceInspect(ctx, store)
		if err != nil {
			return nil, err
		}
	}

	st, err := d.Type(ctx)
	if err != nil {
		return nil, err
	}
	nd, err := d.NextDeviceInfo(ctx)
	if err != nil {
		return nil, err
	}

	return &types.ServiceInfo{
		Name:     service.Name(),
		Instance: instance,
		Driver: &types.DriverInfo{
			Name:       d.Name(),
			Type:       st,
			NextDevice: nd,
		},
	}, nil
}

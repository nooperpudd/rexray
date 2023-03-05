package libstorage

import (
	"strings"

	"github.com/nooperpudd/rexray/libstorage/api/context"
	"github.com/nooperpudd/rexray/libstorage/api/types"
)

func (c *client) requireCtx(ctx types.Context) types.Context {
	if ctx == nil {
		ctx = c.ctx
	} else {
		ctx = ctx.Join(c.ctx)
	}
	return context.RequireTX(ctx)
}

func (c *driver) requireCtx(ctx types.Context) types.Context {
	if ctx == nil {
		ctx = c.ctx
	} else {
		ctx = ctx.Join(c.ctx)
	}
	return context.RequireTX(ctx)
}

func (c *client) withAllInstanceIDs(ctx types.Context) types.Context {

	if c.isController() {
		return ctx
	}

	iidm := types.InstanceIDMap{}
	for _, k := range c.instanceIDCache.Keys() {
		iidm[k] = c.instanceIDCache.GetInstanceID(k)
	}

	if len(iidm) == 0 {
		return ctx
	}

	return ctx.WithValue(context.AllInstanceIDsKey, iidm)
}

func (c *client) withInstanceID(
	ctx types.Context, service string) types.Context {

	ctx = ctx.WithValue(context.ServiceKey, service)

	if c.isController() {
		return ctx
	}

	si := c.serviceCache.GetServiceInfo(service)
	if si == nil {
		return ctx
	}

	if !c.instanceIDCache.IsSet(si.Name) {
		return ctx
	}

	iid := c.instanceIDCache.GetInstanceID(si.Name)
	return ctx.WithValue(context.InstanceIDKey, iid)
}

func (c *client) withAllLocalDevices(ctx types.Context) (types.Context, error) {

	if c.isController() {
		return ctx, nil
	}

	ldm := types.LocalDevicesMap{}
	hit := map[string]bool{}

	for _, service := range c.serviceCache.Keys() {
		si := c.serviceCache.GetServiceInfo(service)
		dn := strings.ToLower(si.Driver.Name)
		if _, ok := hit[dn]; ok {
			continue
		}
		ctx := ctx.WithValue(context.ServiceKey, service)
		ld, err := c.LocalDevices(ctx, &types.LocalDevicesOpts{})
		if err != nil {
			if err == errExecutorNotSupported ||
				err == types.ErrNotImplemented {
				ctx.WithError(err).Warn("cannot get local deviecs")
				continue
			}
			return nil, err
		}
		hit[dn] = true
		ldm[dn] = ld
	}

	return ctx.WithValue(context.AllLocalDevicesKey, ldm), nil
}

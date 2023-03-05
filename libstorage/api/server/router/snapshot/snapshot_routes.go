package snapshot

import (
	"net/http"

	"github.com/akutz/goof"
	"github.com/nooperpudd/rexray/libstorage/api/context"
	"github.com/nooperpudd/rexray/libstorage/api/server/httputils"
	"github.com/nooperpudd/rexray/libstorage/api/server/router/volume"
	"github.com/nooperpudd/rexray/libstorage/api/server/services"
	"github.com/nooperpudd/rexray/libstorage/api/types"
	"github.com/nooperpudd/rexray/libstorage/api/utils"
	"github.com/nooperpudd/rexray/libstorage/api/utils/schema"
)

func (r *router) snapshots(
	ctx types.Context,
	w http.ResponseWriter,
	req *http.Request,
	store types.Store) error {

	var (
		tasks   = map[string]*types.Task{}
		taskIDs []int
		reply   = types.ServiceSnapshotMap{}
	)

	for service := range services.StorageServices(ctx) {

		run := func(
			ctx types.Context,
			svc types.StorageService) (interface{}, error) {

			ctx = context.WithStorageService(ctx, svc)
			var err error
			if ctx, err = context.WithStorageSession(ctx); err != nil {
				return nil, err
			}

			objs, err := svc.Driver().Snapshots(ctx, store)
			if err != nil {
				return nil, err
			}

			objMap := map[string]*types.Snapshot{}
			for _, obj := range objs {
				objMap[obj.ID] = obj
			}
			return objMap, nil
		}

		task := service.TaskEnqueue(ctx, run, schema.SnapshotMapSchema)
		taskIDs = append(taskIDs, task.ID)
		tasks[service.Name()] = task
	}

	run := func(ctx types.Context) (interface{}, error) {

		services.TaskWaitAll(ctx, taskIDs...)

		for k, v := range tasks {
			if v.Error != nil {
				return nil, utils.NewBatchProcessErr(reply, v.Error)
			}

			objMap, ok := v.Result.(map[string]*types.Snapshot)
			if !ok {
				return nil, utils.NewBatchProcessErr(
					reply, goof.New("error casting to []*types.Snapshot"))
			}
			reply[k] = objMap
		}

		return reply, nil
	}

	return httputils.WriteTask(
		ctx,
		r.config,
		w,
		store,
		services.TaskEnqueue(ctx, run, schema.ServiceSnapshotMapSchema),
		http.StatusOK)
}

func (r *router) snapshotsForService(
	ctx types.Context,
	w http.ResponseWriter,
	req *http.Request,
	store types.Store) error {

	service := context.MustService(ctx)

	run := func(
		ctx types.Context,
		svc types.StorageService) (interface{}, error) {

		var reply types.SnapshotMap = map[string]*types.Snapshot{}

		objs, err := svc.Driver().Snapshots(ctx, store)
		if err != nil {
			return nil, err
		}

		for _, obj := range objs {
			reply[obj.ID] = obj
		}
		return reply, nil
	}

	return httputils.WriteTask(
		ctx,
		r.config,
		w,
		store,
		service.TaskEnqueue(ctx, run, schema.SnapshotMapSchema),
		http.StatusOK)
}

func (r *router) snapshotInspect(
	ctx types.Context,
	w http.ResponseWriter,
	req *http.Request,
	store types.Store) error {

	service := context.MustService(ctx)

	run := func(
		ctx types.Context,
		svc types.StorageService) (interface{}, error) {

		return svc.Driver().SnapshotInspect(
			ctx,
			store.GetString("snapshotID"),
			store)
	}

	return httputils.WriteTask(
		ctx,
		r.config,
		w,
		store,
		service.TaskEnqueue(ctx, run, schema.SnapshotSchema),
		http.StatusOK)
}

func (r *router) snapshotRemove(
	ctx types.Context,
	w http.ResponseWriter,
	req *http.Request,
	store types.Store) error {

	service := context.MustService(ctx)

	run := func(
		ctx types.Context,
		svc types.StorageService) (interface{}, error) {

		return nil, svc.Driver().SnapshotRemove(
			ctx,
			store.GetString("snapshotID"),
			store)
	}

	return httputils.WriteTask(
		ctx,
		r.config,
		w,
		store,
		service.TaskEnqueue(ctx, run, nil),
		http.StatusResetContent)
}

func (r *router) volumeCreate(
	ctx types.Context,
	w http.ResponseWriter,
	req *http.Request,
	store types.Store) error {

	service := context.MustService(ctx)

	run := func(
		ctx types.Context,
		svc types.StorageService) (interface{}, error) {

		v, err := svc.Driver().VolumeCreateFromSnapshot(
			ctx,
			store.GetString("snapshotID"),
			store.GetString("name"),
			&types.VolumeCreateOpts{
				AvailabilityZone: store.GetStringPtr("availabilityZone"),
				IOPS:             store.GetInt64Ptr("iops"),
				Size:             store.GetInt64Ptr("size"),
				Type:             store.GetStringPtr("type"),
				Opts:             store,
			})

		if err != nil {
			return nil, err
		}

		if volume.OnVolume != nil {
			ok, err := volume.OnVolume(ctx, req, store, v)
			if err != nil {
				return nil, err
			}
			if !ok {
				return nil, utils.NewNotFoundError(v.ID)
			}
		}

		return v, nil
	}

	return httputils.WriteTask(
		ctx,
		r.config,
		w,
		store,
		service.TaskEnqueue(ctx, run, schema.VolumeSchema),
		http.StatusCreated)
}

func (r *router) snapshotCopy(
	ctx types.Context,
	w http.ResponseWriter,
	req *http.Request,
	store types.Store) error {

	service := context.MustService(ctx)

	run := func(
		ctx types.Context,
		svc types.StorageService) (interface{}, error) {

		return svc.Driver().SnapshotCopy(
			ctx,
			store.GetString("snapshotID"),
			store.GetString("snapshotName"),
			store.GetString("destinationID"),
			store)
	}

	return httputils.WriteTask(
		ctx,
		r.config,
		w,
		store,
		service.TaskEnqueue(ctx, run, schema.SnapshotSchema),
		http.StatusCreated)
}

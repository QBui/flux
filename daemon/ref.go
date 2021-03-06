package daemon

import (
	"context"
	"sync"

	"github.com/weaveworks/flux/api"
	"github.com/weaveworks/flux/api/v6"
	"github.com/weaveworks/flux/api/v9"
	"github.com/weaveworks/flux/job"
	"github.com/weaveworks/flux/update"
)

// Ref is a cell containing a server implementation, that we can
// update atomically. The point of this is to be able to have a
// server in use (e.g., answering RPCs), and swap it later when the
// state changes.
type Ref struct {
	sync.RWMutex
	server api.UpstreamServer
}

func NewRef(server api.UpstreamServer) *Ref {
	return &Ref{server: server}
}

func (r *Ref) Server() api.UpstreamServer {
	r.RLock()
	defer r.RUnlock()
	return r.server
}

func (r *Ref) UpdateServer(server api.UpstreamServer) {
	r.Lock()
	r.server = server
	r.Unlock()
}

// api.Server implementation so clients don't need to be refactored around
// Server() API

func (r *Ref) Ping(ctx context.Context) error {
	return r.Server().Ping(ctx)
}

func (r *Ref) Version(ctx context.Context) (string, error) {
	return r.Server().Version(ctx)
}

func (r *Ref) Export(ctx context.Context) ([]byte, error) {
	return r.Server().Export(ctx)
}

func (r *Ref) ListServices(ctx context.Context, namespace string) ([]v6.ControllerStatus, error) {
	return r.Server().ListServices(ctx, namespace)
}

func (r *Ref) ListImages(ctx context.Context, spec update.ResourceSpec) ([]v6.ImageStatus, error) {
	return r.Server().ListImages(ctx, spec)
}

func (r *Ref) UpdateManifests(ctx context.Context, spec update.Spec) (job.ID, error) {
	return r.Server().UpdateManifests(ctx, spec)
}

func (r *Ref) NotifyChange(ctx context.Context, change v9.Change) error {
	return r.Server().NotifyChange(ctx, change)
}

func (r *Ref) JobStatus(ctx context.Context, id job.ID) (job.Status, error) {
	return r.Server().JobStatus(ctx, id)
}

func (r *Ref) SyncStatus(ctx context.Context, ref string) ([]string, error) {
	return r.Server().SyncStatus(ctx, ref)
}

func (r *Ref) GitRepoConfig(ctx context.Context, regenerate bool) (v6.GitConfig, error) {
	return r.Server().GitRepoConfig(ctx, regenerate)
}

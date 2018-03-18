package util

import (
	"context"
	"sync"

	"github.com/thetechnick/hcloud-ansible/pkg/hcloud"
)

// WaitFn waits for the competion of an action
type WaitFn func(ctx context.Context, client *hcloud.Client, action *hcloud.Action) error

var actionWaitLock sync.Mutex

// WaitForAction makes sure we only watch one action at a time
func WaitForAction(ctx context.Context, client *hcloud.Client, action *hcloud.Action) error {
	actionWaitLock.Lock()
	defer actionWaitLock.Unlock()
	_, errCh := client.Action.WatchProgress(ctx, action)
	return <-errCh
}

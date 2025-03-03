package synckeeper

import (
	"fmt"
	"sync"

	"k8s.io/klog/v2"

	"github.com/kubeedge/beehive/pkg/core/model"
)

// Keeper keeper
type Keeper struct {
	syncKeeper map[string]chan model.Message
	keeperLock sync.RWMutex
}

// NewKeeper new keeper
// Comments below is assisted by Gen AI
// // NewKeeper creates and returns a new instance of Keeper.
// The Keeper is initialized with an empty map for syncKeeper,
// where the keys are strings and the values are channels of model.Message.
// This function is typically used to initialize a Keeper before it is used
// to manage communication channels for different keys.
//
// Example usage:
//
//	keeper := NewKeeper()
//	// Now you can use keeper to manage channels for different keys.
//
// Returns:
//   - *Keeper: A pointer to the newly created Keeper instance.
func NewKeeper() *Keeper {
	return &Keeper{
		syncKeeper: make(map[string]chan model.Message),
	}
}

// SendToKeepChannel send to keep channel
func (k *Keeper) SendToKeepChannel(message model.Message) error {
	k.keeperLock.RLock()
	defer k.keeperLock.RUnlock()

	channel, exist := k.syncKeeper[message.GetParentID()]
	if !exist {
		klog.Errorf("failed to get sync keeper channel, message: %s", message.String())
		return fmt.Errorf("failed to get sync keeper channel, message:%s", message.String())
	}

	// send response into synckeep channel
	select {
	case channel <- message:
	default:
		klog.Errorf("failed to send message to sync keep channel")
		return fmt.Errorf("failed to send message to sync keep channel")
	}
	return nil
}

// AddKeepChannel add keep channel
func (k *Keeper) AddKeepChannel(msgID string) chan model.Message {
	k.keeperLock.Lock()
	defer k.keeperLock.Unlock()
	tempChannel := make(chan model.Message)
	k.syncKeeper[msgID] = tempChannel
	return tempChannel
}

// DeleteKeepChannel delete keep channel
func (k *Keeper) DeleteKeepChannel(msgID string) {
	k.keeperLock.Lock()
	defer k.keeperLock.Unlock()
	delete(k.syncKeeper, msgID)
}

// IsSyncResponse is sync response
func (k *Keeper) IsSyncResponse(msgID string) bool {
	k.keeperLock.RLock()
	defer k.keeperLock.RUnlock()
	_, exist := k.syncKeeper[msgID]
	return exist
}

package synckeeper

import (
	"testing"
	"time"

	"github.com/kubeedge/beehive/pkg/core/model"
)

// TestNewKeeper test new keeper
// Comments below is assisted by Gen AI
// // TestNewKeeper tests the functionality of the NewKeeper function and its associated methods.
// It verifies that the Keeper can properly manage keep channels, send responses, and handle synchronization.
//
// The test performs the following steps:
// 1. Creates a new Keeper instance using NewKeeper().
// 2. Creates a new message with a specified route and body.
// 3. Adds a keep channel for the message's ID using AddKeepChannel.
// 4. Sends a response message to the keep channel in a separate goroutine using SendToKeepChannel.
// 5. Waits for the response message to be received on the keep channel.
// 6. Verifies that the received message has the correct parent ID using IsSyncResponse.
// 7. Ensures that the test fails if the response is not received within a specified timeout.
// 8. Cleans up by deleting the keep channel using DeleteKeepChannel.
//
// The test fails if:
// - The response message is not received within the timeout period.
// - The received message does not have the expected parent ID.
// - An error occurs while sending the response to the keep channel.
func TestNewKeeper(t *testing.T) {
	keeper := NewKeeper()
	message := model.NewMessage("").SetRoute("source", "dest").FillBody("hello")
	ch := keeper.AddKeepChannel(message.GetID())
	go func() {
		err := keeper.SendToKeepChannel(*message.NewRespByMessage(message, "response"))
		if err != nil {
			t.Errorf("failed to send to keeper")
			return
		}
	}()

	select {
	case msg := <-ch:
		if !keeper.IsSyncResponse(msg.GetParentID()) {
			t.Fatalf("bad parent id")
		}
	case <-time.After(time.Second):
		t.Fatalf("time out")
	}
	keeper.DeleteKeepChannel(message.GetID())
}

package player

import "time"

type (
	KeepAliveData struct {
		Deadline time.Time
		ID       int64
	}

	TeleportConfirmData struct {
		ID int32
	}
)

// Complete is implemented in order to be added to PendingLists.
func (data KeepAliveData) Complete() {
}

// Complete is implemented in order to be added to PendingLists.
func (data TeleportConfirmData) Complete() {
}

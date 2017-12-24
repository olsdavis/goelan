package util

// PendingList is a struct which represents actions that have not
// been completed yet. For instance, until the player sends the
// keep alive back, the pending keep alive is added to this list
// with Append(Completable), and removed once the response to the
// keep alive is received with QueryAndComplete(Matcher).
type PendingList struct {
	pending []Completable
}

// NewPendingList creates and returns a new PendingList.
func NewPendingList() *PendingList {
	return &PendingList{make([]Completable, 0)}
}

// Elements returns the elements that PendingList currently contains.
func (list *PendingList) Elements() []Completable {
	return list.pending
}

// Append adds to the PendingList the given completable.
// (This function ignores the eventual duplicates.)
func (list *PendingList) Append(completable Completable) {
	list.pending = append(list.pending, completable)
}

// QueryAndComplete looks up for the first element which matches
// the query's criteria, and, if found, completes and removes
// the pending action from the PendingList. Returns true if
// an element has been removed; otherwise, returns false.
// First in, first out.
func (list *PendingList) QueryAndComplete(query Matcher) bool {
	newList := make([]Completable, 0, 1)
	deleted := false
	for _, completable := range list.pending {
		// check if already deleted, in order to remove only the first match
		if !deleted && query(completable) {
			completable.Complete()
			deleted = true
		} else {
			newList = append(newList, completable)
		}
	}
	list.pending = newList
	return deleted
}

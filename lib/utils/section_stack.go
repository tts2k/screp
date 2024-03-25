package container

import (
	"container/list"
	"sync"

	"screp/lib/model"
)

type SectionStack struct {
	data  *list.List
	mutex sync.Mutex
}

// IsEmpty: check if stack is empty
func (st *SectionStack) IsEmpty() bool {
	return st.data == nil || st.data.Len() == 0
}

// Push: push a new element to stack
func (st *SectionStack) Push(section *model.Section) {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	if st.data == nil {
		st.data = list.New()
	}

	st.data.PushFront(section)
}

// Peek: Get the top item on stack without removing it
func (st *SectionStack) Peek() *model.Section {
	return st.data.Front().Value.(*model.Section)
}

// Pop: Remove the top item on the stack and return it
func (st *SectionStack) Pop() *model.Section {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	if st.IsEmpty() {
		return nil
	}

	return st.data.Remove(st.data.Front()).(*model.Section)
}

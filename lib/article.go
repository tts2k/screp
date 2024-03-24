package lib

import (
	"container/list"
	"sync"
)

type Section struct {
	Title       string
	Content     []string
	SubSections []Section
}

type Article struct {
	Title    string
	Preamble string
	Sections []Section
}

type SectionStack struct {
	data  *list.List
	mutex sync.Mutex
}

// IsEmpty: check if stack is empty
func (st *SectionStack) IsEmpty() bool {
	return st.data == nil || st.data.Len() == 0
}

// Push: push a new element to stack
func (st *SectionStack) Push(section *Section) {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	if st.data == nil {
		st.data = list.New()
	}

	st.data.PushFront(section)
}

// Peek: Get the top item on stack without removing it
func (st *SectionStack) Peek() *Section {
	return st.data.Front().Value.(*Section)
}

// Pop: Remove the top item on the stack and return it
func (st *SectionStack) Pop() *Section {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	if st.IsEmpty() {
		return nil
	}

	return st.data.Remove(st.data.Front()).(*Section)
}

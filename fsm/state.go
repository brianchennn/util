package fsm

import "sync"

type StateType string

// State is a thread-safe structure that represents the state of a object
// and it can be used for FSM
type State struct {
	// Current state of the State object
	Current StateType `json:"Current"`
	// StateMutex ensures that all operations to Current is thread-safe
	StateMutex sync.RWMutex `json:"-"`
}

// NewState create a State object with Current state set to initState
func NewState(initState StateType) *State {
	s := &State{Current: initState}
	return s
}

// Current get the Current state
func (state *State) GetCurrent() StateType {
	state.StateMutex.RLock()
	defer state.StateMutex.RUnlock()
	return state.Current
}

// Is return true if the Current state is equal to target
func (state *State) Is(target StateType) bool {
	state.StateMutex.RLock()
	defer state.StateMutex.RUnlock()
	return state.Current == target
}

// Set Current state to next
func (state *State) Set(next StateType) {
	state.StateMutex.Lock()
	defer state.StateMutex.Unlock()
	state.Current = next
}

package states

// import "github.com/Jasrags/NewMUD/users"

// type Machine struct {
// 	User   *users.User
// 	States map[string]func(input string)
// }

// func NewMachine(u users.User) *Machine {
// 	return &Machine{
// 		User:   &u,
// 		States: make(map[string]func(input string)),
// 	}
// }

// func (sm *Machine) RegisterState(state string, handler func(input string)) {
// 	sm.States[state] = handler
// }

// func (sm *Machine) HandleInput(input string) {
// 	if handler, exists := sm.States[sm.User.State]; exists {
// 		handler(input)
// 	} else {
// 		// sm.user.Conn.Write([]byte("Invalid state.\n"))
// 	}
// }

// func (sm *Machine) TransitionTo(state string) {
// 	sm.User.State = state
// 	if handler, exists := sm.States[state]; exists {
// 		handler("")
// 	}
// }

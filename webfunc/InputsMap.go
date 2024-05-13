package webfunc

import (
	"GT/games"
	"sync"
)

type UsersInputs struct {
	Inputs map[string]games.Input
	mutex  sync.Mutex
}

func NewUsersInputs() UsersInputs {
	return UsersInputs{
		Inputs: make(map[string]games.Input),
		mutex:  sync.Mutex{},
	}
}

func (ui *UsersInputs) Add(user string, inputs games.Input) {
	ui.mutex.Lock()
	defer ui.mutex.Unlock()
	ui.Inputs[user] = inputs
}

func (ui *UsersInputs) Delete(user string) {
	ui.mutex.Lock()
	defer ui.mutex.Unlock()
	delete(ui.Inputs, user)
}

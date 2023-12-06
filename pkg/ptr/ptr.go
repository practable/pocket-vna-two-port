// package ptr is a test package to check out how to handle rfusb.Switch interface
// when using pointer receiver functions
package ptr

type Knife struct {
	//mu   *sync.Mutex
	Name  string
	State bool
}

type Button struct {
	//mu     *sync.Mutex
	Colour string
	State  bool
}

type Switch interface {
	Toggle() error
	GetState() bool
}

func (k *Knife) Toggle() error {
	k.State = !k.State
	return nil
}

func (b *Button) Toggle() error {
	b.State = !b.State
	return nil
}

func (k *Knife) GetState() bool {
	return k.State
}

func (b *Button) GetState() bool {
	return b.State
}

type Decider struct {
	Switch      Switch
	Description string
}

func NewDecider(s Switch, d string) *Decider {
	return &Decider{
		Switch:      s,
		Description: d,
	}
}

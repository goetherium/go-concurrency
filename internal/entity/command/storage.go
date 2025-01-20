package command

type Key string

type Value string

func (v Value) String() string {
	return string(v)
}

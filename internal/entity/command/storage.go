package command

type Key Arg

func NewKey(key Arg) Key {
	return Key(key)
}

type Value Arg

func NewValue(value Arg) Value {
	return Value(value)
}

func (v Value) String() string {
	return string(v)
}

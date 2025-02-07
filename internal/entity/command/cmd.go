package command

import (
	"strconv"
)

const (
	UnknownCommandID CmdID = iota
	SetCommandID
	GetCommandID
	DelCommandID
)

var cmdArgsCount = map[CmdID]int{
	SetCommandID: 2,
	GetCommandID: 1,
	DelCommandID: 1,
}

var cmdToRawCmd = map[CmdID]RawCmd{
	SetCommandID: setCommand,
	GetCommandID: getCommand,
	DelCommandID: delCommand,
}

// walCmd список команд, которые нужно записывать в wal.
var walCmd = map[CmdID]struct{}{
	SetCommandID: {},
	DelCommandID: {},
}

type CmdID int

func (id CmdID) ArgsCount() int {
	return cmdArgsCount[id]
}

func (id CmdID) Int() int {
	return int(id)
}

func (id CmdID) String() string {
	return strconv.Itoa(int(id))
}

func (id CmdID) IsWalCmd() bool {
	_, ok := walCmd[id]

	return ok
}

// Cmd команда после разбора RawCmd.
type Cmd struct {
	id   CmdID
	args Args
}

func NewCmd(id CmdID, args Args) Cmd {
	return Cmd{
		id:   id,
		args: args,
	}
}

func (c Cmd) ID() CmdID {
	return c.id
}

func (c Cmd) Args() Args {
	return c.args
}

func (c Cmd) Arg(idx int) Arg {
	return c.args[idx]
}

const marshalDelimiter = " "

func (c Cmd) Marshal() string {
	str := cmdToRawCmd[c.id].String()

	for _, arg := range c.args {
		str += marshalDelimiter + arg.String()
	}

	return str
}

type Args []Arg

func NewArgs(in []string) Args {
	args := make([]Arg, len(in))

	for i := range in {
		args[i] = Arg(in[i])
	}

	return args
}

type Arg string

func (a Arg) String() string {
	return string(a)
}

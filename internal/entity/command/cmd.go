package command

import (
	"strconv"
)

type RawCmd string

const (
	setCommand RawCmd = "SET"
	getCommand RawCmd = "GET"
	delCommand RawCmd = "DEL"
)

func NewRawCmd(rawCmd string) RawCmd {
	return RawCmd(rawCmd)
}

func (r RawCmd) ToCmdID() CmdID {
	cmdID, ok := rawCmdToCmdID[r]
	if !ok {
		return UnknownCommandID
	}

	return cmdID
}

func (r RawCmd) String() string {
	return string(r)
}

var rawCmdToCmdID = map[RawCmd]CmdID{
	setCommand: SetCommandID,
	getCommand: GetCommandID,
	delCommand: DelCommandID,
}

type CmdID int

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

func (c CmdID) ArgsCount() int {
	return cmdArgsCount[c]
}

func (c CmdID) Int() int {
	return int(c)
}

func (c CmdID) String() string {
	return strconv.Itoa(int(c))
}

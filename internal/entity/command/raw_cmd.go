package command

// RawCmd текст команды, полученный от пользователя или считанный из WAL.
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

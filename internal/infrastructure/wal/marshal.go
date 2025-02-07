package wal

import (
	"bytes"
	"strconv"

	"database/internal/entity/command"
)

const (
	cmdDelimiter = " "
	newLine      = "\n"
)

func marshalCmd(buf *bytes.Buffer, lsn Lsn, cmd command.Cmd) {
	lsnStr := lsnPfx + strconv.FormatUint(uint64(lsn), 10) + "]"
	buf.WriteString(lsnStr + cmdDelimiter + cmd.Marshal() + newLine)
}

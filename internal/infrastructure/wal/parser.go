package wal

import (
	"errors"
	"strconv"
	"strings"

	"database/internal/entity/command"
)

var (
	errLSNNotFound     = errors.New("LSN not found")
	errCommandNotFound = errors.New("command not found")
)

const lsnPfx = "[LSN="

func parseLine(line string) (Lsn, command.RawCmd, error) {
	if len(line) < len(lsnPfx) {
		return 0, "", errLSNNotFound
	}

	tail := line[len(lsnPfx):]

	if len(tail) < 2 {
		return 0, "", errLSNNotFound
	}

	idx := strings.Index(tail, "]")
	if idx < 0 {
		return 0, "", errLSNNotFound
	}

	lsnStr := tail[:idx]
	if lsnStr == "" {
		return 0, "", errLSNNotFound
	}

	lsn, err := strconv.ParseUint(lsnStr, 10, 64)
	if err != nil {
		return 0, "", errLSNNotFound
	}

	cmd := tail[idx+1:]
	cmd = strings.Trim(cmd, " ")
	if cmd == "" {
		return 0, "", errCommandNotFound
	}

	return newLsn(lsn), command.NewRawCmd(cmd), nil
}

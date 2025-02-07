package wal

import (
	"testing"

	"database/internal/entity/command"
	wal2 "database/internal/entity/wal"
	"github.com/stretchr/testify/assert"
)

func Test_parseLine(t *testing.T) {
	tests := []struct {
		name    string
		line    string // Строка WAL-файла
		wantLsn uint64 // Найденный lsn
		wantCmd string // Найденная команда
		err     error
	}{
		{
			name:    "empty line",
			line:    "",
			wantLsn: 0,
			err:     errLSNNotFound,
		},
		{
			name:    "no_lsn",
			line:    "[]",
			wantLsn: 0,
			err:     errLSNNotFound,
		},
		{
			name:    "invalid_lsn",
			line:    "[LSN=",
			wantLsn: 0,
			err:     errLSNNotFound,
		},
		{
			name:    "invalid_lsn_2",
			line:    "[LSN=1",
			wantLsn: 0,
			err:     errLSNNotFound,
		},
		{
			name:    "lsn_wo_command",
			line:    "[LSN=123]",
			wantLsn: 0,
			err:     errCommandNotFound,
		},
		{
			name:    "lsn_and_command",
			line:    "[LSN=123]set key value",
			wantLsn: 123,
			wantCmd: "set key value",
			err:     nil,
		},
		{
			name:    "lsn_and_command",
			line:    "[LSN=456] set key value",
			wantLsn: 456,
			wantCmd: "set key value",
			err:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLsn, gotCmd, err := parseLine(tt.line)
			assert.Equal(t, wal2.newLsn(tt.wantLsn), gotLsn)
			assert.Equal(t, command.NewRawCmd(tt.wantCmd), gotCmd)
			assert.Equal(t, tt.err, err)
		})
	}
}

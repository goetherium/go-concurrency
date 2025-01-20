package storage

import (
	"database/internal/entity/command"
)

// ExecResult результат выполнения команды слоем storage
type ExecResult struct {
	Result string
	Value  command.Value
}

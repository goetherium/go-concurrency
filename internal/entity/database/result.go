package database

// HandleCmdResult результат обработки команды слоем database
type HandleCmdResult struct {
	Result string
}

func (r HandleCmdResult) String() string {
	return r.Result
}

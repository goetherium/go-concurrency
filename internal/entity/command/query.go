package command

type Query struct {
	CmdID CmdID
	Args  Args
}

type Args []string

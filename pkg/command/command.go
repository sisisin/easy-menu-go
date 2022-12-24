package command

import (
	"os/exec"

	m "github.com/sisisin/easy-menu-go/pkg/menu"
)

type CommandProcessState uint32

const (
	NotExecuting CommandProcessState = 1 << iota
	Ready
	Executing
	Succeeded
	Failed
)

type CommandState struct {
	Command      string
	ProcessState CommandProcessState
	Err          error
	Cmd          *exec.Cmd
}

func ExecuteCommand(rootMenu m.MenuItem, cursor []int64) CommandState {
	current := getCurrent(rootMenu, cursor)
	command := current.Command.Command
	cmd := exec.Command("sh", "-c", command)

	println("current workdir: ", current.WorkDir)
	return CommandState{
		Command:      command,
		ProcessState: Ready,
		Err:          nil,
		Cmd:          cmd,
	}
}

func getCurrent(rootMenu m.MenuItem, cursor []int64) m.MenuItem {
	target := rootMenu
	for _, v := range cursor {
		target = target.SubMenu.Items[v]
	}
	return target
}

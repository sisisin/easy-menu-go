package command

import (
	"os/exec"
	m "sisisin/easy-menu-go/pkg/menu"
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

func ExecuteCommand(current m.MenuItem) CommandState {
	command := current.Command.Command
	cmd := exec.Command("sh", "-c", command)

	return CommandState{
		Command:      command,
		ProcessState: Ready,
		Err:          nil,
		Cmd:          cmd,
	}
}
func ExecuteCommand2(current m.MenuItem) CommandState {
	command := current.Command.Command
	cmd := exec.Command("sh", "-c", command)

	if err := cmd.Start(); err != nil {
		return CommandState{
			Command:      command,
			ProcessState: Failed,
			Err:          err,
		}
	}
	return CommandState{
		Command:      command,
		ProcessState: Executing,
		Err:          nil,
	}
}

package command

import (
	"os/exec"
	m "sisisin/easy-menu-go/pkg/menu"
)

type CommandProcessState uint32

const (
	NotExecuting CommandProcessState = 1 << iota
	Executing
	Succeeded
	Failed
)

type CommandState struct {
	Command      string
	ProcessState CommandProcessState
	Err          error
}

func ExecuteCommand(current m.MenuItem) CommandState {
	command := current.Command.Command
	cmd := exec.Command("sh", "-c", command)

	if err := cmd.Start(); err != nil {
		return CommandState{
			Command:      command,
			ProcessState: Failed,
			Err:          err,
		}
	}
	if err := cmd.Wait(); err != nil {
		return CommandState{
			Command:      command,
			ProcessState: Failed,
			Err:          err,
		}
	}
	return CommandState{
		Command:      command,
		ProcessState: Succeeded,
		Err:          nil,
	}
}

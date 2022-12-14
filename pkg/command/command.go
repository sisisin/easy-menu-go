package command

import (
	"fmt"
	"os"
	"os/exec"
	"path"

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

func GetSelectedCommandState(rootMenu m.MenuItem, cursor []int64, configFile string) CommandState {
	current := getCurrent(rootMenu, cursor)
	command := current.Command.Command
	cmd := exec.Command("sh", "-c", command)
	dir, err := getDir(rootMenu, cursor, configFile)

	if err != nil {
		return CommandState{
			Command:      command,
			ProcessState: Failed,
			Err:          err,
			Cmd:          nil,
		}
	}
	cmd.Dir = dir

	cmd.Env = getEnv(rootMenu, cursor)
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

func getEnv(rootMenu m.MenuItem, cursor []int64) []string {
	var envMap map[string]string = nil
	target := rootMenu
	if target.Env != nil {
		envMap = target.Env
	}
	for _, v := range cursor {
		target = target.SubMenu.Items[v]
		if target.Env != nil {
			envMap = target.Env
		}
	}

	var envArr []string
	for k, v := range envMap {
		envArr = append(envArr, fmt.Sprintf("%s=%s", k, v))
	}

	return envArr
}

func getDir(rootMenu m.MenuItem, cursor []int64, configFile string) (string, error) {
	base := path.Dir(configFile)

	workDir := ""
	target := rootMenu
	if target.WorkDir != "" {
		workDir = target.WorkDir
	}
	for _, v := range cursor {
		target = target.SubMenu.Items[v]
		if target.WorkDir != "" {
			workDir = target.WorkDir
		}
	}

	var dir string
	if path.IsAbs(workDir) {
		dir = workDir
	} else {

		dir = path.Join(base, workDir)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return "", fmt.Errorf("invalid config `work_dir`, no such directory.\nused work_dir: `%v`\nresolved path: `%v`", workDir, dir)
	}

	return dir, nil
}

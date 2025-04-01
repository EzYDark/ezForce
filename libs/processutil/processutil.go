package processutil

import (
	"fmt"
	"os/exec"

	"github.com/ezydark/warpenforcer/libs/logger"
	"github.com/shirou/gopsutil/process"
)

func IsProcessRunningByName(name string) (bool, error) {
	processes, err := process.Processes()
	if err != nil {
		return false, fmt.Errorf("error while getting list of processes:\n %w", err)
	}

	for _, p := range processes {
		processName, err := p.Name()
		if err != nil {
			return false, fmt.Errorf("Could not get process name:\n %w", err)
		}
		if processName == name {
			return true, nil
		}
	}
	return false, nil
}

func StartProcess(path string, args ...string) error {
	cmd := exec.Command(path, args...)

	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start process:\n %w", err)
	}

	processID := cmd.Process.Pid

	log, err := logger.Get()
	if err != nil {
		return fmt.Errorf("failed to get logger:\n %w", err)
	}

	log.Info().Msgf("Started process '%s' with PID '%d'", path, processID)

	if err := cmd.Process.Release(); err != nil {
		return fmt.Errorf("failed to release process:\n %w", err)
	}

	return nil
}

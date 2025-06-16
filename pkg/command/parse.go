package command

import (
	"bufio"
	"os/exec"
)

func RunAndParse(c *exec.Cmd) ([]string, error) {
	stdout, err := c.StdoutPipe()
	if err != nil {
		return nil, err
	}
	var result []string
	var done bool
	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			if done {
				break
			}
			result = append(result, scanner.Text())
		}
	}()
	err = c.Run()
	if err != nil {
		return nil, err
	}
	done = true
	return result, nil
}

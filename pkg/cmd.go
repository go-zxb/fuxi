package pkg

import (
	"errors"
	"os"
	"os/exec"
)

func RunCommand(name string, args ...string) error {
	command := exec.Command(name, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err := command.Start()
	if err != nil {
		return errors.New("⚠️❎ 启动命令时出错:" + err.Error())
	}

	err = command.Wait()
	if err != nil {
		return errors.New("⚠️❎ 命令已完成，但有错误:" + err.Error())
	}
	return err
}

func RunCommandNoOutput(name string, args ...string) error {
	command := exec.Command(name, args...)
	err := command.Start()
	if err != nil {
		return errors.New("⚠️❎ 启动命令时出错:" + err.Error())
	}

	err = command.Wait()
	if err != nil {
		return errors.New("⚠️❎ 命令已完成，但有错误:" + err.Error())
	}
	return err
}

package pkg

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
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

// CommandInfo 定义一个通道类型，用于发送命令执行信息
type CommandInfo struct {
	Message string
	Error   error
}

func RunCommandChannel(ctx context.Context, infoChan chan<- CommandInfo, name string, args ...string) {
	command := exec.CommandContext(ctx, name, args...)
	defer close(infoChan) // 确保在函数返回时关闭通道

	// 创建一个管道用于捕获标准输出和标准错误
	stdout, err := command.StdoutPipe()
	if err != nil {
		fmt.Println(err.Error())
		infoChan <- CommandInfo{Message: "⚠️❎ 运行命令出错了", Error: err}
		return
	}

	stderr, err := command.StderrPipe()
	if err != nil {
		fmt.Println(err.Error())
		infoChan <- CommandInfo{Message: "⚠️❎ 运行命令出错了", Error: err}
		return
	}

	// 合并标准输出和标准错误
	mergedOutput := io.MultiReader(stdout, stderr)

	// 启动命令
	err = command.Start()
	if err != nil {
		fmt.Println(err.Error())
		infoChan <- CommandInfo{Message: "⚠️❎ 运行命令出错了", Error: err}
		return
	}

	// 使用 WaitGroup 来等待所有发送操作完成
	var wg sync.WaitGroup
	wg.Add(1)

	// 创建一个 goroutine 来读取标准输出和标准错误并发送给通道
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(mergedOutput)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				fmt.Println("收到结束命令")
				return
			case infoChan <- CommandInfo{Message: scanner.Text(), Error: nil}:
			}
		}
	}()

	// 等待命令完成
	err = command.Wait()
	if err != nil {
		infoChan <- CommandInfo{Message: "⚠️❎ 命令已完成，但有错误", Error: err}
	} else {
		infoChan <- CommandInfo{Message: "命令执行成功", Error: nil}
	}

	// 等待所有发送操作完成后再关闭通道
	wg.Wait()
}

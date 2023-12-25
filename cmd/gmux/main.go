package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"
)

// Config 结构体用于匹配 YAML 配置文件的结构
type Config struct {
	Name    string   `yaml:"name"`
	Root    string   `yaml:"root"`
	Windows []Window `yaml:"windows"`
}

// Window 结构体表示一个 tmux 窗口的配置
type Window struct {
	Name  string   `yaml:"name"`
	Root  string   `yaml:"root"`
	Panes []string `yaml:"panes"`
}

func main() {
	// 读取配置文件
	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}

	// 解析配置文件
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	// 创建 tmux 会话
	sessionName := config.Name
	exec.Command("tmux", "new-session", "-d", "-s", sessionName).Run()

	for i, window := range config.Windows {
		// 对于每个窗口，创建一个新窗口
		windowName := fmt.Sprintf("%s:%d", sessionName, i)
		exec.Command("tmux", "new-window", "-t", windowName, "-n", window.Name).Run()

		// 设置窗口的根目录
		if window.Root != "" {
			exec.Command("tmux", "send-keys", "-t", windowName, "cd "+window.Root, "C-m").Run()
		}

		// 分割窗口并运行命令
		for _, paneCmd := range window.Panes {
			exec.Command("tmux", "split-window", "-t", windowName).Run()
			exec.Command("tmux", "send-keys", "-t", windowName, paneCmd, "C-m").Run()
		}

		// 最后一个 split-window 是多余的，所以关闭它
		// exec.Command("tmux", "kill-pane", "-t", windowName).Run()
	}

	// 附加到 tmux 会话
	// exec.Command("tmux", "attach-session", "-t", sessionName).Run()
	attachCmd := exec.Command("tmux", "attach-session", "-t", sessionName)
	attachCmd.Stdin = os.Stdin
	attachCmd.Stdout = os.Stdout
	attachCmd.Stderr = os.Stderr

	// 使用 Start 开始命令并等待其完成
	if err := attachCmd.Start(); err != nil {
		panic(err)
	}
	if err := attachCmd.Wait(); err != nil {
		panic(err)
	}
}

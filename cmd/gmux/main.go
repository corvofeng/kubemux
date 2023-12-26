package main

import (
	"bytes"
	"fmt"
	"gmux/lib"
	"gmux/lib/asset"

	log "github.com/sirupsen/logrus"

	"os"
	"os/exec"
	"text/template"
)

// Window 结构体表示一个 tmux 窗口的配置

func main() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
		DisableQuote:  true,
	})

	log.SetLevel(log.DebugLevel)
	// 读取配置文件
	// data, err := os.ReadFile("config.yaml")
	// if err != nil {
	// 	panic(err)
	// }

	// 解析配置文件
	// var config Config
	// err = yaml.Unmarshal(data, &config)
	// if err != nil {
	// 	panic(err)
	// }

	// 假设我们有以下配置
	config := lib.Config{
		Name: "my_session",
		Root: "~/GitRepo",
		Windows: []lib.Window{
			{
				Name: "editor",
				Root: "~/GitRepo",
				Panes: []lib.Pane{
					{Command: "vim"},
					{Command: "ls -alh"},
					{Command: "pwd"},
				},
			},
			{
				Name: "shell",
				Root: "/tmp",
				Panes: []lib.Pane{
					{Command: "htop"},
				},
			},
		},
	}

	config.Tmux = fmt.Sprintf("tmux -L %s", config.Name)

	funcMap := template.FuncMap{
		// The name "inc" is what the function will be called in the template text.
		"TmuxHasSession": lib.TmuxHasSession,
		"inc": func(i int) int {
			return i + 1
		},
	}
	// 创建一个新模板并解析模板字符串
	tmpl, err := template.New("bashScript").
		Funcs(funcMap).Parse(asset.BashScriptTemplate)
	if err != nil {
		log.Errorf("parsing: %s", err)
	}
	// 将配置数据应用于模板
	var script bytes.Buffer
	err = tmpl.Execute(&script, config)
	if err != nil {
		log.Errorf("execution: %s", err)
	}
	log.Debugf(script.String())

	// 将生成的脚本保存到临时文件
	tmpfile, err := os.CreateTemp("", "tmux-script-*.sh")
	if err != nil {
		log.Error(err)
	}

	defer os.Remove(tmpfile.Name()) // 清理临时文件

	if _, err := tmpfile.Write(script.Bytes()); err != nil {
		log.Error(err)
	}
	if err := tmpfile.Chmod(0755); err != nil { // 使脚本可执行
		log.Error(err)
	}
	tmpfile.Close()

	// 执行生成的脚本
	cmd := exec.Command(tmpfile.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return

	if err := cmd.Run(); err != nil {
		log.Error(err)
	}

	attachCmd := exec.Command("tmux", "-L", config.Name, "attach-session", "-t", config.Name)
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

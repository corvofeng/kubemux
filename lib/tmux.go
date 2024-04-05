package lib

import (
	"bytes"
	"gmux/lib/asset"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"
)

func ParseConfigPath(flagDirectory, flagProject string) string {

	if strings.HasPrefix(flagDirectory, "~/") {
		dirname, _ := os.UserHomeDir()
		flagDirectory = filepath.Join(dirname, flagDirectory[2:])
	}
	// if we have the full path for project
	if strings.HasPrefix(flagProject, ".") || strings.HasPrefix(flagProject, "/") {
		flagDirectory = ""
	}

	if strings.HasSuffix(flagProject, "yml") {
		flagProject = strings.TrimSuffix(flagProject, ".yml")
	}

	return path.Join(flagDirectory, flagProject+".yml")
}

func TmuxHasSession(sessionName string) bool {
	cmd := exec.Command("tmux", "-L", sessionName, "ls")
	out, err := cmd.Output()
	if err != nil {
		// Handle error if command execution fails
		// For simplicity, returning false in case of error
		return false
	}

	// Split the output into lines and check if sessionName exists
	outputLines := strings.Split(string(out), "\n")
	for _, line := range outputLines {
		if strings.HasPrefix(line, sessionName+":") {
			return true
		}
	}

	return false
}

// https://github.com/ruby/shellwords
// https://apidock.com/ruby/v2_5_5/Shellwords/shellescape
func shellescape(str string) string {
	str = regexp.MustCompile(`[^A-Za-z0-9_\-.,:+/@\n]`).ReplaceAllStringFunc(str, func(s string) string {
		return "\\" + s
	})
	str = strings.ReplaceAll(str, "\n", `'\n'`)
	if str == "" {
		return "''"
	}
	return str
}

func RunTmux(log log.FieldLogger, config *Config) {
	funcMap := template.FuncMap{
		"TmuxHasSession": TmuxHasSession,
		"Safe":           shellescape,
		"StringsJoin":    strings.Join,
		"Inc": func(i int) int {
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

	cmd := exec.Command(tmpfile.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Error(err)
	}

	var args []string
	if len(config.TmuxArgs) == 0 { // default to attach-session
		args = []string{"-L", config.Name, "attach-session", "-t", config.Name}
		args = append(args, config.TmuxArgs...)
	} else { // If we have custom args, use them
		args = []string{"-L", config.Name}
		args = append(args, config.TmuxArgs...)

	}

	attachCmd := exec.Command("tmux", args...)
	attachCmd.Stdin = os.Stdin
	attachCmd.Stdout = os.Stdout
	attachCmd.Stderr = os.Stderr

	// 使用 Start 开始命令并等待其完成
	if err := attachCmd.Start(); err != nil {
		panic(err)
	}
	if err := attachCmd.Wait(); err != nil {
		log.Error(err)
	}

}

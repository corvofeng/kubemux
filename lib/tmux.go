package lib

import (
	"bytes"
	"fmt"
	"kubemux/lib/asset"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"

	log "github.com/sirupsen/logrus"
)

func HasTmux() bool {
	return isCommandAvailable("tmux")
}

func GetConfigList(flagDirectory string) []string {
	if strings.HasPrefix(flagDirectory, "~/") {
		dirname, _ := os.UserHomeDir()
		flagDirectory = filepath.Join(dirname, flagDirectory[2:])
	}
	dirs, err := os.ReadDir(flagDirectory)
	if os.IsNotExist(err) {
		log.Error("Can't read", err)
		return []string{}
	}
	configList := []string{}
	for _, cfg := range dirs {
		if !strings.HasSuffix(cfg.Name(), ".yml") {
			continue
		}

		name := strings.TrimSuffix(cfg.Name(), ".yml")
		configList = append(
			configList,
			fmt.Sprintf("%s\t%s", name, name),
		)
	}

	return configList
}

func GetKubeConfigList() []string {
	homeDir, _ := os.UserHomeDir()
	kubeConfigDir := filepath.Join(homeDir, ".kube")
	dirs, err := os.ReadDir(kubeConfigDir)
	if os.IsNotExist(err) {
		log.Error("Can't read", err)
		return []string{}
	}
	configList := []string{}
	for _, cfg := range dirs {
		if cfg.IsDir() {
			continue
		}
		configList = append(
			configList,
			cfg.Name(),
		)
	}

	return configList
}

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

func ZellijHasSession(sessionName string) bool {
	cmd := exec.Command("zellij", "list-sessions", "-n")
	out, err := cmd.Output()
	if err != nil {
		// Handle error if command execution fails
		// For simplicity, returning false in case of error
		return false
	}

	// Split the output into lines and check if sessionName exists
	outputLines := strings.Split(string(out), "\n")
	for _, line := range outputLines {
		if strings.HasPrefix(line, sessionName) {
			return true
		}
	}

	return false
}

/*
The zellij may have exited session, it will block the create action, so we need to delete it.
Zellij list sessions output:

kubemux [Created 0s ago] (EXITED - attach to resurrect)
*/
func RemoveExitedZellijSession(sessionName string) {
	cmd := exec.Command("zellij", "list-sessions", "-n")
	out, err := cmd.Output()
	if err != nil {
		// Handle error if command execution fails
		// For simplicity, returning false in case of error
		log.Errorf("Can't list zellij sessions: %v", err)
	}

	outputLines := strings.Split(string(out), "\n")
	for _, line := range outputLines {
		if strings.HasPrefix(line, sessionName) {
			if strings.Contains(line, "(EXITED - attach to resurrect)") {
				log.Debug("Zellij session is exited, we need to delete it")
				cmd := exec.Command("zellij", "delete-session", sessionName)
				output, err := cmd.Output()
				log.Debugf("Run zellij delete session: %s %v", string(output), err)
			}
		}
	}
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

/*
*
  - ConfigCheck for tmux config
    We have such script template:
    {{$.Tmux}} split-window -c {{$window.Root}} -t {{$.Name}}:{{$winId}}
    We can't predict the program behavior if $window.Root is empty.
*/
func ConfigCheck(config *Config) error {
	if config.Root == "" {
		return fmt.Errorf("tmux Root is required")
	}
	for _, window := range config.Windows {
		if window.Root == "" {
			return fmt.Errorf("window root for %s is required", window.Name)
		}
	}
	return nil
}

func RenderCMDTemplate(
	baseTemplate string, config *Config, funcMap template.FuncMap,
	destFile *os.File,
) *exec.Cmd {

	tmpl, err := template.New("bashScript").
		Funcs(funcMap).Parse(baseTemplate)
	if err != nil {
		log.Errorf("parsing: %s", err)
	}
	var script bytes.Buffer
	err = tmpl.Execute(&script, config)
	if err != nil {
		log.Errorf("execution: %s", err)
	}
	log.Debugf("Start command: %s", script.String())

	if _, err := destFile.Write(script.Bytes()); err != nil {
		log.Error(err)
	}
	if err := destFile.Chmod(0755); err != nil {
		log.Error(err)
	}
	destFile.Close()
	return exec.Command(destFile.Name())
}

func RunTmux(config *Config) {
	// config.PlexerTool = KZellij // Debug
	// config.PlexerTool = KTmux   // Debug

	funcMap := template.FuncMap{
		"Safe":        shellescape,
		"StringsJoin": strings.Join,
		"Inc": func(i int) int {
			return i + 1
		},
	}
	if err := ConfigCheck(config); err != nil {
		log.Error(err)
		return
	}

	// how do we create a new session?
	startCmd, err := os.CreateTemp("", "kubemux-start-script-*.sh")
	if err != nil {
		log.Error(err)
	}
	defer os.Remove(startCmd.Name())

	// how do we attach to the session?
	attachCmd, err := os.CreateTemp("", "kubemux-attach-script-*.sh")
	if err != nil {
		log.Error(err)
	}
	defer os.Remove(attachCmd.Name())

	// how do we init the session layout
	prepareCmd, err := os.CreateTemp("", "kubemux-prepare-script-*.sh")
	if err != nil {
		log.Error(err)
	}
	defer os.Remove(prepareCmd.Name())

	runStartCmd := func() *exec.Cmd {
		startTemplate := asset.TmuxSessionCreateTemplate
		if config.PlexerTool == KZellij {
			startTemplate = asset.ZellijSessionCreateTemplate
		}
		return RenderCMDTemplate(startTemplate, config, funcMap, startCmd)
	}

	runAttachCmd := func() *exec.Cmd {
		attachTemplate := asset.TmuxSessionAttachTemplate
		if config.PlexerTool == KZellij {
			attachTemplate = asset.ZellijSessionAttachTemplate
		}
		return RenderCMDTemplate(attachTemplate, config, funcMap, attachCmd)
	}

	runPrepareCmd := func() *exec.Cmd {
		prepareTemplate := asset.TmuxSessionPrepareTemplate
		if config.PlexerTool == KZellij {
			prepareTemplate = asset.ZellijSessionPrepareTemplate
		}
		return RenderCMDTemplate(prepareTemplate, config, funcMap, prepareCmd)
	}

	var plexerCmd *exec.Cmd
	// After we create the session, we need to create related window and panes for user.
	// since the zellij can't create window and panes if no client attached,
	// we need to create them in a new goroutine, after we attach the window.
	// I also modify the code for tmux.
	var needPrepareSession bool = false

	if config.PlexerTool == KTmux {
		if !TmuxHasSession(config.Name) {
			plexerCmd = runStartCmd()
			needPrepareSession = true
		} else {
			plexerCmd = runAttachCmd()
		}

	} else if config.PlexerTool == KZellij {
		RemoveExitedZellijSession(config.Name) // Ugly, but we need to check it.

		if !ZellijHasSession(config.Name) {
			plexerCmd = runStartCmd()
			needPrepareSession = true
		} else {
			plexerCmd = runAttachCmd()
		}
	}

	if plexerCmd == nil {
		log.Error("Can't find plexer tool")
		return
	}

	plexerCmd.Stdin = os.Stdin
	plexerCmd.Stdout = os.Stdout
	plexerCmd.Stderr = os.Stderr

	// 使用 Start 开始命令并等待其完成
	if err := plexerCmd.Start(); err != nil {
		panic(err)
	}

	if needPrepareSession {
		go func() {
			time.Sleep(2 * time.Second) // wait for session created
			log.Debug("Create prepare session command")
			prepareCmd := runPrepareCmd()
			prepareCmd.Stdout = os.Stdout
			prepareCmd.Stderr = os.Stderr
			if err := prepareCmd.Run(); err != nil {
				log.Error("Can't run prepare command", err)
			}
		}()
	}
	log.Debugf("Waiting for command to finish...")

	if err := plexerCmd.Wait(); err != nil {
		log.Error(err)
	}
}

func TryToAttachTmux(log log.FieldLogger, config *Config) *exec.Cmd {
	var args []string
	if len(config.TmuxArgs) == 0 { // default to attach-session
		args = []string{"-L", config.Name, "attach-session", "-t", config.Name}
		args = append(args, config.TmuxArgs...)
	} else { // If we have custom args, use them
		args = []string{"-L", config.Name}
		args = append(args, config.TmuxArgs...)
	}

	log.Debugf("Run tmux %s", args)
	attachCmd := exec.Command("tmux", args...)
	return attachCmd
}

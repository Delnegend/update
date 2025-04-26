package app

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	c "update/app/internal"
)

type AppName string

const (
	AdvancedRenamer    AppName = "Advanced Renamer"
	Caddy              AppName = "Caddy"
	Everything         AppName = "Everything"
	FirewallAppBlocker AppName = "FirewallAppBlocker"
	GoCompiler         AppName = "Go"
	HandBrake          AppName = "HandBrake"
	HTML2Markdown      AppName = "html2markdown"
	Just               AppName = "just"
	Lazygit            AppName = "lazygit"
	LibJXL             AppName = "libjxl"
	MadVR              AppName = "madVR"
	Nrr                AppName = "nrr"
	NvidiaICAT         AppName = "NvidiaICAT"
	Onefetch           AppName = "onefetch"
	Photoshop          AppName = "Photoshop"
	PlatformTools      AppName = "Platform Tool"
	Python             AppName = "Python"
	RevoUninstaller    AppName = "Revo Uninstaller"
	RsRPC              AppName = "rsRPC"
	Upx                AppName = "upx"
)

type Checker func() c.CheckResult

type AppInfo struct {
	Name       AppName
	Exec       string
	rawChecker func(func() (string, error)) c.CheckResult
}

type CheckResult struct {
	Name AppName
	c.CheckResult
}

var apps = []AppInfo{
	{AdvancedRenamer, "arenc", c.AdvancedRenamer},
	{Caddy, "caddy", c.Caddy},
	{Everything, "Everything", c.Everything},
	{FirewallAppBlocker, "Fab_x64", c.FirewallAppBlocker},
	{GoCompiler, "go", c.GoCompiler},
	{HandBrake, "HandBrake", c.Handbrake},
	{HTML2Markdown, "html2markdown", c.HTML2Markdown},
	{Just, "just", c.Just},
	{Lazygit, "lazygit", c.Lazygit},
	{LibJXL, "cjxl", c.LibJXL},
	{MadVR, "", c.MadVR},
	{Nrr, "nrr", c.Nrr},
	{NvidiaICAT, "ICAT", c.NvidiaICAT},
	{Onefetch, "onefetch", c.Onefetch},
	{Photoshop, "Photoshop", c.Photoshop},
	{PlatformTools, "adb", c.PlatformTools},
	{Python, "python", c.Python},
	{RevoUninstaller, "RevoUnPro", c.RevoUninstaller},
	{RsRPC, "rsrpc-cli", c.RsRPC},
	{Upx, "upx", c.Upx},
}

type Apps struct {
	Apps          []AppInfo
	execOverrides map[AppName]string
}

func InitApps() Apps {
	execOverrides := func() map[AppName]string {
		acc := make(map[AppName]string)

		exePath, err := os.Executable()
		if err != nil {
			return acc
		}
		exeDir := strings.TrimRightFunc(exePath, func(r rune) bool { return r != '\\' && r != '/' })
		updateTxtPath := exeDir + string(os.PathSeparator) + "update.txt"

		file, err := os.Open(updateTxtPath)
		if err != nil {
			return acc
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.Split(line, "=")
			if len(parts) != 2 {
				panic(fmt.Errorf("invalid line in update.txt: %s", line))
			}
			acc[AppName(strings.TrimSpace(parts[0]))] = strings.TrimSpace(parts[1])
		}
		return acc
	}()

	return Apps{apps, execOverrides}
}

func (a *Apps) CheckApp(appName AppName) CheckResult {
	appName = strings.ToLower(appName)

	for _, app := range a.Apps {
		if strings.ToLower(app.Name) == appName {
			return CheckResult{
				app.Name,
				app.checker(func() (string, error) {
					if _, err := os.Stat(a.execOverrides[appName]); err == nil {
						return a.execOverrides[appName], nil
					}

					if _, err := exec.LookPath(app.Exec); err == nil {
						return app.Exec, nil
					}

					return "", fmt.Errorf("can't find executable")
				}),
			}
		}
	}
	panic("app not found")
}

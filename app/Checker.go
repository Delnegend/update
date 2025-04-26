package app

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	c "update/app/internal"
)

type AppName = string

type Checker func() c.CheckResult

type AppInfo struct {
	Name    AppName
	Exec    string
	checker func(func() (string, error)) c.CheckResult
}

type CheckResult struct {
	Name AppName
	c.CheckResult
}

var apps = []AppInfo{
	{"Advanced Renamer", "arenc", c.AdvancedRenamer},
	{"Caddy", "caddy", c.Caddy},
	{"EqualizerAPO", "EqualizerAPO", c.EqualizerAPO},
	{"Everything", "Everything", c.Everything},
	{"Firewall App Blocker", "Fab_x64", c.FirewallAppBlocker},
	{"Git", "git", c.Git},
	{"Go", "go", c.GoCompiler},
	{"HandBrake", "HandBrake", c.Handbrake},
	{"html2markdown", "html2markdown", c.HTML2Markdown},
	{"just", "just", c.Just},
	{"LAVFilter", "", c.LAVFilter},
	{"Lazygit", "lazygit", c.Lazygit},
	{"LibJXL", "cjxl", c.LibJXL},
	{"madVR", "", c.MadVR},
	{"NanaZip", "NanaZip", c.NanaZip},
	{"NeatDM", "NeatDM", c.NeatDM},
	{"Node", "node", c.Node},
	{"nrr", "nrr", c.Nrr},
	{"NvidiaICAT", "ICAT", c.NvidiaICAT},
	{"Onefetch", "onefetch", c.Onefetch},
	{"Photoshop", "Photoshop", c.Photoshop},
	{"Platform Tools", "adb", c.PlatformTools},
	{"Pnpm", "pnpm", c.Pnpm},
	{"Python", "python", c.Python},
	{"Revo Uninstaller", "RevoUnPro", c.RevoUninstaller},
	{"RsRPC", "rsrpc-cli", c.RsRPC},
	{"Upx", "upx", c.Upx},
	{"Volta", "volta", c.Volta},
	{"ZeroTier", "zerotier-cli", c.ZeroTier},
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
			file, err = os.Open("update.txt")
			if err != nil {
				return acc
			}
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
			acc[strings.ToLower(strings.TrimSpace(parts[0]))] = strings.TrimSpace(parts[1])
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
	panic(fmt.Sprintf("app not found: %s", appName))
}

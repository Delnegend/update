package check

import (
	"fmt"
	"regexp"
	"strings"
	"update/utils"
)

var revoUninstallerLatestVerRe = regexp.MustCompile(`Latest Version: ([\d\.]+)`)

func RevoUninstaller(getExec func() (string, error)) CheckResult {
	currVer, err := func() (string, error) {
		execPath, err := getExec()
		if err != nil {
			return "", err
		}
		currVer, err := utils.GetExeVersion(execPath)
		if err != nil {
			return "", fmt.Errorf("can't get executable version: %v", err)
		}
		return strings.TrimSuffix(currVer, ".0"), nil
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get current version: %v", err)}
	}

	latestVer, err := func() (string, error) {
		content, err := utils.Fetch("https://www.revouninstaller.com/products/revo-uninstaller-pro/")
		if err != nil {
			return "", fmt.Errorf("can't fetch the homepage: %v", err)
		}
		match := revoUninstallerLatestVerRe.FindStringSubmatch(string(content))
		if match == nil {
			return "", fmt.Errorf("unexpected website content")
		}
		return strings.TrimSpace(match[1]), nil
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get latest version: %v", err)}
	}

	return CheckResult{
		CurrentVersion:  currVer,
		LatestVersion:   latestVer,
		DirectURLAlive:  false,
		DownloadPageURL: "https://www.revouninstaller.com/products/revo-uninstaller-pro/",
	}
}

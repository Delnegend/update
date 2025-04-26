package check

import (
	"fmt"
	"regexp"
	"strings"
	"update/utils"
)

var photoshopLatestVerRe = regexp.MustCompile(`The current version of Photoshop.+?\(desktop\) is.<b>([\d\.]+)<\/b>`)

func Photoshop(getExec func() (string, error)) CheckResult {
	currVer, err := func() (string, error) {
		execPath, err := getExec()
		if err != nil {
			return "", err
		}

		currVer, err := utils.GetExeVersion(execPath)
		if err != nil {
			return "", fmt.Errorf("can't get executable version: %v", err)
		}
		slices := strings.Split(currVer, ".")
		return strings.Join(slices[:len(slices)-2], "."), nil
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get current version: %v", err)}
	}

	latestVer, err := func() (string, error) {
		webpageContent, err := utils.Fetch("https://helpx.adobe.com/photoshop/kb/uptodate.html")
		if err != nil {
			return "", fmt.Errorf("can't fetch webpage: %v", err)
		}
		match := photoshopLatestVerRe.FindStringSubmatch(string(webpageContent))
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
		DownloadPageURL: "https://www.adobe.com/creativecloud/desktop-app.html",
	}
}

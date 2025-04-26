package main

import (
	"fmt"
	"strings"
	"sync"
	"update/app"
)

var (
	RED       = "\033[91m"
	GREEN     = "\033[92m"
	BLUE      = "\033[94m"
	YELLOW    = "\033[93m"
	END       = "\033[0m"
	UNDERLINE = "\033[4m"
)

func main() {
	apps := app.InitApps()
	results := make(chan app.CheckResult)
	defer close(results)
	done := make(chan struct{})
	defer close(done)
	var wg sync.WaitGroup

	for _, appInfo := range apps.Apps {
		wg.Add(1)
		go func(appName app.AppName) {
			defer wg.Done()
			results <- apps.CheckApp(appName)
		}(appInfo.Name)
	}

	go func() {
		wg.Wait()
		done <- struct{}{}
	}()

scoped:
	for {
		select {
		case r := <-results:
			var sb strings.Builder

			sb.WriteString(r.Name + " ")
			if r.Error != nil {
				sb.WriteString("⛔  " + r.Error.Error())
				fmt.Println(sb.String())
				continue
			}
			if r.CurrentVersion == "" {
				sb.WriteString("⚠️  current version is empty")
				continue
			}
			if r.LatestVersion == "" {
				sb.WriteString("⚠️  latest version is empty")
				fmt.Println(sb.String())
				continue
			}
			if r.CurrentVersion > r.LatestVersion {
				sb.WriteString("❓  installed version is newer")
				fmt.Println(sb.String())
				continue
			}

			if r.CurrentVersion == r.LatestVersion {
				sb.WriteString("✔️")
				fmt.Println(sb.String())
				continue
			}
			sb.WriteString(YELLOW)
			sb.WriteString(r.CurrentVersion)
			sb.WriteString(" -> ")
			sb.WriteString(r.LatestVersion)
			sb.WriteString(END)
			if r.DirectURL == "" || !r.DirectURLAlive {
				sb.WriteString(" 🏠  ")
				sb.WriteString(BLUE)
				sb.WriteString(UNDERLINE)
				sb.WriteString(r.DownloadPageURL)
			} else {
				sb.WriteString(" 🔗  ")
				sb.WriteString(GREEN)
				sb.WriteString(UNDERLINE)
				sb.WriteString(r.DirectURL)
			}
			sb.WriteString(END)
			sb.WriteString(END)
			fmt.Println(sb.String())
		case <-done:
			break scoped
		}
	}
}

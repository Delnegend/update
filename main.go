package main

import (
	"fmt"
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
			results <- apps.CheckApp(appInfo.Name)
		}(appInfo.Name)
	}

	go func() {
		wg.Wait()
		done <- struct{}{}
	}()

	fmt.Println()

scoped:
	for {
		select {
		case r := <-results:
			fmt.Printf("%s: ", r.Name)
			if r.Error != nil {
				fmt.Printf("%v\n\n", r.Error)
				continue
			}
			if r.CurrentVersion == "" {
				fmt.Printf("current version is empty\n\n")
				continue
			}
			if r.LatestVersion == "" {
				fmt.Printf("latest version is empty\n\n")
				continue
			}
			if r.CurrentVersion == r.LatestVersion {
				fmt.Printf("up to date\n\n")
				continue
			}
			fmt.Printf("%s%s -> %s%s\n", YELLOW, r.CurrentVersion, r.LatestVersion, END)
			if r.DirectURL == "" || !r.DirectURLAlive {
				fmt.Printf("  - 🏠 %s%s%s%s%s\n\n", BLUE, UNDERLINE, r.DownloadPageURL, END, END)
				continue
			}
			fmt.Printf("  - 🔗 %s%s%s%s%s\n\n", GREEN, UNDERLINE, r.DirectURL, END, END)
		case <-done:
			break scoped
		}
	}
}

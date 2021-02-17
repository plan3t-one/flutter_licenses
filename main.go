package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"golang.org/x/sync/semaphore"

	"github.com/cheggaaa/pb/v3"
	"github.com/olekukonko/tablewriter"
)

var fullReport = flag.Bool("full-report", false, "generate a full report where each dependency is listed with their respective license")
var format = flag.String("format", "table", "output format, either 'table' (default) or 'csv'")

func main() {
	flag.Parse()
	err := run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, args []string) error {
	pathToLockFile := args[len(args)-1]

	lockFileContent, err := ioutil.ReadFile(pathToLockFile)
	if err != nil {
		return err
	}

	lockFile, err := parseLockFile(lockFileContent)
	if err != nil {
		return err
	}

	licenseCount := map[License]int{}
	licenses := map[string]License{}
	l := sync.Mutex{}

	bar := pb.StartNew(len(lockFile.Packages))
	bar.SetWriter(os.Stderr)

	sem := semaphore.NewWeighted(4)
	wg := sync.WaitGroup{}

	for _, p := range lockFile.Packages {
		wg.Add(1)
		err := sem.Acquire(ctx, 1)
		if err != nil {
			return err
		}
		go func(name string) {
			license, err := getLicense(ctx, name)
			if err != nil {
				log.Printf("failed to get license for %s: %s", name, err)
				license = "ERROR"
			}

			l.Lock()
			licenses[name] = license
			licenseCount[license] = licenseCount[license] + 1
			l.Unlock()

			bar.Increment()
			sem.Release(1)
			wg.Done()
		}(p.Description.Name)
	}
	wg.Wait()
	bar.Finish()

	if *format == "table" {
		if *fullReport {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Dependency", "Version", "License"})
			for name, l := range licenses {
				table.Append([]string{name, lockFile.Packages[name].Version, string(l)})
			}
			table.Render()
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"License", "# Found"})
			table.SetFooter([]string{"Total", fmt.Sprintf("%d", len(licenseCount))})
			for l, c := range licenseCount {
				table.Append([]string{string(l), fmt.Sprintf("%d", c)})
			}

			table.Render()
		}
	} else if *format == "csv" {
		if *fullReport {
			fmt.Println("Dependency,Version,License")
			for name, l := range licenses {
				fmt.Printf("%s,%s,%s\n", name, lockFile.Packages[name].Version, string(l))
			}
		} else {
			fmt.Println("License, # Found")
			for l, c := range licenseCount {
				fmt.Printf("%s,%d\n", string(l), c)
			}
		}
	}

	return nil
}

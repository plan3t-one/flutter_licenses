package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"golang.org/x/sync/semaphore"

	"github.com/cheggaaa/pb/v3"
	"github.com/olekukonko/tablewriter"
)

func main() {
	err := run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, args []string) error {
	pathToLockFile := args[1]

	lockFileContent, err := ioutil.ReadFile(pathToLockFile)
	if err != nil {
		return err
	}

	packages, err := parsePackages(lockFileContent)
	if err != nil {
		return err
	}

	licenses := map[License]int{}
	l := sync.Mutex{}

	bar := pb.StartNew(len(packages))

	sem := semaphore.NewWeighted(4)
	wg := sync.WaitGroup{}

	for _, name := range packages {
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
			licenses[license] = licenses[license] + 1
			l.Unlock()

			bar.Increment()
			sem.Release(1)
			wg.Done()
		}(name)
	}
	wg.Wait()
	bar.Finish()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"License", "# Found"})
	table.SetFooter([]string{"Total", fmt.Sprintf("%d", len(packages))})
	for l, c := range licenses {
		table.Append([]string{string(l), fmt.Sprintf("%d", c)})
	}

	table.Render()

	return nil
}

package search

import (
	"bufio"
	"context"
	"os"
	"strings"
	"sync"
)
type Result struct {
	Phrase 	string
	Line 		string
	LineNum	int64
	ColNum	int64
}
func All(ctx context.Context, phrase string, files []string) <-chan []Result{
	
	ch := make(chan []Result)
	wg := sync.WaitGroup{}

	ctx, cancel := context.WithCancel(ctx)

	for i := 0; i < len(files); i++ {

		wg.Add(1)
		go func(ctx context.Context, file string, ch chan<- []Result, index int) {
			defer wg.Done()
			lines := []string{}
			resultMain := []Result{}
			f, err := os.Open(file)
			if err != nil {
				return
			}

			defer f.Close()

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}

			for j := 0; j < len(lines); j++ {
				if strings.Contains(lines[j], phrase) {
					result := Result {
						Phrase: 	phrase,
						Line:		lines[j],
						LineNum: int64(j + 1),
						ColNum: 	int64(strings.Index(lines[j], phrase)) + 1,
					}
					resultMain = append(resultMain, result)
				}
			}
			if len(resultMain) > 0 {
				ch <- resultMain
			}

		}(ctx, files[i], ch, i)
	}
	go func() {
		defer close(ch)
		wg.Wait()
	}()
	
	cancel()
	return ch
}

func Any(ctx context.Context, phrase string, files []string) <-chan Result {
	ch := make(chan Result)
	wg := sync.WaitGroup{}

	ctx, cancel := context.WithCancel(ctx)

		wg.Add(1)
		go func(ctx context.Context, file string, ch chan<- Result) {
			defer wg.Done()
			lines := []string{}
			resultMain := []Result{}
			f, err := os.Open(file)
			if err != nil {
				return
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}

			for j := 0; j < len(lines); j++ {
				if strings.Contains(lines[j], phrase) {
					result := Result {
						Phrase: 	phrase,
						Line:		lines[j],
						LineNum: int64(j + 1),
						ColNum: 	int64(strings.Index(lines[j], phrase)) + 1,	
					}
					resultMain = append(resultMain, result)
				}
			}
			if len(resultMain) > 0 {
				ch <- resultMain[0]
			}
		}(ctx, files[0], ch)

	go func() {
		defer close(ch)
		wg.Wait()
	}()

	cancel()
	return ch
}
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/grailbio/bigslice"
	"github.com/grailbio/bigslice/sliceconfig"
)

var wordCount = bigslice.Func(func(url string) bigslice.Slice {
	slice := bigslice.ScanReader(8, func() (io.ReadCloser, error) {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("get %v: %v", url, resp.Status)
		}
		return resp.Body, nil
	})
	slice = bigslice.Flatmap(slice, strings.Fields)
	slice = bigslice.Map(slice, func(token string) (string, int) {
		return token, 1
	})
	slice = bigslice.Reduce(slice, func(a, e int) int {
		return a + e
	})
	return slice
})

const shakespeare = "https://ocw.mit.edu/ans7870/6" +
	"/6.006/s08/lecturenotes/files/t8.shakespeare.txt"

func main() {
	sess := sliceconfig.Parse()
	defer sess.Shutdown()

	ctx := context.Background()
	tokens, err := sess.Run(ctx, wordCount, shakespeare)
	if err != nil {
		log.Fatal(err)
	}
	scanner := tokens.Scanner()
	defer scanner.Close()
	type counted struct {
		token string
		count int
	}
	var (
		token  string
		count  int
		counts []counted
	)
	for scanner.Scan(ctx, &token, &count) {
		counts = append(counts, counted{token, count})
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	sort.Slice(counts, func(i, j int) bool {
		return counts[i].count > counts[j].count
	})
	if len(counts) > 10 {
		counts = counts[:10]
	}
	for _, count := range counts {
		fmt.Println(count.token, count.count)
	}
}

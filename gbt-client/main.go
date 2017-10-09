package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/valyala/gorpc"
)

func main() {
	pairs := newPairs(total, msgSize)

	fmt.Fprintf(os.Stderr, "creating client for %s server\n", server)
	c := gorpc.NewTCPClient(server)

	fmt.Fprintf(os.Stderr, "starting client\n")
	c.Start()
	defer c.Stop()

	fmt.Fprintf(os.Stderr, "sending messages\n")
	count := len(pairs)
	for count > 0 {
		n, err := sendBatch(c, len(pairs)-count, limit, pairs)
		if err != nil {
			panic(err)
		}

		count -= n
	}

	dump(pairs, "")
}

type pair struct {
	req []byte

	sent time.Time
	recv *time.Time
	dup  int
}

func newPairs(n, size int) []*pair {
	out := make([]*pair, n)

	if size > 0 {
		fmt.Fprintf(os.Stderr, "making messages to send:\n")
	}

	fold := 3
	last := len(out) - 1

	for i := range out {
		if size > 0 {
			buf := make([]byte, size)
			for i := 6; i < len(buf); i++ {
				buf[i] = 0xaa
			}

			if i < fold || i == fold && fold >= last {
				fmt.Fprintf(os.Stderr, "\t%d: % x\n", i, buf)
			} else if i == fold {
				fmt.Fprintf(os.Stderr, "\t%d: ...\n", i)
			}

			out[i] = &pair{req: buf}
		} else {
			out[i] = &pair{}
		}
	}

	if size > 0 && fold < last {
		fmt.Fprintf(os.Stderr, "\t%d: % x\n", last, out[last].req)
	}

	return out
}

func sendBatch(c *gorpc.Client, start, chunk int, pairs []*pair) (int, error) {
	if start >= len(pairs) || chunk <= 0 {
		return 0, nil
	}

	b := c.NewBatch()

	end := start + chunk
	if end >= len(pairs) {
		end = len(pairs)
	}

	var wg sync.WaitGroup
	for _, p := range pairs[start:end] {
		p.sent = time.Now()
		wg.Add(1)
		go func(p *pair, r *gorpc.BatchResult) {
			defer wg.Done()
			<-r.Done

			t := time.Now()
			p.recv = &t
		}(p, b.Add(p.req))
	}

	// fmt.Fprintf(os.Stderr, "sending requests %d - %d\n", start, end)
	err := b.Call()
	if err != nil {
		return 0, err
	}

	wg.Wait()

	return end - start, nil
}

package main

import (
	"bufio"
	"io"
	"math/rand"
	"os/exec"
	"time"

	"github.com/creack/pty"
	"go.coder.com/flog"
)

// randchar generates a random alphabetic character
func randchar() byte {
	return 'A' + byte(rand.Intn(26))
}


type timings struct {
	connect timer
	// input contains a timer for each character inputted.
	input []timer
}

func test(cmd *exec.Cmd, iterations int) timings {
	t := timings{
		connect: startTimer(),
	}

	tty, err := pty.Start(cmd)
	if err != nil {
		flog.Fatal("start tty: ", err)
	}
	defer tty.Close()


	rd := bufio.NewReader(tty)
	const shellIndicator = '$'
	prelude, err := rd.ReadBytes(shellIndicator)
	if err != nil {
		flog.Fatal("couldn't find indication that shell started (no %q in %s): %+v", shellIndicator, prelude, err)
	}
	t.connect.end()

	// Discard any terminal bullshit.
	time.Sleep(time.Millisecond*500)
	rd.Discard(rd.Buffered())

	for i := 0; i < iterations; i++ {
		ct := startTimer()
		c := randchar()
		io.WriteString(tty, string(c))
		got, err := rd.ReadByte()
		if err != nil {
			flog.Fatal("read back byte: %v", err)
		}
		if got != c {
			flog.Fatal("sent %q, got %q back", c, got)
		}
		ct.end()
		t.input = append(t.input, ct)
	}

	return t
}

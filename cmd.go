package main

import (
	"bufio"
	"io"
	"math/rand"
	"os/exec"
	"time"

	"github.com/coder/flog"
	"github.com/creack/pty"
	"github.com/thanhpk/randstr"
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

func readUntil(rd io.Reader, token string) error {
	br := bufio.NewReader(rd)
	for i := range token {
		_, err := br.ReadBytes(token[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func test(cmd *exec.Cmd, iterations int) timings {
	tty, err := pty.Start(cmd)
	if err != nil {
		flog.Fatalf("start tty: ", err)
	}
	defer tty.Close()

	// Discard any terminal bullshit.
	// time.Sleep(time.Millisecond * 500)

	// Reset buffer
	// rd.Discard(rd.Buffered())

	// We wait for connectNonce to echo bank to indicate the TTY has connected.
	t := timings{
		connect: startTimer(),
	}

	tty.SetReadDeadline(time.Now().Add(10 * time.Second))
	// Detect the other end of the TTY connecting. If we send before,
	// we may just process the local echo.
	tty.Read(make([]byte, 1))

	// Skip the bytes that represent the shell line (e.g. "$").
	connectNonce := randstr.String(4)
	tty.WriteString(connectNonce)
	err = readUntil(tty, connectNonce)
	if err != nil {
		flog.Fatalf("did not find nonce (%q) echoed back", connectNonce)
	}
	t.connect.end()

	gotBuf := make([]byte, 1)
	for i := 0; i < iterations; i++ {
		// rd.Reset()
		tty.SetReadDeadline(time.Now().Add(10 * time.Second))

		ct := startTimer()
		c := randchar()
		io.WriteString(tty, string(c))
		_, err := tty.Read(gotBuf)
		if err != nil {
			flog.Fatalf("read back byte: %v", err)
		}
		got := gotBuf[0]
		if got != c {
			flog.Fatalf("sent %q, got %q back", c, got)
		}
		ct.end()
		t.input = append(t.input, ct)
	}

	return t
}

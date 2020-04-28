package main

import "time"

type timer struct {
	startedAt, endAt time.Time
}

func (t *timer) end() {
	t.endAt = time.Now()
}

func (t *timer) took() time.Duration {
	return t.endAt.Sub(t.startedAt)
}

func startTimer() timer {
	return timer{startedAt: time.Now()}
}

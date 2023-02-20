package terminal

import "time"

var spinnerIcons = []string{
	"⢿", "⣻", "⣽", "⣾", "⣷", "⣯", "⣟", "⡿",
}

type Spinner struct {
	state int
	Tick  <-chan time.Time
}

func NewSpinner(tickInterval time.Duration) *Spinner {
	ticker := time.NewTicker(tickInterval)
	return &Spinner{
		Tick:ticker.C,
	}
}

func (s *Spinner) Spin() {
	s.state++
	s.state = s.state % len(spinnerIcons)
}

func (s *Spinner) String() string {
	return spinnerIcons[s.state]
}

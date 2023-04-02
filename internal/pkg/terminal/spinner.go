package terminal

import "time"

// spinnerIcons is a list of spinner icons
var spinnerIcons = []string{
	"⢿", "⣻", "⣽", "⣾", "⣷", "⣯", "⣟", "⡿",
}

// Spinner is a struct for a terminal spinner
type Spinner struct {
	state int
	Tick  <-chan time.Time
}

// NewSpinner returns a new Spinner
func NewSpinner(tickInterval time.Duration) *Spinner {
	ticker := time.NewTicker(tickInterval)
	return &Spinner{
		Tick: ticker.C,
	}
}

// Spin updates the spinner state
func (s *Spinner) Spin() {
	s.state++
	s.state = s.state % len(spinnerIcons)
}

// String returns the current spinner icon
func (s *Spinner) String() string {
	return spinnerIcons[s.state]
}

package terminal

import (
	"math"
	"strings"
)

var (
	fullBlock     = "█"
	partialBlocks = []string{
		"", "▏", "▎", "▍", "▌", "▋", "▊", "▉",
	}
)

type ProgressBar struct {
	increments      int
	totalIncrements int
	width           int
	progressBar     string
	barColor        string
	resetColor      string
}

func NewProgressBar(foregroundColor, backgroundColor, resetColor string) *ProgressBar {
	progressBarColor := foregroundColor + backgroundColor
	pb := &ProgressBar{
		barColor:        progressBarColor,
		resetColor:      resetColor,
		totalIncrements: 1,
		width:           5,
	}
	pb.updateState()
	return pb
}

func floor(x float64) int {
	return int(math.Floor(x))
}

func (p *ProgressBar) Increment(i int) {
	if p.increments >= p.totalIncrements {
		return
	}
	p.increments += i
	if p.increments > p.totalIncrements {
		p.increments = p.totalIncrements
	}
	p.updateState()
}

func (p *ProgressBar) updateState() {
	incrementPercent := float64(p.increments) / float64(p.totalIncrements)
	nbBlocksPercent := incrementPercent * float64(p.width)
	nbFullBlocks := floor(nbBlocksPercent)
	remainder := nbBlocksPercent - float64(nbFullBlocks)
	partialBlockIndex := floor(remainder * float64(len(partialBlocks)))
	fullBlocks := strings.Repeat(fullBlock, nbFullBlocks)
	partialBlock := partialBlocks[partialBlockIndex]
	nbPartialBlocks := 1
	if len(partialBlock) == 0 {
		nbPartialBlocks = 0
	}
	nbEmptyBlocks := p.width - nbFullBlocks - nbPartialBlocks
	emptyBlocks := strings.Repeat(" ", nbEmptyBlocks)
	p.progressBar = p.barColor + fullBlocks + partialBlock + emptyBlocks + p.resetColor
}

func (p *ProgressBar) SetTotalIncrements(increments int) {
	p.totalIncrements = increments
	p.updateState()
}

func (p *ProgressBar) SetWidth(width int) {
	p.width = width
	p.updateState()
}

func (p *ProgressBar) String() string {
	return p.progressBar
}

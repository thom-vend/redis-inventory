package adapter

import (
	"io"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/mattn/go-isatty"
)

// ProgressWriter abstraction of progress writer
type ProgressWriter interface {
	// Start initiates progress writing progress, if total is unknown should be zero
	Start(total int64)
	// Increment increments progress
	Increment()
	// Stop labels progress as finished and stops updating progress
	Stop()
}

// NewPrettyProgressWriter creates PrettyProgressWriter
func NewPrettyProgressWriter(output io.Writer) *PrettyProgressWriter {
	p := &PrettyProgressWriter{pw: progress.NewWriter()}
	p.init(output)

	return p
}

// PrettyProgressWriter progress writer using go-pretty/progress library
type PrettyProgressWriter struct {
	pw      progress.Writer
	tracker Tracker
}

// Tracker is abstraction over libraries "Tracker" struct
type Tracker interface {
	Increment(value int64)
	MarkAsDone()
}

func (p *PrettyProgressWriter) init(output io.Writer) {
	p.pw.SetAutoStop(false)
	p.pw.SetTrackerLength(50)
	p.pw.ShowETA(true)
	p.pw.ShowOverallTracker(false)
	p.pw.ShowTime(true)
	p.pw.ShowTracker(true)
	p.pw.ShowValue(true)
	p.pw.SetMessageWidth(13)
	p.pw.SetNumTrackersExpected(1)
	p.pw.SetSortBy(progress.SortByPercentDsc)
	p.pw.SetStyle(progress.StyleDefault)
	p.pw.SetTrackerPosition(progress.PositionRight)
	p.pw.Style().Options.PercentFormat = "%4.1f%%"

	if p.isTerminal() {
		p.pw.SetUpdateFrequency(time.Millisecond * 100)
		p.pw.Style().Colors = progress.StyleColorsExample
	} else {
		p.pw.SetUpdateFrequency(time.Second * 30)
		p.pw.SetStyle(progress.Style{
			Name:   "log-friendly",
			Chars:  progress.StyleCharsDefault,
			Colors: progress.StyleColors{},
			Options: progress.StyleOptions{
				Separator: " ",
			},
		})
	}
	p.pw.SetOutputWriter(output)
}

func (p *PrettyProgressWriter) isTerminal() bool {
	return isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
}

// Start initiates progress writing progress, if total is unknown should be zero
func (p *PrettyProgressWriter) Start(total int64) {
	scanningKeysTracker := &progress.Tracker{Message: "Scanning keys", Total: total, Units: progress.UnitsDefault}
	p.pw.AppendTracker(scanningKeysTracker)
	p.tracker = scanningKeysTracker

	go p.pw.Render()
}

// Increment increments progress
func (p *PrettyProgressWriter) Increment() {
	p.tracker.Increment(1)
}

// Stop labels progress as finished and stops updating progress
func (p *PrettyProgressWriter) Stop() {
	p.tracker.MarkAsDone()
	p.pw.Stop()
}

// Originally created by Brian Downs
// https://github.com/briandowns/spinner
//
// Stripped down and modified to work better for Home Assistant CLI
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package spinner is a simple package to add a spinner / progress indicator to any terminal application.
package spinner

import (
	"fmt"
	"io"
	"sync"
	"time"
	"unicode/utf8"
)

// Spinner struct to hold the provided options
type Spinner struct {
	Delay      time.Duration // Delay is the speed of the indicator
	chars      []string      // chars holds the chosen character set
	Prefix     string        // Prefix is the text prepended to the indicator
	Suffix     string        // Suffix is the text appended to the indicator
	FinalMSG   string        // string displayed after Stop() is called
	lastOutput string        // last character(set) written
	lock       *sync.RWMutex //
	Writer     io.Writer     // to make testing better, exported so users have access
	active     bool          // active holds the state of the spinner
	stopChan   chan struct{} // stopChan is a channel used to stop the indicator
	HideCursor bool          // hideCursor determines if the cursor is visible
}

// New provides a pointer to an instance of Spinner with the supplied options
func New(cs []string, d time.Duration, options ...Option) *Spinner {
	s := &Spinner{
		Delay:    d,
		chars:    cs,
		lock:     &sync.RWMutex{},
		Writer:   io.Discard,
		active:   false,
		stopChan: make(chan struct{}, 1),
	}

	for _, option := range options {
		option(s)
	}
	return s
}

// Option is a function that takes a spinner and applies
// a given configuration
type Option func(*Spinner)

// Options contains fields to configure the spinner
type Options struct {
	Suffix     string
	FinalMSG   string
	HideCursor bool
}

// WithFinalMSG adds the given string to the spinner
// as the final message to be written
func WithFinalMSG(finalMsg string) Option {
	return func(s *Spinner) {
		s.FinalMSG = finalMsg
	}
}

// WithHiddenCursor hides the cursor
// if hideCursor = true given
func WithHiddenCursor(hideCursor bool) Option {
	return func(s *Spinner) {
		s.HideCursor = hideCursor
	}
}

// Active will return whether or not the spinner is currently active
func (s *Spinner) Active() bool {
	return s.active
}

// Start will start the indicator
func (s *Spinner) Start() {
	s.lock.Lock()
	if s.active {
		s.lock.Unlock()
		return
	}
	s.active = true
	s.lock.Unlock()

	go func() {
		for {
			for i := 0; i < len(s.chars); i++ {
				select {
				case <-s.stopChan:
					return
				default:
					s.lock.Lock()
					s.erase()
					if !s.active {
						return
					}
					outPlain := fmt.Sprintf("%s%s%s ", s.Prefix, s.chars[i], s.Suffix)
					fmt.Fprint(s.Writer, outPlain)
					s.lastOutput = outPlain
					delay := s.Delay
					s.lock.Unlock()

					time.Sleep(delay)
				}
			}
		}
	}()
}

// Stop stops the indicator
func (s *Spinner) Stop() {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.active {
		s.active = false
		s.erase()
		if s.FinalMSG != "" {
			fmt.Fprintf(s.Writer, "%s", s.FinalMSG)
		}
		s.stopChan <- struct{}{}
	}
}

// Restart will stop and start the indicator
func (s *Spinner) Restart() {
	s.Stop()
	s.Start()
}

// UpdateSpeed will set the indicator delay to the given value
func (s *Spinner) UpdateSpeed(d time.Duration) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Delay = d
}

// UpdateCharSet will change the current character set to the given one
func (s *Spinner) UpdateCharSet(cs []string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.chars = cs
}

// erase deletes written characters
//
// Caller must already hold s.lock.
func (s *Spinner) erase() {
	n := utf8.RuneCountInString(s.lastOutput)
	clearString := "\r"
	for i := 0; i < n; i++ {
		clearString += " "
	}
	clearString += "\r"
	fmt.Fprint(s.Writer, clearString)
	s.lastOutput = ""
}

// Lock allows for manual control to lock the spinner
func (s *Spinner) Lock() {
	s.lock.Lock()
}

// Unlock allows for manual control to unlock the spinner
func (s *Spinner) Unlock() {
	s.lock.Unlock()
}

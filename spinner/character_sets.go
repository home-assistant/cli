// Package spinner Originally created by Brian Downs
// https://github.com/briandowns/spinner
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
package spinner

const (
	clockOneOClock = '\U0001F550'
	clockOneThirty = '\U0001F55C'
)

// CharSets contains the available character sets
var CharSets = map[int][]string{
	0: {"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"},
	1: {"⢹", "⢺", "⢼", "⣸", "⣇", "⡧", "⡗", "⡏"},
}

func init() {
	for i := rune(0); i < 12; i++ {
		CharSets[37] = append(CharSets[37], string([]rune{clockOneOClock + i}))
		CharSets[38] = append(CharSets[38], string([]rune{clockOneOClock + i}), string([]rune{clockOneThirty + i}))
	}
}

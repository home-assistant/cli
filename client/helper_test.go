package client

import (
	"fmt"
	"testing"
)

var urlTests = []struct {
	base    string
	section string
	command string
	out     string
}{
	{"hassio", "section", "command", "http://hassio/section/command"},
	{"hassio:80", "section", "command", "http://hassio:80/section/command"},
	{"hassio.example.org:8080", "section", "command", "http://hassio.example.org:8080/section/command"},
	{"https://hassio", "section", "command", "https://hassio/section/command"},
	{"https://hassio:8080", "section", "command", "https://hassio:8080/section/command"},
	{"https://hassio", "section", "command", "https://hassio/section/command"},
	{"https://hassio:8080", "section", "command", "https://hassio:8080/section/command"},
	{"hassio", "section", "", "http://hassio/section"},
	{"hassio", "section/../othersection", "", "http://hassio/othersection"},
	{"hassio/api/", "section", "command", "http://hassio/api/section/command"},
}

func TestURLHelper(t *testing.T) {
	for _, tt := range urlTests {
		t.Run(fmt.Sprintf("[%s][%s][%s]", tt.base, tt.section, tt.command), func(t *testing.T) {
			s, _ := URLHelper(tt.base, tt.section, tt.command)
			if s != tt.out {
				t.Errorf("got %q, want %q", s, tt.out)
			}
		})
	}
}

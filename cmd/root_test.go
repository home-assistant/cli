package cmd

import "testing"

func TestDefaultEndpoint(t *testing.T) {
	tests := []struct {
		name   string
		socket string
		want   string
	}{
		{"socket unset", "", "supervisor"},
		{"socket set", "/run/supervisor/api/supervisor.sock", "unix:///run/supervisor/api/supervisor.sock"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("SUPERVISOR_SOCKET", tt.socket)
			if got := defaultEndpoint(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

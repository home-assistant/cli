package client

import "testing"

func TestToYAML(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  string
	}{
		{
			"scalars",
			map[string]any{"version": "2026.11.0", "count": 12345678901, "ratio": 1.5, "healthy": true, "channel": nil},
			"channel: null\ncount: 12345678901\nhealthy: true\nratio: 1.5\nversion: 2026.11.0\n",
		},
		{
			"ambiguous strings are quoted",
			map[string]any{"bool": "yes", "number": "123", "version": "1.10", "empty": ""},
			"bool: \"yes\"\nempty: \"\"\nnumber: \"123\"\nversion: \"1.10\"\n",
		},
		{
			"lists are indented under their key",
			map[string]any{"addons": []any{
				map[string]any{"slug": "core_ssh", "ports": map[string]any{"22/tcp": 22}},
				map[string]any{"slug": "core_samba"},
			}},
			"addons:\n  - ports:\n      22/tcp: 22\n    slug: core_ssh\n  - slug: core_samba\n",
		},
		{
			"empty collections",
			map[string]any{"list": []any{}, "map": map[string]any{}},
			"list: []\nmap: {}\n",
		},
		{
			"multiline string",
			map[string]any{"text": "first\nsecond"},
			"text: |-\n  first\n  second\n",
		},
		{
			"struct uses JSON field names",
			&Response{Result: "unknown", Message: "what"},
			"message: what\nresult: unknown\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toYAML(tt.value)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != tt.want {
				t.Errorf("got:\n%s\nwant:\n%s", got, tt.want)
			}
		})
	}
}

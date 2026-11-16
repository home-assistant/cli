package cmd

import (
	"sort"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func TestMissingPositionalArgCompletions(t *testing.T) {
	var cmds []*cobra.Command
	var walk func(c *cobra.Command)
	walk = func(c *cobra.Command) {
		if !c.Hidden {
			cmds = append(cmds, c)
		}
		for _, sub := range c.Commands() {
			walk(sub)
		}
	}
	walk(rootCmd)

	names := make([]string, 0, len(cmds))
	for _, c := range cmds {
		if c.HasAvailableSubCommands() {
			continue
		}
		if c.ValidArgsFunction == nil {
			names = append(names, c.CommandPath())
		}
	}
	sort.Strings(names)
	for _, n := range names {
		t.Errorf("missing positional arg completion: %s", n)
	}
}

func TestMissingFlagCompletions(t *testing.T) {
	var cmds []*cobra.Command
	var walk func(c *cobra.Command)
	walk = func(c *cobra.Command) {
		if !c.Hidden {
			cmds = append(cmds, c)
		}
		for _, sub := range c.Commands() {
			walk(sub)
		}
	}
	walk(rootCmd)

	missing := map[string]bool{}
	for _, c := range cmds {
		flags := c.LocalFlags()
		flags.VisitAll(func(f *pflag.Flag) {
			if f.Hidden {
				return
			}
			if _, ok := c.GetFlagCompletionFunc(f.Name); !ok {
				if _, present := f.Annotations[cobra.BashCompFilenameExt]; !present {
					missing[c.CommandPath()+" --"+f.Name] = true
				}
			}
		})
		c.InheritedFlags().VisitAll(func(f *pflag.Flag) {
			if f.Hidden {
				return
			}
			if _, ok := c.GetFlagCompletionFunc(f.Name); !ok {
				if _, present := f.Annotations[cobra.BashCompFilenameExt]; !present {
					missing[c.CommandPath()+" --"+f.Name] = true
				}
			}
		})
	}
	names := make([]string, 0, len(missing))
	for k := range missing {
		names = append(names, k)
	}
	sort.Strings(names)
	for _, n := range names {
		t.Errorf("missing flag completion: %s", n)
	}
}

package main

import (
	"reflect"
	"testing"
)

func TestCmdRoot(t *testing.T) {
	validCommands := []string{
		"atc",
		"mayday",
		"unlatch",
		"schematics",
		"verify",
		"version",
		"takeoff",
		"descent",
		"blackbox",
		"turbulence",
		"stow",
		"sign",
	}
	for _, cmd := range validCommands {
		if _, ok := CmdRoot.SubCommands[cmd]; !ok {
			t.Fatalf("expected %s to be a subcommand of the root command", cmd)
		}
	}
}

func TestCmdSeek(t *testing.T) {
	for _, s := range []struct {
		Want CmdRunner
		Args []string
	}{
		{
			Want: CmdRoot.Runner,
			Args: []string{"yoke"},
		},
	} {
		runner, _ := Seek(s.Args)
		if reflect.ValueOf(runner) != reflect.ValueOf(s.Want) {
			t.Fatalf("got function mismatch for %v", s.Args)
		}
	}
}

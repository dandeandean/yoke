package main

import (
	"context"
	"flag"
)

type YokeCommandRunner interface {
	Run(ctx context.Context, settings GlobalSettings, subCommands string) error
}

// The YokeCommand struct represents a cli commmand
// It should have a name, alias, and subcommands
type YokeCommand struct {
	Name           string
	Aliases        []string
	FlagSet        *flag.FlagSet
	SubCommands    map[string]*YokeCommand
	CompletionFunc func([]string)
	Parent         *YokeCommand
	Runner         CmdRunner
}

// We might actually want to implement this as a map to make it blazingly fast
func (y *YokeCommand) AddCommand(sub *YokeCommand) {
	sub.Parent = y
	y.SubCommands[sub.Name] = sub
	for _, alias := range sub.Aliases {
		_, alreadyThere := y.SubCommands[alias]
		if !alreadyThere {
			y.SubCommands[alias] = sub
		}
	}
}

func (y YokeCommand) AllNames() []string {
	return append(y.Aliases, y.Name)
}

type CmdRunner func(ctx context.Context, settings GlobalSettings, args []string) error
type CmdBuilder func(ctx context.Context) (*flag.FlagSet, CmdRunner)

func NewCommand(name string, aliases []string, builder CmdBuilder) *YokeCommand {
	flagset, runner := builder(context.Background())
	return &YokeCommand{
		Name:        name,
		Aliases:     aliases,
		FlagSet:     flagset,
		Runner:      runner,
		SubCommands: make(map[string]*YokeCommand),
	}
}

func Seek(args []string) CmdRunner {
	cmdPtr := CmdRoot
	for _, arg := range args {
		nextCmd, ok := cmdPtr.SubCommands[arg]
		if ok {
			cmdPtr = nextCmd
		}
	}
	return cmdPtr.Runner
}

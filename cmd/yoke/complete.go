package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
)

func cleanArg(argIn string) string {
	return argIn
}

func printFlagCompletion(args []string, cmd *YokeCommand) {
	for _, flag := range getFlagCompletion(args, cmd) {
		fmt.Println(flag)
	}
}

// get the flags associated with a yokeCommand
// it takes all of args slice and the YokeCommand
func getFlagCompletion(args []string, cmd *YokeCommand) []string {
	flagSetAll := make(map[string]bool)
	out := make([]string, 0)
	partial := strings.TrimLeft(args[len(args)-1], "-")
	if cmd.FlagSet == nil {
		return out
	}
	appendWithPrefix := func(f *flag.Flag, p string) {
		if p == "" || strings.HasPrefix(f.Name, p) {
			if !slices.Contains(args, f.Name) {
				flagSetAll["-"+f.Name] = true
			}
		}
	}
	// Iterate through all of the places we get flags from
	cur := cmd
	for cur != nil && cur.FlagSet != nil {
		cur.FlagSet.VisitAll(func(f *flag.Flag) {
			appendWithPrefix(f, partial)
		})
		cur = cur.Parent
	}

	for k := range flagSetAll {
		out = append(out, k)
	}
	return out
}

// given the args passed, yield all of the valid next top level cocmmands
func getCommandCompletions(args []string) []*YokeCommand {
	out := make([]*YokeCommand, 0)
	if len(args) == 0 {
		return out
	}
	partial := args[len(args)-1]
	if partial == "complete" || partial == "yoke" {
		partial = ""
	}

	cmd, rest := Seek(args)
	// we've hit the end, set the partial to ""
	if len(rest) == 0 {
		partial = ""
	}
	fmt.Println("DEBUG: got to ", cmd.Name, "partial=", partial, "<-")
	for k, v := range cmd.SubCommands {
		if strings.HasPrefix(k, partial) || partial == "" {
			fmt.Println("DEBUG: appending ", k)
			out = append(out, v)
		}
	}
	return out
}

func printCommandCompletions(args []string) {
	for _, cmd := range getCommandCompletions(args) {
		fmt.Println(cmd.Name)
	}
}

func Complete() {
	if len(os.Args) < 2 {
		return
	}
	argSet := make(map[string]bool)
	for _, arg := range os.Args {
		argSet[arg] = true
	}
	argsAfterComp := os.Args[2:]
	if len(argsAfterComp) > 1 && argsAfterComp[0] == "yoke" {
		argsAfterComp = argsAfterComp[1:]
	}
	fmt.Println("DEBUG: seeking for ", argsAfterComp)
	cmd, rest := Seek(argsAfterComp)
	fmt.Println("DEBUG: result for seek: ", cmd.Name, rest)
	partial := ""
	if len(rest) > 0 {
		partial = rest[len(rest)-1]
	}
	fmt.Println("DEBUG: ", cmd.Name, rest, argsAfterComp)
	if strings.HasPrefix(partial, "-") {
		printFlagCompletion(rest, cmd)
	}
	//if it's not a full command, get top-level completions
	//printCommandCompletions(os.Args)
}

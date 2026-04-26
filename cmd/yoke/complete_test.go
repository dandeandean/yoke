package main

import (
	"os"
	"testing"
)

// returns true if a is a subset of b
func isSubset(a, b []string) bool {
	set := make(map[string]bool)
	for _, s := range b {
		set[s] = true
	}
	for _, s := range a {
		_, ok := set[s]
		if !ok {
			return false
		}
	}
	return true
}

func TestCompFlags(t *testing.T) {
	comps := getFlagCompletion(
		[]string{"yoke", "descent", "-"},
		CmdRoot.SubCommands["descent"],
	)
	cmpFlags := isSubset([]string{
		"-debug",
		"-kube-context",
		"-namespace",
		"-poll",
		"-remove-all",
		"-remove-crds",
		"-wait",
		"-kubeconfig",
		"-lock",
		"-remove-namespaces",
	}, comps)
	if cmpFlags {
		t.Fatal("TestDescentFlagCompletions did not yield expected flags, got: ", comps, "ARGS: ", os.Args)
	}
}

package main

import (
	"cmp"
	"context"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/davidmdm/x/xcontext"

	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/yokecd/yoke/internal"
	"github.com/yokecd/yoke/internal/home"
	"github.com/yokecd/yoke/pkg/yoke"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		if internal.IsWarning(err) {
			return
		}
		os.Exit(1)
	}
}

//go:embed cmd_help.txt
var rootHelp string

var CmdRoot = NewCommand("yoke", []string{}, func(ctx context.Context) (*flag.FlagSet, CmdRunner) {
	flagset := flag.NewFlagSet("yoke", flag.ExitOnError)
	flagset.Usage = func() {
		rootHelp = strings.TrimSpace(internal.Colorize(rootHelp))
		fmt.Fprintln(flag.CommandLine.Output(), rootHelp)
		flagset.PrintDefaults()
		fmt.Fprintln(os.Stderr)
	}
	runner := func(ctx context.Context, settings GlobalSettings, args []string) error {
		RegisterGlobalFlags(flagset, &settings)
		return nil
	}
	return flagset, runner
})

var settings = GlobalSettings{
	Debug: new(bool),
	Kube:  genericclioptions.NewConfigFlags(false),
}

func init() {
	CmdRoot.AddCommand(CmdATC)
	CmdRoot.AddCommand(CmdBlackbox)
	CmdRoot.AddCommand(CmdDescent)
	CmdRoot.AddCommand(CmdMayday)
	CmdRoot.AddCommand(CmdSchematics)
	CmdRoot.AddCommand(CmdSign)
	CmdRoot.AddCommand(CmdStow)
	CmdRoot.AddCommand(CmdTakeoff)
	CmdRoot.AddCommand(CmdTurbulence)
	CmdRoot.AddCommand(CmdVersion)
}

func run() error {

	if len(os.Args) > 1 && os.Args[1] == "complete" {
		Complete()
		return nil
	}

	CmdRoot.FlagSet.Parse(os.Args)

	ctx, cancel := xcontext.WithSignalCancelation(context.Background(), syscall.SIGINT)
	defer cancel()

	ctx = internal.WithDebugFlag(ctx, settings.Debug)

	if len(CmdRoot.FlagSet.Args()) == 0 {
		CmdRoot.FlagSet.Usage()
		return fmt.Errorf("no command provided")
	}

	subcmdArgs := CmdRoot.FlagSet.Args()[2:]

	switch cmd := CmdRoot.FlagSet.Arg(1); cmd {
	case "atc":
		return CmdATC.Runner(ctx, settings, subcmdArgs)
	case "blackbox", "inspect":
		return CmdBlackbox.Runner(ctx, settings, subcmdArgs)
	case "descent", "down", "restore":
		return CmdDescent.Runner(ctx, settings, subcmdArgs)
	case "mayday", "delete":
		return CmdMayday.Runner(ctx, settings, subcmdArgs)
	case "schematics", "meta":
		{
			runner := Seek(CmdRoot.FlagSet.Args())
			// FIXME: doesn't work if we just pass in args
			// it seems to be because subcmdArgs  passes the flag first,
			// which works... maybe we could handle that in seek somehow
			return runner(ctx, settings, subcmdArgs)
			// It might also make more sense for wasmfile to be positional
		}
	case "sign":
		return CmdSign.Runner(ctx, settings, subcmdArgs)
	case "stow", "push":
		return CmdStow.Runner(ctx, settings, subcmdArgs)
	case "takeoff", "up", "apply":
		CmdTakeoff.Runner(ctx, settings, subcmdArgs)
	case "turbulence", "drift", "diff":
		CmdTurbulence.Runner(ctx, settings, subcmdArgs)
	case "unlatch", "unlock":
		{
			params, err := GetUnlatchParams(settings, subcmdArgs)
			if err != nil {
				return err
			}
			return Unlatch(ctx, *params)
		}

	case "verify":
		{
			params, err := GetVerifyParams(subcmdArgs)
			if err != nil {
				return err
			}
			return yoke.Verify(*params)
		}
	case "version":
		return CmdVersion.Runner(ctx, settings, subcmdArgs)
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}
	return nil
}

type GlobalSettings struct {
	Kube  *genericclioptions.ConfigFlags
	Debug *bool
}

func RegisterGlobalFlags(flagset *flag.FlagSet, settings *GlobalSettings) {
	flagset.StringVar(settings.Kube.KubeConfig, "kubeconfig", cmp.Or(*settings.Kube.KubeConfig, os.Getenv("KUBECONFIG"), home.Kubeconfig), "path to kube config")
	flagset.StringVar(settings.Kube.Context, "kube-context", *settings.Kube.Context, "kubernetes context to use")
	flagset.BoolVar(settings.Debug, "debug", *settings.Debug, "debug output mode")
}

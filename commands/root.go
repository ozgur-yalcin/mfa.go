package commands

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/ozgur-yalcin/mfa/internal/initialize"
)

type rootCommand struct {
	name     string
	use      string
	commands []Commander
	fs       *flag.FlagSet
}

func (r *rootCommand) Name() string {
	return r.name
}

func (r *rootCommand) Use() string {
	return r.use
}

func (r *rootCommand) Init(cd *Ancestor) error {
	r.fs = flag.NewFlagSet(r.name, flag.ExitOnError)
	return nil
}

func (r *rootCommand) Run(ctx context.Context, cd *Ancestor, args []string) error {
	slog.Debug(fmt.Sprintf("mfa version %q finishing with parameters %q", initialize.Version, os.Args))
	return nil
}

func (r *rootCommand) Commands() []Commander {
	return r.commands
}

func (r *Exec) Execute(ctx context.Context, args []string) (*Ancestor, error) {
	if err := r.c.init(); err != nil {
		return nil, err
	}
	cd := r.c
	if len(args) > 0 {
		for _, subcmd := range r.c.ancestors {
			if subcmd.Commander.Name() == args[0] {
				cd = subcmd
				break
			}
		}
	}
	if err := cd.Command.Parse(args); err != nil {
		return cd, err
	}
	if err := cd.Commander.Run(ctx, cd, cd.Command.Args()[1:]); err != nil {
		return cd, err
	}
	return cd, nil
}

func Execute(args []string) error {
	x, err := newExec()
	if err != nil {
		return err
	}
	if _, err := x.Execute(context.Background(), args); err != nil {
		return err
	}
	return err
}

func New(rootCmd Commander) (*Exec, error) {
	rootCd := &Ancestor{
		Commander: rootCmd,
	}
	rootCd.Root = rootCd
	var addCommands func(cd *Ancestor, cmd Commander)
	addCommands = func(cd *Ancestor, cmd Commander) {
		cd2 := &Ancestor{
			Root:      rootCd,
			Parent:    cd,
			Commander: cmd,
		}
		cd.ancestors = append(cd.ancestors, cd2)
		for _, c := range cmd.Commands() {
			addCommands(cd2, c)
		}
	}
	for _, cmd := range rootCmd.Commands() {
		addCommands(rootCd, cmd)
	}
	if err := rootCd.compile(); err != nil {
		return nil, err
	}
	return &Exec{c: rootCd}, nil
}

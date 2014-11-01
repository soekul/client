package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/keybase/go-libkb"
	"github.com/keybase/go-jsonw"
)

type CmdTrack struct {
	user      string
	assertion string
	track     bool
}

func (v *CmdTrack) ParseArgv(ctx *cli.Context) error {
	nargs := len(ctx.Args())
	var err error
	v.track = ctx.Bool("track-statement")
	if nargs == 1 {
		v.user = ctx.Args()[0]
	} else {
		err = fmt.Errorf("track takes one arg -- the user to track")
	}
	return err
}

func (v *CmdTrack) Run() error {

	u2, err := libkb.LoadUser(libkb.LoadUserArg{
		Name:             v.user,
		RequirePublicKey: true,
		Self:             false,
		LoadSecrets:      false,
		ForceReload:      false,
		SkipVerify:       false,
	})

	if err != nil {
		return err
	}

	me, err := libkb.LoadMe()

	if err != nil {
		return err
	}

	if me.Equal(*u2) {
		return fmt.Errorf("Cannot track yourself")
	}

	err = u2.Identify()
	if err != nil {
		G.Log.Warning("Some proofs failed")
	}

	var jw *jsonw.Wrapper 
	jw, err = me.TrackingProofFor(u2)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", jw.MarshalPretty())
	fmt.Printf("%v\n", u2.IdTable.MakeTrackSet())


	return nil
}

func NewCmdTrack(cl *CommandLine) cli.Command {
	return cli.Command{
		Name:        "track",
		Usage:       "keybase track <username>",
		Description: "verify a user's authenticity and optionally track them",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "assert, a",
				Usage: "a boolean expression on this identity",
			},
		},
		Action: func(c *cli.Context) {
			cl.ChooseCommand(&CmdTrack{}, "track", c)
		},
	}
}

func (v *CmdTrack) UseConfig() bool   { return true }
func (v *CmdTrack) UseKeyring() bool  { return true }
func (v *CmdTrack) UseAPI() bool      { return true }
func (v *CmdTrack) UseTerminal() bool { return true }

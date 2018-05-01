// Copyright 2015 Keybase, Inc. All rights reserved. Use of
// this source code is governed by the included BSD license.

package engine

import (
	"fmt"
	"time"

	"github.com/keybase/client/go/libkb"
	keybase1 "github.com/keybase/client/go/protocol/keybase1"
	"golang.org/x/net/context"
)

type Context struct {
	GPGUI       libkb.GPGUI
	LogUI       libkb.LogUI
	LoginUI     libkb.LoginUI
	SecretUI    libkb.SecretUI
	IdentifyUI  libkb.IdentifyUI
	PgpUI       libkb.PgpUI
	ProveUI     libkb.ProveUI
	ProvisionUI libkb.ProvisionUI

	LoginContext libkb.LoginContext
	NetContext   context.Context
	SaltpackUI   libkb.SaltpackUI

	// Usually set to `NONE`, meaning none specified.
	// But if we know it, specify the end client type here
	// since some things like GPG shell-out work differently
	// depending.
	ClientType keybase1.ClientType

	// Special-case flag for identifyUI -- if it's been delegated
	// to the electron UI, then it's rate-limitable
	IdentifyUIIsDelegated bool

	SessionID int
}

func engineContextFromMetaContext(m libkb.MetaContext) *Context {
	uis := m.UIs()
	return &Context{
		GPGUI:                 uis.GPGUI,
		LogUI:                 uis.LogUI,
		LoginUI:               uis.LoginUI,
		SecretUI:              uis.SecretUI,
		IdentifyUI:            uis.IdentifyUI,
		PgpUI:                 uis.PgpUI,
		ProveUI:               uis.ProveUI,
		ProvisionUI:           uis.ProvisionUI,
		LoginContext:          m.LoginContext(),
		NetContext:            m.Ctx(),
		SaltpackUI:            uis.SaltpackUI,
		ClientType:            uis.ClientType,
		IdentifyUIIsDelegated: uis.IdentifyUIIsDelegated,
		SessionID:             uis.SessionID,
	}
}

func metaContextFromEngineContext(g *libkb.GlobalContext, ctx *Context) libkb.MetaContext {
	uis := libkb.UIs{
		GPGUI:                 ctx.GPGUI,
		LogUI:                 ctx.LogUI,
		LoginUI:               ctx.LoginUI,
		SecretUI:              ctx.SecretUI,
		IdentifyUI:            ctx.IdentifyUI,
		PgpUI:                 ctx.PgpUI,
		ProveUI:               ctx.ProveUI,
		ProvisionUI:           ctx.ProvisionUI,
		SaltpackUI:            ctx.SaltpackUI,
		ClientType:            ctx.ClientType,
		IdentifyUIIsDelegated: ctx.IdentifyUIIsDelegated,
		SessionID:             ctx.SessionID,
	}
	return libkb.NewMetaContext(ctx.GetNetContext(), g).WithLoginContext(ctx.LoginContext).WithUIs(uis)
}

func (c *Context) HasUI(kind libkb.UIKind) bool {
	switch kind {
	case libkb.GPGUIKind:
		return c.GPGUI != nil
	case libkb.LogUIKind:
		return c.LogUI != nil
	case libkb.LoginUIKind:
		return c.LoginUI != nil
	case libkb.SecretUIKind:
		return c.SecretUI != nil
	case libkb.IdentifyUIKind:
		return c.IdentifyUI != nil
	case libkb.PgpUIKind:
		return c.PgpUI != nil
	case libkb.ProveUIKind:
		return c.ProveUI != nil
	case libkb.ProvisionUIKind:
		return c.ProvisionUI != nil
	case libkb.SaltpackUIKind:
		return c.SaltpackUI != nil
	}
	panic(fmt.Sprintf("unhandled kind:  %d", kind))
}

func (c *Context) GetNetContext() context.Context {
	if c.NetContext == nil {
		return context.Background()
	}
	return c.NetContext
}

func (c *Context) SetNetContext(netCtx context.Context) {
	c.NetContext = netCtx
}

// A copy of the Context with the NetContext swapped out
func (c *Context) WithNetContext(netCtx context.Context) *Context {
	c2 := *c
	c2.NetContext = netCtx
	return &c2
}

func (c *Context) WithCancel() (*Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(c.GetNetContext())
	return c.WithNetContext(ctx), cancel
}

func (c *Context) WithTimeout(timeout time.Duration) (*Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(c.GetNetContext(), timeout)
	return c.WithNetContext(ctx), cancel
}

func SecretKeyPromptArg(ui libkb.SecretUI, ska libkb.SecretKeyArg, reason string) libkb.SecretKeyPromptArg {
	return libkb.SecretKeyPromptArg{
		SecretUI: ui,
		Ska:      ska,
		Reason:   reason,
	}
}

func (c *Context) SecretKeyPromptArg(ska libkb.SecretKeyArg, reason string) libkb.SecretKeyPromptArg {
	return SecretKeyPromptArg(c.SecretUI, ska, reason)
}

func (c *Context) CloneGlobalContextWithLogTags(g *libkb.GlobalContext, k string) *libkb.GlobalContext {
	netCtx := libkb.WithLogTag(c.GetNetContext(), k)
	c.NetContext = netCtx
	return g.CloneWithNetContextAndNewLogger(netCtx)
}

func NewMetaContext(e Engine, c *Context) libkb.MetaContext {
	return libkb.NewMetaContext(c.GetNetContext(), e.G()).WithLoginContext(c.LoginContext)
}
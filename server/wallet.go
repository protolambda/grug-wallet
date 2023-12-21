package server

import (
	"context"
	"errors"
	"github.com/protolambda/grug-wallet/server/jsonrpc"
	"math/big"
	"sync"
)

type Provider interface {
	ChangeAccount(id ResourceID[*Account]) // metamask_accountsChanged
	ChangeChain(id ResourceID[*Chain])     // metamask_chainChanged
	ChangeState(locked bool)               // metamask_unlockStateChanged
}

type ProviderImpl struct {
	user    jsonrpc.Channel
	backend jsonrpc.Channel

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	chainID          *big.Int
	accounts         []Address
	unlockedAccounts bool
}

func NewProviderImpl(user, backend jsonrpc.Channel) *ProviderImpl {
	ctx, cancel := context.WithCancel(context.Background())
	prov := &ProviderImpl{
		user:    user,
		backend: backend,
		ctx:     ctx,
		cancel:  cancel,
	}
	prov.wg.Add(1)
	go prov.run()
	return prov
}

func (p *ProviderImpl) run() {
	defer p.wg.Done()

	for {
		select {
		case msg := <-p.user.Inward:
			p.handleRPC(msg)
		case <-p.ctx.Done():
			return
		}
	}
}

func (p *ProviderImpl) handleRPC(m *jsonrpc.Message) {
	if m.Request == nil {
		// responses are unexpected

	}
	switch m.Request.Method {
	case "wallet_addEthereumChain":
		// TODO UI
	case "wallet_switchEthereumChain":
		// TODO UI
	case "wallet_requestPermissions":
		// TODO UI
	case "wallet_revokePermissions":
		// TODO UI
	case "wallet_getPermissions":
		if p.unlockedAccounts {

		}
	case "wallet_registerOnboarding":
		p.send(m.RespondErr(jsonrpc.AnnotatedErrorObj(jsonrpc.MethodNotSupported, errors.New("no onboarding support"))))
	case "wallet_watchAsset":
		// TODO UI - add token to tracker (ERC20, ERC721, ERC1155)
		p.send(m.Respond(true))
	case "wallet_scanQRCode":
		p.send(m.RespondErr(jsonrpc.AnnotatedErrorObj(jsonrpc.MethodNotSupported, errors.New("no camera support"))))
	case "web3_clientVersion":
		p.send(m.Respond("Grug/v0.0.1")) // TODO allow override, to mask as metamask
	case "eth_requestAccounts":
		// TODO UI
	case "eth_accounts":
		if p.unlockedAccounts {
			p.send(m.Respond(p.accounts))
		} else {
			p.send(m.RespondErr(jsonrpc.AnnotatedErrorObj(jsonrpc.Unauthorized, errors.New("account not unlocked"))))
		}
	case "eth_chainId":
		p.send(m.Respond((*HexBig)(p.chainID)))
	default:
		// TODO: timeout / UI-cancel, in case backend becomes unresponsive
		select {
		case p.backend.Outward <- m:
		case <-p.ctx.Done():
		}
	}
}

func (p *ProviderImpl) send(m *jsonrpc.Message) {
	select {
	case <-p.ctx.Done():
		// TODO
	case p.user.Outward <- m:
		return
	}
}

func (p *ProviderImpl) Close() error {
	p.cancel()
	p.wg.Wait()
	return nil
}

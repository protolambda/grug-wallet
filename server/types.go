package server

import (
	"math/big"
)

type ResourceKind string

const (
	KindCurrency       ResourceKind = "currency"
	KindExplorer       ResourceKind = "explorer"
	KindChain          ResourceKind = "chain"
	KindRPCVendor      ResourceKind = "rpc-vendor"
	KindRPC            ResourceKind = "rpc"
	KindRPCComposition ResourceKind = "rpc-composition"
	KindAccount        ResourceKind = "account"
	KindWallet         ResourceKind = "wallet"
	KindSigner         ResourceKind = "signer"
	KindHardware       ResourceKind = "hardware"
	KindContractABI    ResourceKind = "contract-abi"
	KindTokenERC20     ResourceKind = "token-erc20"
	KindTokenERC721    ResourceKind = "token-erc721"
	KindTokenERC1155   ResourceKind = "token-erc1155"

	KindHolding ResourceKind = "holding"
	KindScript  ResourceKind = "script"
)

type Resource interface {
	ResourceKind() ResourceKind
}

type ResourceID[R Resource] string

func (r ResourceID[R]) String() string {
	return string(r)
}

type Currency struct {
	Name     string
	Icon     string
	Symbol   string
	Decimals uint
}

func (c *Currency) ResourceKind() ResourceKind {
	return KindCurrency
}

type Explorer struct {
	Name     string
	Icon     string
	Endpoint string
	// standards
}

func (c *Explorer) ResourceKind() ResourceKind {
	return KindExplorer
}

type Chain struct {
	Name      string
	ChainID   *big.Int
	NetworkID *big.Int

	Currency  ResourceID[*Currency]
	Explorers []Explorer
}

func (c *Chain) ResourceKind() ResourceKind {
	return KindChain
}

type RPCVendor struct {
	Name string
	Icon string
}

func (c *RPCVendor) ResourceKind() ResourceKind {
	return KindRPCVendor
}

type MethodRateLimit struct {
	Matcher string
	Limit   int // TODO
}

type RPC struct {
	Name        string
	Description string
	Vendor      ResourceID[*RPCVendor]

	Chain      ResourceID[*Chain]
	Methods    []string
	EngineAuth string

	RateLimit []MethodRateLimit

	AdminPanel string
}

func (c *RPC) ResourceKind() ResourceKind {
	return KindRPC
}

type MethodForward struct {
	Matcher string
	Forward ResourceID[*RPC]
}

type RPCComposition struct {
	Name        string
	Description string
	Methods     []MethodForward
}

func (c *RPCComposition) ResourceKind() ResourceKind {
	return KindRPCComposition
}

type Account struct {
	Name        string
	Description string
	Address     string
	Color       string
}

func (c *Account) ResourceKind() ResourceKind {
	return KindAccount
}

type Wallet struct {
	Name        string
	Description string
	Addresses   []string
}

func (c *Wallet) ResourceKind() ResourceKind {
	return KindWallet
}

type Signer struct {
	Name        string
	Description string
	Matcher     string
	Hardware    ResourceID[*Hardware]
}

func (c *Signer) ResourceKind() ResourceKind {
	return KindSigner
}

type Hardware struct {
	Name        string
	Description string
	Type        string
}

func (c *Hardware) ResourceKind() ResourceKind {
	return KindHardware
}

type ContractABI struct {
	Name    string
	Address string
	Type    string
}

func (c *ContractABI) ResourceKind() ResourceKind {
	return KindContractABI
}

type TokenERC20 struct {
	Name     string
	Icon     string
	Symbol   string
	Decimals uint
	Address  string

	Chain ResourceID[*Chain]
}

func (c *TokenERC20) ResourceKind() ResourceKind {
	return KindTokenERC20
}

type TokenERC721 struct {
	Name    string
	Address string
	TokenID string

	Chain ResourceID[*Chain]
}

func (c *TokenERC721) ResourceKind() ResourceKind {
	return KindTokenERC721
}

type TokenERC1155 struct {
	Name    string
	Address string
	TokenID string

	Chain ResourceID[*Chain]
}

func (c *TokenERC1155) ResourceKind() ResourceKind {
	return KindTokenERC1155
}

type Holding struct {
	Name        string
	Description string
	Account     ResourceID[*Account]
	Token       ResourceID[*TokenERC20]
	Currency    ResourceID[*Currency]
}

func (c *Holding) ResourceKind() ResourceKind {
	return KindHolding
}

type KV struct {
	K string
	V string
}

type RPCDependency struct {
	K     string
	Chain ResourceID[*Chain]
}

type Script struct {
	Name       string
	Inputs     []KV // env
	Invocation string
	RPCs       []RPCDependency
}

func (c *Script) ResourceKind() ResourceKind {
	return KindScript
}

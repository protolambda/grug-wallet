package server

import (
	"fmt"
	"github.com/google/uuid"
	"math/big"
)

type ResourceKind string

const (
	KindCurrency     ResourceKind = "currency"
	KindExplorer     ResourceKind = "explorer"
	KindChain        ResourceKind = "chain"
	KindRPCVendor    ResourceKind = "rpc-vendor"
	KindRPC          ResourceKind = "rpc"
	KindCompositeRPC ResourceKind = "composite-rpc"
	KindProviderRPC  ResourceKind = "provider-rpc"
	KindAccount      ResourceKind = "account"
	KindWallet       ResourceKind = "wallet"
	KindSigner       ResourceKind = "signer"
	KindHardware     ResourceKind = "hardware"
	KindContractABI  ResourceKind = "contract-abi"
	KindTokenERC20   ResourceKind = "token-erc20"
	KindTokenERC721  ResourceKind = "token-erc721"
	KindTokenERC1155 ResourceKind = "token-erc1155"

	KindHolding ResourceKind = "holding"
	KindScript  ResourceKind = "script"
)

type Resource interface {
	ResourceKind() ResourceKind
	Init(basePath string, mngr *ResourceManager) error
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

func (c *Currency) Init(basePath string, mngr *ResourceManager) error {
	return nil
}

type Explorer struct {
	Name     string
	Icon     string
	Endpoint string
	Chain    ResourceID[*Chain]
}

func (c *Explorer) ResourceKind() ResourceKind {
	return KindExplorer
}

func (c *Explorer) Init(basePath string, mngr *ResourceManager) error {
	if err := c.Chain.Load(basePath, mngr); err != nil {
		return fmt.Errorf("failed to load chain of explorer: %w", err)
	}
	return nil
}

type Chain struct {
	Name      string
	ChainID   *big.Int
	NetworkID *big.Int

	Currency ResourceID[*Currency]
}

func (c *Chain) ResourceKind() ResourceKind {
	return KindChain
}

func (c *Chain) Init(basePath string, mngr *ResourceManager) error {
	if err := c.Currency.Load(basePath, mngr); err != nil {
		return fmt.Errorf("failed to load currency of chain: %w", err)
	}
	return nil
}

type RPCVendor struct {
	Name string
	Icon string
}

func (c *RPCVendor) ResourceKind() ResourceKind {
	return KindRPCVendor
}

func (c *RPCVendor) Init(basePath string, mngr *ResourceManager) error {
	return nil
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
	Endpoint   string

	RateLimit []MethodRateLimit

	AdminPanel string
}

func (c *RPC) ResourceKind() ResourceKind {
	return KindRPC
}

func (c *RPC) Init(basePath string, mngr *ResourceManager) error {
	if err := c.Vendor.Load(basePath, mngr); err != nil {
		return fmt.Errorf("failed to load vendor of RPC: %w", err)
	}
	if err := c.Chain.Load(basePath, mngr); err != nil {
		return fmt.Errorf("failed to load chain of RPC: %w", err)
	}
	return nil
}

type MethodForward struct {
	Matcher string
	Forward ResourceID[*RPC]
}

type CompositeRPC struct {
	Name        string
	Description string
	Methods     []MethodForward
}

func (c *CompositeRPC) ResourceKind() ResourceKind {
	return KindCompositeRPC
}

func (c *CompositeRPC) Init(basePath string, mngr *ResourceManager) error {
	for _, m := range c.Methods {
		if err := m.Forward.Load(basePath, mngr); err != nil {
			return fmt.Errorf("failed to load forward-to RPC of composite RPC: %w", err)
		}
	}
	return nil
}

type ProviderRPC struct {
	Name     string                    `yaml:"name"`
	Backend  ResourceID[*CompositeRPC] `yaml:"backend"`
	Accounts []ResourceID[*Account]    `yaml:"accounts"`
	UUIDv4   uuid.UUID                 `yaml:"-"`
}

func (c *ProviderRPC) ResourceKind() ResourceKind {
	return KindProviderRPC
}

func (c *ProviderRPC) Init(basePath string, mngr *ResourceManager) error {
	if err := c.Backend.Load(basePath, mngr); err != nil {
		return fmt.Errorf("failed to load provider RPC backend: %w", err)
	}
	for i, id := range c.Accounts {
		if err := id.Load(basePath, mngr); err != nil {
			return fmt.Errorf("failed to load provider account %d: %w", i, err)
		}
	}
	c.UUIDv4 = uuid.New()
	return nil
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

func (c *Account) Init(basePath string, mngr *ResourceManager) error {
	return nil
}

type Wallet struct {
	Name        string
	Description string
	Accounts    []ResourceID[*Account]
}

func (c *Wallet) ResourceKind() ResourceKind {
	return KindWallet
}

func (c *Wallet) Init(basePath string, mngr *ResourceManager) error {
	for i, acc := range c.Accounts {
		if err := acc.Load(basePath, mngr); err != nil {
			return fmt.Errorf("failed to load account %d of wallet: %w", i, err)
		}
	}
	return nil
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

func (c *Signer) Init(basePath string, mngr *ResourceManager) error {
	if err := c.Hardware.Load(basePath, mngr); err != nil {
		return fmt.Errorf("failed to load hardware of signer: %w", err)
	}
	return nil
}

type Hardware struct {
	Name        string
	Description string
	Type        string
}

func (c *Hardware) ResourceKind() ResourceKind {
	return KindHardware
}

func (c *Hardware) Init(basePath string, mngr *ResourceManager) error {
	return nil
}

type ContractABI struct {
	Name    string
	Address string
	Type    string
}

func (c *ContractABI) ResourceKind() ResourceKind {
	return KindContractABI
}

func (c *ContractABI) Init(basePath string, mngr *ResourceManager) error {
	return nil
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

func (c *TokenERC20) Init(basePath string, mngr *ResourceManager) error {
	if err := c.Chain.Load(basePath, mngr); err != nil {
		return fmt.Errorf("faileed to load chain of erc20 token: %w", err)
	}
	return nil
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

func (c *TokenERC721) Init(basePath string, mngr *ResourceManager) error {
	if err := c.Chain.Load(basePath, mngr); err != nil {
		return fmt.Errorf("faileed to load chain of erc721 token: %w", err)
	}
	return nil
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

func (c *TokenERC1155) Init(basePath string, mngr *ResourceManager) error {
	if err := c.Chain.Load(basePath, mngr); err != nil {
		return fmt.Errorf("faileed to load chain of erc1155 token: %w", err)
	}
	return nil
}

type Holding struct {
	Name        string
	Description string
	Account     ResourceID[*Account]
	Token       ResourceID[*TokenERC20] `yaml:"token,omitempty"`
	Currency    ResourceID[*Currency]   `yaml:"currency,omitempty"`
}

func (c *Holding) ResourceKind() ResourceKind {
	return KindHolding
}

func (c *Holding) Init(basePath string, mngr *ResourceManager) error {
	if err := c.Account.Load(basePath, mngr); err != nil {
		return fmt.Errorf("faileed to load account of holding: %w", err)
	}
	if c.Token != "" {
		if err := c.Token.Load(basePath, mngr); err != nil {
			return fmt.Errorf("faileed to load erc20 token of holding: %w", err)
		}
	}
	if c.Currency != "" {
		if err := c.Currency.Load(basePath, mngr); err != nil {
			return fmt.Errorf("faileed to load currency of holding: %w", err)
		}
	}
	if c.Token == "" && c.Currency == "" {
		return fmt.Errorf("holding must have either erc20 token or currency")
	}
	if c.Token != "" && c.Currency != "" {
		return fmt.Errorf("holding cannot both have erc20 token and currency")
	}
	return nil
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

func (c *Script) Init(basePath string, mngr *ResourceManager) error {
	for _, dep := range c.RPCs {
		if err := dep.Chain.Load(basePath, mngr); err != nil {
			return fmt.Errorf("failed to load chain of script dep %q: %w", dep.K, err)
		}
	}
	return nil
}

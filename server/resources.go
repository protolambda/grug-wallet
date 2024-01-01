package server

import (
	"context"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/fs"
	"path"
)

type ResourceMap[R Resource] map[ResourceID[R]]R

type ResourceManager struct {
	FS fs.FS // TODO: FS should provide special path prefix for global things (to deduplicate common data)

	Currencies    ResourceMap[*Currency]
	Explorers     ResourceMap[*Explorer]
	Chains        ResourceMap[*Chain]
	RPCVendor     ResourceMap[*RPCVendor]
	RPCs          ResourceMap[*RPC]
	CompositeRPCs ResourceMap[*CompositeRPC]
	ProviderRPCs  ResourceMap[*ProviderRPC]
	Accounts      ResourceMap[*Account]
	Wallets       ResourceMap[*Wallet]
	Signers       ResourceMap[*Signer]
	Hardware      ResourceMap[*Hardware]
	ContractABIs  ResourceMap[*ContractABI]
	TokenERC20s   ResourceMap[*TokenERC20]
	TokenERC721s  ResourceMap[*TokenERC721]
	TokenERC1155s ResourceMap[*TokenERC1155]
	Holdings      ResourceMap[*Holding]
	Scripts       ResourceMap[*Script]
}

// Load loads the resource, and updates the resource-ID to an absolute ID
func (id *ResourceID[R]) Load(basePath string, resMngr *ResourceManager) error {
	if id == nil {
		return fmt.Errorf("cannot load nil resource! (type: %T)", id)
	}

	// update ID to resolve against base path (if not already absolute path)
	if !path.IsAbs(string(*id)) {
		*id = ResourceID[R](path.Join(basePath, string(*id)))
	}

	if r := id.get(resMngr); r != nil {
		return nil
	}
	if err := id.load(resMngr); err != nil {
		return fmt.Errorf("failed to load %q (%s): %w", string(*id), id.Kind(), err)
	}
	return nil
}

func (id ResourceID[R]) Kind() ResourceKind {
	var r R
	return r.ResourceKind()
}

func loadRawResource(resMngr *ResourceManager, filename string, dest any) error {
	f, err := resMngr.FS.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open resource file: %w", err)
	}
	defer f.Close()

	dec := yaml.NewDecoder(f)
	dec.KnownFields(true)
	if err := dec.Decode(dest); err != nil {
		return fmt.Errorf("failed to decode YAML: %w", err)
	}
	return nil
}

func (id ResourceID[R]) load(resMngr *ResourceManager) error {
	var r R
	if err := loadRawResource(resMngr, string(id), &r); err != nil {
		return err
	}

	if err := r.Init(path.Dir(string(id)), resMngr); err != nil {
		return fmt.Errorf("failed to init: %w", err)
	}

	var v any = r
	switch x := v.(type) {
	case *Currency:
		resMngr.Currencies[ResourceID[*Currency](id)] = x
	case *Explorer:
		resMngr.Explorers[ResourceID[*Explorer](id)] = x
	case *Chain:
		resMngr.Chains[ResourceID[*Chain](id)] = x
	case *RPCVendor:
		resMngr.RPCVendor[ResourceID[*RPCVendor](id)] = x
	case *RPC:
		resMngr.RPCs[ResourceID[*RPC](id)] = x
	case *CompositeRPC:
		resMngr.CompositeRPCs[ResourceID[*CompositeRPC](id)] = x
	case *ProviderRPC:
		resMngr.ProviderRPCs[ResourceID[*ProviderRPC](id)] = x
	case *Account:
		resMngr.Accounts[ResourceID[*Account](id)] = x
	case *Wallet:
		resMngr.Wallets[ResourceID[*Wallet](id)] = x
	case *Signer:
		resMngr.Signers[ResourceID[*Signer](id)] = x
	case *Hardware:
		resMngr.Hardware[ResourceID[*Hardware](id)] = x
	case *ContractABI:
		resMngr.ContractABIs[ResourceID[*ContractABI](id)] = x
	case *TokenERC20:
		resMngr.TokenERC20s[ResourceID[*TokenERC20](id)] = x
	case *TokenERC721:
		resMngr.TokenERC721s[ResourceID[*TokenERC721](id)] = x
	case *TokenERC1155:
		resMngr.TokenERC1155s[ResourceID[*TokenERC1155](id)] = x
	case *Holding:
		resMngr.Holdings[ResourceID[*Holding](id)] = x
	case *Script:
		resMngr.Scripts[ResourceID[*Script](id)] = x
	default:
		panic(fmt.Errorf("unrecognized type %T", x))
	}
	return nil
}

type resManagerKeyType struct{}

var resManagerKey = resManagerKeyType{}

func WithResourceManager(ctx context.Context, resMngr *ResourceManager) context.Context {
	return context.WithValue(ctx, resManagerKey, resMngr)
}

func ResourceManagerFromContext(ctx context.Context) *ResourceManager {
	v := ctx.Value(resManagerKey)
	if v == nil {
		return nil
	}
	return v.(*ResourceManager)
}

func (id ResourceID[R]) Get(ctx context.Context) R {
	resMngr := ResourceManagerFromContext(ctx)
	return id.get(resMngr)
}

func (id ResourceID[R]) get(resMngr *ResourceManager) R {
	// Yes, this is a very hacky way to work around
	// Go generics type-switching limitations, oh well.
	r := new(R)
	var v any = r
	switch x := v.(type) {
	case **Currency:
		*x = resMngr.Currencies[ResourceID[*Currency](id)]
	case **Explorer:
		*x = resMngr.Explorers[ResourceID[*Explorer](id)]
	case **Chain:
		*x = resMngr.Chains[ResourceID[*Chain](id)]
	case **RPCVendor:
		*x = resMngr.RPCVendor[ResourceID[*RPCVendor](id)]
	case **RPC:
		*x = resMngr.RPCs[ResourceID[*RPC](id)]
	case **CompositeRPC:
		*x = resMngr.CompositeRPCs[ResourceID[*CompositeRPC](id)]
	case **ProviderRPC:
		*x = resMngr.ProviderRPCs[ResourceID[*ProviderRPC](id)]
	case **Account:
		*x = resMngr.Accounts[ResourceID[*Account](id)]
	case **Wallet:
		*x = resMngr.Wallets[ResourceID[*Wallet](id)]
	case **Signer:
		*x = resMngr.Signers[ResourceID[*Signer](id)]
	case **Hardware:
		*x = resMngr.Hardware[ResourceID[*Hardware](id)]
	case **ContractABI:
		*x = resMngr.ContractABIs[ResourceID[*ContractABI](id)]
	case **TokenERC20:
		*x = resMngr.TokenERC20s[ResourceID[*TokenERC20](id)]
	case **TokenERC721:
		*x = resMngr.TokenERC721s[ResourceID[*TokenERC721](id)]
	case **TokenERC1155:
		*x = resMngr.TokenERC1155s[ResourceID[*TokenERC1155](id)]
	case **Holding:
		*x = resMngr.Holdings[ResourceID[*Holding](id)]
	case **Script:
		*x = resMngr.Scripts[ResourceID[*Script](id)]
	default:
		panic(fmt.Errorf("unrecognized type %T", x))
	}
	return *r
}

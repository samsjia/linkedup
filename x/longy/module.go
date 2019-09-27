package longy

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/eco/longy/x/longy/internal"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	//_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic is
type AppModuleBasic struct{}

//Name returns the name of the module
func (a AppModuleBasic) Name() string {
	return ModuleName
}

//RegisterCodec registers module to the codec
func (a AppModuleBasic) RegisterCodec(*codec.Codec) {

}

//DefaultGenesis returns the default genesis for this module if any
func (a AppModuleBasic) DefaultGenesis() json.RawMessage {
	return nil
}

//ValidateGenesis validates that the json genesis is valid to our module
func (a AppModuleBasic) ValidateGenesis(json.RawMessage) error {
	return nil
}

//RegisterRESTRoutes registers our module rest endpoints
func (a AppModuleBasic) RegisterRESTRoutes(context.CLIContext, *mux.Router) {

}

//GetTxCmd returns any tx commands from this module to the parent command in the cli
func (a AppModuleBasic) GetTxCmd(*codec.Codec) *cobra.Command {
	return nil
}

//GetQueryCmd returns any query commands from this module to the parent command in the cli
func (a AppModuleBasic) GetQueryCmd(*codec.Codec) *cobra.Command {
	return nil
}

// AppModule structure holding or keepers together
type AppModule struct {
	AppModuleBasic
	accountKeeper types.AccountKeeper
	bankKeeper    bank.Keeper
	longyKeeper   internal.Keeper
}

// NewAppModule creates a new AppModule object
// nolint: gocritic
func NewAppModule(account types.AccountKeeper, bank bank.Keeper, longy Keeper) module.AppModule {

	return module.NewGenesisOnlyAppModule(AppModule{
		AppModuleBasic: AppModuleBasic{},
		accountKeeper:  account,
		bankKeeper:     bank,
		longyKeeper:    longy,
	})
}

// InitGenesis init-genesis
// nolint: gocritic
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// ExportGenesis export genesis
// nolint: gocritic
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	return make([]byte, 0)
}

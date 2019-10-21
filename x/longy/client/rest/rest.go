package rest

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	"github.com/eco/longy/x/longy/client/rest/query"
	"github.com/eco/longy/x/longy/internal/querier"
	"github.com/gorilla/mux"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
//nolint:gocritic
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	//longy/attendees/{attendee_id}
	r.HandleFunc(fmt.Sprintf("/%s/%s/{%s}", storeName, querier.QueryAttendees, query.AttendeeIDKey),
		attendeeHandler(cliCtx, storeName)).Methods("GET")

	//longy/attendees/address/{address_id}
	r.HandleFunc(fmt.Sprintf("/%s/%s/%s/{%s}", storeName, querier.QueryAttendees, querier.AddressKey,
		query.AddressIDKey), attendeeAddressHandler(cliCtx, storeName)).Methods("GET")

	//longy/scans/{scan_id}
	r.HandleFunc(fmt.Sprintf("/%s/%s/{%s}", storeName, querier.QueryScans, query.ScanIDKey),
		scanGetHandler(cliCtx, storeName)).Methods("GET")

	//longy/prizes
	r.HandleFunc(fmt.Sprintf("/%s/%s", storeName, querier.PrizesKey),
		prizesGetHandler(cliCtx, storeName)).Methods("GET")

	//longy/bonus
	r.HandleFunc(fmt.Sprintf("/%s/%s", storeName, querier.QueryBonus),
		bonusGetHandler(cliCtx, storeName)).Methods("GET")

	//longy/leader
	r.HandleFunc(fmt.Sprintf("/%s/%s", storeName, querier.LeaderKey),
		query.LeaderBoardHandler(cliCtx, storeName)).Methods("GET")

	//longy/redeem?address_id={address_id}
	r.HandleFunc(fmt.Sprintf("/%s/%s", storeName, querier.RedeemKey),
		query.RedeemHandler(cliCtx, storeName)).
		Queries(query.AddressIDKey, fmt.Sprintf("{%s}", query.AddressIDKey)).Methods("GET")

	//open endpoint to post transactions directly to full node
	r.HandleFunc("/longy/txs", rest.BroadcastTxRequest(cliCtx)).Methods("POST")
}

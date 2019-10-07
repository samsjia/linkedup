package handler

import (
	"github.com/eco/longy/x/longy/internal/keeper"
	"github.com/eco/longy/x/longy/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// HandleMsgQrScan processes MsgScanQr message
//nolint:gocritic
func HandleMsgQrScan(ctx sdk.Context, k keeper.Keeper, msg types.MsgScanQr) sdk.Result {
	//get the address for the scanned qr code
	attendee, ok := k.GetAttendeeWithID(ctx, msg.ScannedQR)
	if !ok {
		return types.ErrAttendeeNotFound("cannot find the attendee").Result()
	}
	//get the id for the scan event
	id, err := types.GenScanID(msg.Sender, attendee.Address)
	if err != nil {
		return err.Result()
	}

	//get the scan event
	scan, err := k.GetScanByID(ctx, id)
	if err != nil { //if new scan, create it
		scan, err = handleNewScan(ctx, k, msg, attendee)
		if err != nil {
			return err.Result()
		}
	}

	err = handleShareInfo(ctx, k, scan, msg.Sender, attendee, msg.Data)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

//nolint:gocritic
func handleShareInfo(ctx sdk.Context, k keeper.Keeper, scan types.Scan, sender sdk.AccAddress,
	attendee types.Attendee, data []byte) sdk.Error {
	//add share ids, skips if the ids are already added
	err := k.AddSharedID(ctx, sender, attendee.Address, scan.ID)
	if err != nil {
		return err
	}

	//check who is in what position
	var oldData *[]byte
	if scan.S1.Equals(sender) {
		oldData = &scan.D1
	} else {
		oldData = &scan.D2
	}

	if len(*oldData) == 0 && len(data) > 0 {
		err := k.AwardShareInfoPoints(ctx, sender, attendee.Address)
		if err != nil {
			return err
		}

		//set new data into scan and save scan
		*oldData = data
		k.SetScan(ctx, &scan)
	}
	return nil
}

//nolint:gocritic
func handleNewScan(ctx sdk.Context, k keeper.Keeper, msg types.MsgScanQr,
	attendee types.Attendee) (scan types.Scan, err sdk.Error) {
	scan, err = types.NewScan(msg.Sender, attendee.Address, nil, nil) //dont pass data here

	if err != nil {
		return
	}
	err = k.AwardScanPoints(ctx, scan)
	if err != nil {
		return
	}

	k.SetScan(ctx, &scan)
	return
}

package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = MsgRedeem{}

// MsgRedeem is used to claim prizes by the booth operators
type MsgRedeem struct {
	Sender    sdk.AccAddress `json:"sender"`    //Standard for all messages
	ScannedQR string         `json:"scannedQR"` //the string representation of the other attendee's QR badge
}

// NewMsgRedeem in the constructor for `MsgRedeem`
func NewMsgRedeem(sender sdk.AccAddress, scannedQr string) MsgRedeem {
	return MsgRedeem{
		Sender:    sender,
		ScannedQR: scannedQr,
	}
}

// Route defines the route for this message
//nolint:gocritic
func (msg MsgRedeem) Route() string {
	return RouterKey
}

// Type is the message type
//nolint:gocritic
func (msg MsgRedeem) Type() string {
	return "claim_prize"
}

// ValidateBasic performs sanity checks on the message
//nolint:gocritic
func (msg MsgRedeem) ValidateBasic() sdk.Error {
	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}

	if !ValidQrCode(msg.ScannedQR) {
		return ErrQRCodeInvalid("message QR code is invalid, should be a string of a positive integer")
	}

	return nil
}

// GetSignBytes returns the byte array that is signed over
//nolint:gocritic
func (msg MsgRedeem) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns the
//nolint:gocritic
func (msg MsgRedeem) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

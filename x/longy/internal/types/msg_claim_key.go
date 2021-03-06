package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = MsgClaimKey{}

// MsgClaimKey is used to claim a prior rekey message
type MsgClaimKey struct {
	AttendeeAddress sdk.AccAddress `json:"attendeeAddress"`
	Name            string         `json:"name"`
	Secret          string         `json:"secret"`
	RsaPublicKey    string         `json:"rsaPublicKey"`
	EncryptedInfo   []byte         `json:"encryptedInfo"`
}

// NewMsgClaimKey in the constructor for `MsgClaimKey`
func NewMsgClaimKey(address sdk.AccAddress, name string, secret string, rsaPublicKey string,
	encryptedInfo []byte) MsgClaimKey {
	return MsgClaimKey{
		AttendeeAddress: address,
		Name:            name,
		Secret:          secret,
		RsaPublicKey:    rsaPublicKey,
		EncryptedInfo:   encryptedInfo,
	}
}

// Route defines the route for this message
//nolint:gocritic
func (msg MsgClaimKey) Route() string {
	return RouterKey
}

// Type is the message type
//nolint:gocritic
func (msg MsgClaimKey) Type() string {
	return "claim_key"
}

// ValidateBasic performs sanity checks on the message
//nolint:gocritic
func (msg MsgClaimKey) ValidateBasic() sdk.Error {
	switch {
	case msg.AttendeeAddress.Empty():
		return sdk.ErrInvalidAddress("empty attendee address")
	case msg.Name == "":
		return ErrEmptyName("name cannot be empty")
	case len(msg.Secret) == 0:
		return ErrEmptySecret("secret cannot be empty")
	case len(msg.RsaPublicKey) == 0:
		return ErrEmptyRsaKey("rsa public key cannot be empty")
	case len(msg.EncryptedInfo) == 0:
		return ErrEmptyEncryptedInfo("encrypted info cannot be empty")
	default:
		return nil
	}
}

// GetSignBytes returns the byte array that is signed over
//nolint:gocritic
func (msg MsgClaimKey) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns the
//nolint:gocritic
func (msg MsgClaimKey) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.AttendeeAddress}
}

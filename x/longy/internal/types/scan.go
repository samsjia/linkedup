package types

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//Scan represents an unique scan between two parties, can only be one scan between the same parties
type Scan struct {
	//ID is the key we use to store this struct in the keyStore, it is unique per S1-S2 pair
	ID []byte
	//S1 is the scanner that initiates the scan
	S1 sdk.AccAddress
	//S2 is the person who's QR code is scanned
	S2 sdk.AccAddress
	//D1 is the encrypted data shared by S1 with S2
	D1 []byte
	//D2 is the encrypted data shared by S2 with S1
	D2 []byte
	//P1 are the points earned by S1 for a scan
	P1 uint
	//P2 are the points earned by S2 for a scan
	P2 uint
	//UnixTimeSec is the unix time in seconds of the block header of when this scan was created
	UnixTimeSec int64
}

//NewScan creates a new scan and sets its id
func NewScan(s1 sdk.AccAddress, s2 sdk.AccAddress, d1 []byte, d2 []byte, p1 uint, p2 uint) (*Scan, sdk.Error) {
	id, err := GenScanID(s1, s2)
	if err != nil {
		return &Scan{}, err
	}
	return &Scan{
		ID: id,
		S1: s1,
		S2: s2,
		D1: d1,
		D2: d2,
		P1: p1,
		P2: p2,
	}, nil
}

//AddPoints adds points to the s1 and s2 respectively
func (s *Scan) AddPoints(p1 uint, p2 uint) {
	s.P1 += p1
	s.P2 += p2
}

//AddPointsToAccount Adds points to the given account, assumes address is one of S1 or S2
func (s *Scan) AddPointsToAccount(address sdk.AccAddress, points uint) {
	if address.Equals(s.S1) {
		s.AddPoints(points, 0)
	} else {
		s.AddPoints(0, points)
	}
}

//SetTimeUnixSeconds sets the unix time in seconds of the block header of when this scan was created
func (s *Scan) SetTimeUnixSeconds(unix int64) {
	s.UnixTimeSec = unix
}

//GenScanID creates the unique id between a scan pair, regardless of the order of the account addresses passed into it
func GenScanID(s1, s2 sdk.AccAddress) (id []byte, err sdk.Error) {
	err = CheckSameness(s1, s2)
	if err != nil {
		return
	}

	val := bytes.Compare(s1, s2)

	//nolint:gocritic
	if val > 0 {
		id = append(s1, s2...)
	} else {
		id = append(s2, s1...)
	}
	//append the key so we dont have to do this everywhere
	id = ScanKey(id)
	return
}

//CheckSameness checks the address are not empty or the same
//nolint:interfacer
func CheckSameness(s1 sdk.AccAddress, s2 sdk.AccAddress) (err sdk.Error) {
	if s1.Empty() || s2.Empty() {
		err = ErrAccountAddressEmpty("cannot have empty addresses")
		return
	}

	if s1.Equals(s2) {
		err = ErrAccountsSame("addresses cannot be the same")
		return
	}

	return
}

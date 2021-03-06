package masterkey

import (
	"errors"
	"fmt"
	"sync"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	longyApp "github.com/eco/longy"
	longyClnt "github.com/eco/longy/key-service/longyclient"
	"github.com/eco/longy/util"
	"github.com/eco/longy/x/longy"
	"github.com/sirupsen/logrus"
	tmcrypto "github.com/tendermint/tendermint/crypto"
)

var log = logrus.WithField("module", "masterkey")

var (
	// ErrAlreadyKeyed denotes that this address has already been key'd
	ErrAlreadyKeyed = errors.New("account already key'ed")
)

// MasterKey encapslates the master key for the longy game
type MasterKey struct {
	privKey tmcrypto.PrivKey
	pubKey  tmcrypto.PubKey
	address sdk.AccAddress

	chainID string

	accNum      uint64
	sequenceNum uint64
	seqLock     *sync.Mutex

	cdc *codec.Codec
}

// NewMasterKey is the constructor for `Key`. A new secp256k1 is generated if empty.
// The `chainID` is used when generating RekeyTransactions to prevent cross-chain replay attacks
func NewMasterKey(privateKey tmcrypto.PrivKey, chainID string) (MasterKey, error) {

	// retrieve details about the master account from the rest endpoint
	sdkAddr := sdk.AccAddress(privateKey.PubKey().Address())
	masterAccount, err := longyClnt.GetAccount(sdkAddr)
	if err != nil {
		return MasterKey{}, fmt.Errorf("masterkey account retrieval: %s", err)
	}

	k := MasterKey{
		privKey: privateKey,
		pubKey:  privateKey.PubKey(),
		address: sdkAddr,

		chainID: chainID,

		accNum:      masterAccount.GetAccountNumber(),
		sequenceNum: masterAccount.GetSequence(),
		seqLock:     &sync.Mutex{},

		cdc: longyApp.MakeCodec(),
	}

	log.Infof("constructed master key. Chain-Id=%s, AccountNum=%d, SequenceNum=%d", k.chainID, k.accNum, k.sequenceNum)
	return k, nil
}

// SendKeyTransaction generates a `RekeyMsg`, authorized by the master key. The transaction bytes
// generated are created using the cosmos-sdk/x/auth module's StdSignDoc.
func (mk *MasterKey) SendKeyTransaction(
	attendeeAddr sdk.AccAddress,
	newPublicKey tmcrypto.PubKey,
	commitment util.Commitment,
) error {

	/** Block until we submit the transaction **/
	mk.seqLock.Lock()

	// create and broadcast the transaction
	keyMsg := longy.NewMsgKey(attendeeAddr, mk.address, newPublicKey, commitment)
	tx := mk.createKeyTx(keyMsg)
	res, err := longyClnt.BroadcastAuthTx(tx, "block")
	if err != nil { // nolint
		log.WithError(err).Info("failed transaction submission")
	} else {
		if res.Code != 0 {
			if res.Code == uint32(longy.CodeAttendeeKeyed) {
				err = ErrAlreadyKeyed
			} else {
				log.WithField("raw_log", res.RawLog).Info("failed tx response")
				err = fmt.Errorf("failed tx")
			}
		}

		mk.sequenceNum++
	}

	mk.seqLock.Unlock()

	return err
}

//nolint
func (mk *MasterKey) createKeyTx(keyMsg longy.MsgKey) *auth.StdTx {
	msgs := []sdk.Msg{keyMsg}

	nilFee := auth.NewStdFee(50000, sdk.NewCoins(sdk.NewInt64Coin("longy", 0)))
	signBytes := auth.StdSignBytes(mk.chainID, mk.accNum, mk.sequenceNum, nilFee, msgs, "")

	// sign the message with the master private key
	sig, err := mk.privKey.Sign(signBytes)
	if err != nil {
		panic(err)
	}
	stdSig := auth.StdSignature{PubKey: mk.pubKey, Signature: sig}
	tx := auth.NewStdTx(msgs, nilFee, []auth.StdSignature{stdSig}, "")

	return &tx
}

// Copyright Monax Industries Limited
// SPDX-License-Identifier: Apache-2.0

package mock

import (
	"fmt"

	"github.com/hyperledger/burrow/acm"
	"github.com/hyperledger/burrow/crypto"
	"github.com/hyperledger/burrow/keys"
)

//---------------------------------------------------------------------
// Mock client for replacing signing done by monax-keys

// Implementation assertion
var _ keys.KeyClient = (*KeyClient)(nil)

type KeyClient struct {
	knownKeys map[crypto.Address]*Key
}

func NewKeyClient(privateAccounts ...*acm.PrivateAccount) *KeyClient {
	client := &KeyClient{
		knownKeys: make(map[crypto.Address]*Key),
	}
	for _, pa := range privateAccounts {
		client.knownKeys[pa.GetAddress()] = mockKeyFromPrivateAccount(pa)
	}
	return client
}

func (mkc *KeyClient) NewKey(name string) crypto.Address {
	// Only tests ED25519 curve and ripemd160.
	key, err := newKey(name)
	if err != nil {
		panic(fmt.Sprintf("Mocked key client failed on key generation: %s", err))
	}
	mkc.knownKeys[key.Address] = key
	return key.Address
}

func (mkc *KeyClient) Sign(signAddress crypto.Address, message []byte) (*crypto.Signature, error) {
	key := mkc.knownKeys[signAddress]
	if key == nil {
		return nil, fmt.Errorf("unknown address (%s)", signAddress)
	}
	return key.Sign(message)
}

func (mkc *KeyClient) PublicKey(address crypto.Address) (crypto.PublicKey, error) {
	key := mkc.knownKeys[address]
	if key == nil {
		return crypto.PublicKey{}, fmt.Errorf("unknown address (%s)", address)
	}
	return crypto.PublicKeyFromBytes(key.PublicKey, crypto.CurveTypeEd25519)
}

func (mkc *KeyClient) Generate(keyName string, curve crypto.CurveType) (crypto.Address, error) {
	return mkc.NewKey(keyName), nil
}

func (mkc *KeyClient) GetAddressForKeyName(keyName string) (crypto.Address, error) {
	for _, m := range mkc.knownKeys {
		if m.Name == keyName {
			return m.Address, nil
		}
	}

	return crypto.Address{}, nil
}

func (mkc *KeyClient) HealthCheck() error {
	return nil
}

func (mkc *KeyClient) Keys() []*Key {
	var knownKeys []*Key
	for _, key := range mkc.knownKeys {
		knownKeys = append(knownKeys, key)
	}
	return knownKeys
}

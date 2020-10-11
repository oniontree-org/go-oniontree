package oniontree

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/openpgp"
	"strings"
	"unicode"
)

// The type implements openpgp.KeyRing interface.
type PublicKeys []*PublicKey

func (pks PublicKeys) getEntities() openpgp.EntityList {
	el := make(openpgp.EntityList, 0, len(pks))
	for i := range pks {
		ets, err := openpgp.ReadArmoredKeyRing(strings.NewReader(pks[i].Value))
		if err != nil {
			return nil
		}
		el = append(el, ets...)
	}
	return el
}

// KeysById returns the set of keys that have the given key id.
func (pks PublicKeys) KeysById(id uint64) []openpgp.Key {
	el := pks.getEntities()
	return el.KeysById(id)
}

// KeysByIdAndUsage returns the set of keys with the given id
// that also meet the key usage given by requiredUsage.
// The requiredUsage is expressed as the bitwise-OR of
// packet.KeyFlag* values.
func (pks PublicKeys) KeysByIdUsage(id uint64, requiredUsage byte) []openpgp.Key {
	el := pks.getEntities()
	return el.KeysByIdUsage(id, requiredUsage)
}

// DecryptionKeys returns all private keys that are valid for
// decryption.
func (pks PublicKeys) DecryptionKeys() []openpgp.Key {
	el := pks.getEntities()
	return el.DecryptionKeys()
}

type PublicKey struct {
	ID          string `json:"id,omitempty" yaml:"id,omitempty"`
	UserID      string `json:"user_id,omitempty" yaml:"user_id,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty" yaml:"fingerprint,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Value       string `json:"value" yaml:"value"`
}

func NewPublicKey(b []byte) (*PublicKey, error) {
	bClean := bytes.TrimLeftFunc(b, unicode.IsSpace)
	el, err := openpgp.ReadArmoredKeyRing(bytes.NewReader(bClean))
	if err != nil {
		return nil, err
	}

	publicKey := &PublicKey{}
	for _, e := range el {
		userID := ""
		for _, ident := range e.Identities {
			userID = ident.Name
		}
		pk := e.PrimaryKey
		publicKey = &PublicKey{
			Value:       string(bClean),
			ID:          pk.KeyIdString(),
			Fingerprint: fmt.Sprintf("%X", pk.Fingerprint),
			UserID:      userID,
		}
	}
	return publicKey, nil
}

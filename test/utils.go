package testutil

import (
	"math/rand"
	"time"

	// TODO: change back to libp2p for pr
	ci "github.com/c3systems/go-libp2p-crypto"
)

func RandTestKeyPair(typ, bits int) (ci.PrivKey, ci.PubKey, error) {
	return SeededTestKeyPair(typ, bits, time.Now().UnixNano())
}

func SeededTestKeyPair(typ, bits int, seed int64) (ci.PrivKey, ci.PubKey, error) {
	r := rand.New(rand.NewSource(seed))
	return ci.GenerateKeyPairWithReader(typ, bits, r)
}

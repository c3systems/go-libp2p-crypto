package crypto

import (
	"fmt"
	"io"

	// note: changeback to libp2p for pr
	pb "github.com/c3systems/go-libp2p-crypto/pb"

	btcec "github.com/btcsuite/btcd/btcec"
	proto "github.com/gogo/protobuf/proto"
	sha256 "github.com/minio/sha256-simd"
)

// Secp256k1PrivateKey is an Secp256k1 private key
type Secp256k1PrivateKey btcec.PrivateKey

// Secp256k1PublicKey is an Secp256k1 public key
type Secp256k1PublicKey btcec.PublicKey

// GenerateSecp256k1Key generates a new Secp256k1 private and public key pair
func GenerateSecp256k1Key(src io.Reader) (PrivKey, PubKey, error) {
	privk, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return nil, nil, err
	}

	k := (*Secp256k1PrivateKey)(privk)
	return k, k.GetPublic(), nil
}

// UnmarshalSecp256k1PrivateKey returns a private key from bytes
func UnmarshalSecp256k1PrivateKey(data []byte) (PrivKey, error) {
	if len(data) != btcec.PrivKeyBytesLen {
		return nil, fmt.Errorf("expected secp256k1 data size to be %d", btcec.PrivKeyBytesLen)
	}

	privk, _ := btcec.PrivKeyFromBytes(btcec.S256(), data)
	return (*Secp256k1PrivateKey)(privk), nil
}

// UnmarshalSecp256k1PublicKey returns a public key from bytes
func UnmarshalSecp256k1PublicKey(data []byte) (PubKey, error) {
	k, err := btcec.ParsePubKey(data, btcec.S256())
	if err != nil {
		return nil, err
	}

	return (*Secp256k1PublicKey)(k), nil
}

// Bytes returns protobuf bytes from a private key
func (k *Secp256k1PrivateKey) Bytes() ([]byte, error) {
	pbmes := new(pb.PrivateKey)
	typ := pb.KeyType_Secp256k1
	pbmes.Type = &typ
	pbmes.Data = (*btcec.PrivateKey)(k).Serialize()
	return proto.Marshal(pbmes)
}

// Equals compares two private keys
func (k *Secp256k1PrivateKey) Equals(o Key) bool {
	sk, ok := o.(*Secp256k1PrivateKey)
	if !ok {
		return false
	}

	return k.D.Cmp(sk.D) == 0
}

// Sign returns a signature from input data
func (k *Secp256k1PrivateKey) Sign(data []byte) ([]byte, error) {
	hash := sha256.Sum256(data)
	sig, err := (*btcec.PrivateKey)(k).Sign(hash[:])
	if err != nil {
		return nil, err
	}

	return sig.Serialize(), nil
}

// GetPublic returns a public key
func (k *Secp256k1PrivateKey) GetPublic() PubKey {
	return (*Secp256k1PublicKey)((*btcec.PrivateKey)(k).PubKey())
}

// Bytes returns protobuf bytes from a public key
func (k *Secp256k1PublicKey) Bytes() ([]byte, error) {
	pbmes := new(pb.PublicKey)
	typ := pb.KeyType_Secp256k1
	pbmes.Type = &typ
	pbmes.Data = (*btcec.PublicKey)(k).SerializeCompressed()
	return proto.Marshal(pbmes)
}

// Equals compares two public keys
func (k *Secp256k1PublicKey) Equals(o Key) bool {
	sk, ok := o.(*Secp256k1PublicKey)
	if !ok {
		return false
	}

	return (*btcec.PublicKey)(k).IsEqual((*btcec.PublicKey)(sk))
}

// Verify compares a signature against the input data
func (k *Secp256k1PublicKey) Verify(data []byte, sigStr []byte) (bool, error) {
	sig, err := btcec.ParseDERSignature(sigStr, btcec.S256())
	if err != nil {
		return false, err
	}

	hash := sha256.Sum256(data)
	return sig.Verify(hash[:], (*btcec.PublicKey)(k)), nil
}

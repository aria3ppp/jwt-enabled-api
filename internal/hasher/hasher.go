package hasher

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

type Interface interface {
	GenerateHash(value []byte) ([]byte, error)
	CompareHash(hash, value []byte) error
}

func NewBcrypt(cost int) Bcrypt {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		log.Panicf("invalid cost %d", cost)
	}
	return Bcrypt{
		cost: cost,
	}
}

type Bcrypt struct {
	cost int
}

var _ Interface = (Bcrypt{})

func (b Bcrypt) GenerateHash(value []byte) (hash []byte, err error) {
	return bcrypt.GenerateFromPassword(value, b.cost)
}

func (b Bcrypt) CompareHash(hash, value []byte) error {
	err := bcrypt.CompareHashAndPassword(hash, value)
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return ErrMismatchedHash
		}
		return err
	}
	return nil
}

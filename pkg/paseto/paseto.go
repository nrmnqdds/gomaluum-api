package paseto

import (
	"log"
	"os"

	"aidanwoods.dev/go-paseto"
)

type AppPaseto struct {
	PublicKey  *paseto.V4AsymmetricPublicKey
	PrivateKey *paseto.V4AsymmetricSecretKey
	Token      *paseto.Token
}

func New() *AppPaseto {
	publicKey, err := paseto.NewV4AsymmetricPublicKeyFromHex(os.Getenv("PASETO_PUBLIC_KEY"))
	if err != nil {
		log.Fatalf("Failed to create public key: %v", err)
		return nil
	}

	privateKey, err := paseto.NewV4AsymmetricSecretKeyFromHex(os.Getenv("PASETO_SECRET_KEY"))
	if err != nil {
		log.Fatalf("Failed to create secret key: %v", err)
		return nil
	}

	token := paseto.NewToken()
	return &AppPaseto{
		PublicKey:  &publicKey,
		PrivateKey: &privateKey,
		Token:      &token,
	}
}

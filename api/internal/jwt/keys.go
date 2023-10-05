package jwt

import (
	"crypto/rsa"
)

type rsaKeyPair struct {
	name    string
	private *rsa.PrivateKey
	public  *rsa.PublicKey
}

type rsaKeyChain map[string]rsaKeyPair

package common

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"time"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

type JWEEnv struct {
	PrivateKeyFile string `env:"PRIVATE_KEY_JWE" envDefault:""`
}

type JWEGenerator struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func (this *JWEGenerator) GenerateToken(username string) (string, error) {
	signingkey := jose.SigningKey{
		Algorithm: jose.PS512,
		Key:       this.privateKey,
	}

	signer, err := jose.NewSigner(signingkey, nil)
	if err != nil {
		return "", err
	}

	options := &jose.EncrypterOptions{}
	options = options.WithType("JWT").WithContentType("JWT")

	recipient := jose.Recipient{
		Algorithm: jose.RSA_OAEP,
		Key:       this.publicKey,
	}

	encrypter, err := jose.NewEncrypter(jose.A128GCM, recipient, options)
	if err != nil {
		return "", err
	}

	timeout := 7 * 24 * time.Hour
	expired := time.Now().Truncate(time.Second).Add(timeout)

	claims := jwt.Claims{
		Subject: username,
		Issuer:  "JWEGenerator",
		Expiry:  jwt.NewNumericDate(expired),
	}

	jwe, err := jwt.SignedAndEncrypted(signer, encrypter).Claims(claims).CompactSerialize()
	if err != nil {
		return "", err
	}

	return jwe, nil
}

func (this *JWEGenerator) VerifyToken(jwe string) error {
	tok, err := jwt.ParseSignedAndEncrypted(jwe)
	if err != nil {
		return err
	}

	jws, err := tok.Decrypt(this.privateKey)
	if err != nil {
		return err
	}

	var claims jwt.Claims

	err = jws.Claims(this.publicKey, &claims)
	if err != nil {
		return err
	}

	expected := jwt.Expected{
		Time: time.Now(),
	}

	err = claims.Validate(expected)
	if err != nil {
		return err
	}

	return nil
}

func NewJWEGenerator(env *JWEEnv) (*JWEGenerator, error) {
	priv, err := ioutil.ReadFile(env.PrivateKeyFile)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(priv)
	if block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("rsa private key is of the wrong type")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	publicKey := &privateKey.PublicKey

	return &JWEGenerator{
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

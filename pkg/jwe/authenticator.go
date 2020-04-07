package jwe

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"time"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

type Authenticator interface {
	GenerateToken(username string, password string) (string, error)
	VerifyToken(jwe string) error
	VerifyTokenWithUser(jwe string, u User) error
	User() User
}

type authenticator struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	repository UserRepository
	user       User
}

func (a *authenticator) User() User {
	return a.user
}

func (a *authenticator) GenerateToken(username string, password string) (string, error) {
	user, err := a.repository.FindByUsername(username)
	if err != nil {
		return "", errors.New("user does not exists")
	}

	if !user.ComparePasswords(password) {
		return "", errors.New("password incorrect")
	}

	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.PS512, Key: a.privateKey}, nil)
	if err != nil {
		return "", err
	}

	encrypter, err := jose.NewEncrypter(
		jose.A128GCM,
		jose.Recipient{Algorithm: jose.RSA_OAEP, Key: a.publicKey},
		(&jose.EncrypterOptions{}).WithType("JWT").WithContentType("JWT"),
	)

	if err != nil {
		return "", err
	}

	cl := jwt.Claims{
		Subject: username,
		Issuer:  "authenticator",
		Expiry:  jwt.NewNumericDate(time.Now().Add(24 * 7 * time.Hour)),
	}

	user.ClearPassword()

	jwe, err := jwt.SignedAndEncrypted(signer, encrypter).Claims(cl).Claims(user).CompactSerialize()
	if err != nil {
		return "", err
	}

	a.user = user

	return jwe, nil
}

func (a *authenticator) VerifyToken(jwe string) error {
	tok, err := jwt.ParseSignedAndEncrypted(jwe)
	if err != nil {
		return err
	}

	jws, err := tok.Decrypt(a.privateKey)
	if err != nil {
		return err
	}

	cl := jwt.Claims{}
	if err := jws.Claims(a.publicKey, &cl, nil); err != nil {
		return err
	}

	err = cl.Validate(jwt.Expected{
		Time: time.Now(),
	})

	if err != nil {
		return err
	}

	return nil
}

func (a *authenticator) VerifyTokenWithUser(jwe string, u User) error {
	tok, err := jwt.ParseSignedAndEncrypted(jwe)
	if err != nil {
		return err
	}

	jws, err := tok.Decrypt(a.privateKey)
	if err != nil {
		return err
	}

	cl := jwt.Claims{}

	if err := jws.Claims(a.publicKey, &cl, u); err != nil {
		return err
	}

	err = cl.Validate(jwt.Expected{
		Time: time.Now(),
	})

	if err != nil {
		return err
	}

	return nil
}

func NewAuthenticator(repository UserRepository) (*authenticator, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	publicKey := &privateKey.PublicKey

	return &authenticator{privateKey: privateKey, publicKey: publicKey, repository: repository}, nil
}

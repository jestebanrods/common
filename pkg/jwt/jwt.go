package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"assistans-courses/apiserver/internal/models"
	"assistans-courses/apiserver/internal/repositories"

	"github.com/sirupsen/logrus"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

// TODO = Create jwe environments variables
const AuthTimeout = 1 * time.Minute

type JWEAuthenticator interface {
	GenerateToken(username string, pwd string) (string, error)
	VerifyToken(jwe string) (*models.User, error)
	AuthMiddleware(h http.HandlerFunc) http.HandlerFunc
	User() models.User
}

type jweAuthenticator struct {
	privateKey     *rsa.PrivateKey
	publicKey      *rsa.PublicKey
	logger         *logrus.Entry
	userRepository repositories.UsersRepository
	userLogged     models.User
}

func (a jweAuthenticator) User() models.User {
	return a.userLogged
}

func (a *jweAuthenticator) GenerateToken(username string, password string) (string, error) {
	user, err := a.userRepository.FindByUsername(username)
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
		Expiry:  jwt.NewNumericDate(time.Now().Add(AuthTimeout)),
	}

	user.Password = ""

	jwe, err := jwt.SignedAndEncrypted(signer, encrypter).Claims(cl).Claims(user).CompactSerialize()
	if err != nil {
		return "", err
	}

	a.userLogged = user

	return jwe, nil
}

func (a *jweAuthenticator) VerifyToken(jwe string) (*models.User, error) {
	tok, err := jwt.ParseSignedAndEncrypted(jwe)
	if err != nil {
		return nil, err
	}

	jws, err := tok.Decrypt(a.privateKey)
	if err != nil {
		return nil, err
	}

	cl := jwt.Claims{}
	u := models.User{}

	if err := jws.Claims(a.publicKey, &cl, &u); err != nil {
		return nil, err
	}

	err = cl.Validate(jwt.Expected{
		Time: time.Now(),
	})

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (a *jweAuthenticator) AuthMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "authorization heard is required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			http.Error(w, "invalid bearer header authorization", http.StatusBadRequest)
			return
		}

		token := parts[1]

		user, err := a.VerifyToken(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// TODO : Change to session variables
		r.Header.Set("X-User-Data", user.String())

		if user.Role == models.SuperAdminRole {
			h.ServeHTTP(w, r)
			return
		}

		switch r.Method {
		case "POST":
			if user.Role == models.AdminRole {
				h.ServeHTTP(w, r)
				return
			}
		case "PUT", "DELETE":
			if user.Role == models.AdminRole {
				h.ServeHTTP(w, r)
				return
			}
		case "GET", "OPTIONS":
			if user.Role == models.ReadOnlyRole || user.Role == models.AdminRole {
				h.ServeHTTP(w, r)
				return
			}
		}

		http.Error(w, "access denied", http.StatusForbidden)
	}
}

func NewJWEAuthenticator(logger *logrus.Entry, userRepository repositories.UsersRepository) *jweAuthenticator {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("create crypto keys: %v", err)
	}

	publicKey := &privateKey.PublicKey

	return &jweAuthenticator{privateKey: privateKey, publicKey: publicKey, logger: logger, userRepository: userRepository}
}

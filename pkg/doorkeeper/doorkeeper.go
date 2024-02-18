// "Boring cryptography" refers to cryptography designs and
// implementations that are obviously secure. This means having
// at least 2128 bits of security (Ed25519) instead of 1024-bit
// RSA (which is estimated to be approximately 280). Boring
// cryptography means being obviously constant-time. When cryptography
// is boring, there's far less room for implementers to make
// cataclysmic mistakes (such as repeating an ECDSA nonce).

// Remember, "Attacks only get better; they never get worse."

package doorkeeper

import (
	"crypto/sha512"
	"hash"
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rmscoal/moilerplate/internal/utils"
)

type Doorkeeper struct {
	JWT       jwtDoorkeeper
	General   generalDoorkeeper
	Encryptor encryptorDoorkeeper
}

type jwtDoorkeeper struct {
	signMethod      jwt.SigningMethod
	privateKey      any
	publicKey       any
	accessDuration  time.Duration
	refreshDuration time.Duration
	issuer          string
}

type generalDoorkeeper struct {
	hasherFunc                     func() hash.Hash // This will be used as the hashing method (may on top of PBKD2F)
	hashKeyLen                     int              // Special case for PBKD2F key length
	hashIter                       int              // Special case for PBKD2F iterator
	disableRandomGeneratedPassword bool             // Whether to randomly generate password or disable it and return constant string
	defaultGeneratedPassword       string           // The default value if random generated password is disable
}

type encryptorDoorkeeper struct {
	secretKey string
}

var (
	HMAC_SIGN_METHOD_TYPE   = reflect.TypeOf(&jwt.SigningMethodHMAC{})
	RSA_SIGN_METHOD_TYPE    = reflect.TypeOf(&jwt.SigningMethodRSA{})
	RSAPSS_SIGN_METHOD_TYPE = reflect.TypeOf(&jwt.SigningMethodRSAPSS{})
	ECDSA_SIGN_METHOD_TYPE  = reflect.TypeOf(&jwt.SigningMethodECDSA{})
	EdDSA_SIGN_METHOD_TYPE  = reflect.TypeOf(&jwt.SigningMethodEd25519{})
)

var (
	_defaultGeneratedPassword = "verystrongpassword"
	_defaultHasherFunc        = sha512.New384
	_defaultHashKeyLen        = 64
	_defaultHashIter          = 4096
	_defaultSigningMethod     = jwt.SigningMethodHS384
	_defaultAccessDuration    = 15 * time.Minute
	_defaultRefreshDuration   = 1 * time.Hour
)

var (
	once                     sync.Once
	doorkeeperSingleInstance *Doorkeeper
)

func GetDoorkeeper(opts ...Option) *Doorkeeper {
	if doorkeeperSingleInstance == nil {
		once.Do(func() {
			doorkeeperSingleInstance = &Doorkeeper{
				JWT: jwtDoorkeeper{
					accessDuration:  _defaultAccessDuration,
					refreshDuration: _defaultRefreshDuration,
					signMethod:      _defaultSigningMethod,
				},

				General: generalDoorkeeper{
					hasherFunc:               _defaultHasherFunc,
					hashKeyLen:               _defaultHashKeyLen,
					hashIter:                 _defaultHashIter,
					defaultGeneratedPassword: _defaultGeneratedPassword,
				},
			}

			for _, opt := range opts {
				opt(doorkeeperSingleInstance)
			}

			doorkeeperSingleInstance.loadJWTSecretKeys()
		})
	}

	return doorkeeperSingleInstance
}

// --- General ---

func (d *Doorkeeper) GetHasherFunc() func() hash.Hash {
	return d.General.hasherFunc
}

func (d *Doorkeeper) GetHashKeyLen() int {
	return d.General.hashKeyLen
}

func (d *Doorkeeper) GetHashIter() int {
	return d.General.hashIter
}

// --- JWT ---

func (d *Doorkeeper) GetIssuer() string {
	return d.JWT.issuer
}

func (d *Doorkeeper) GetSignMethod() jwt.SigningMethod {
	return d.JWT.signMethod
}

func (d *Doorkeeper) GetPubKey() interface{} {
	return d.JWT.publicKey
}

func (d *Doorkeeper) GetPrivKey() interface{} {
	return d.JWT.privateKey
}

func (d *Doorkeeper) GetJWTAccessDuration() time.Duration {
	return d.JWT.accessDuration
}

func (d *Doorkeeper) GetJWTRefreshDuration() time.Duration {
	return d.JWT.refreshDuration
}

func (d *Doorkeeper) GetConcreteSignMethod() reflect.Type {
	return reflect.TypeOf(d.JWT.signMethod)
}

func (d *Doorkeeper) loadJWTSecretKeys() {
	// If private key in still as string, convert it to []byte
	switch d.JWT.privateKey.(type) {
	case string:
		d.JWT.privateKey = utils.ConvertStringToByteSlice(d.JWT.privateKey.(string))
	}

	// If public key in still as string, convert it to []byte
	switch d.JWT.publicKey.(type) {
	case string:
		d.JWT.publicKey = utils.ConvertStringToByteSlice(d.JWT.publicKey.(string))
	}

	var err error
	switch d.GetConcreteSignMethod() {
	case RSA_SIGN_METHOD_TYPE, RSAPSS_SIGN_METHOD_TYPE:
		d.JWT.privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(d.JWT.privateKey.([]byte))
		if err != nil {
			log.Fatalf("unable to parse RSA private key from given jwt private key: %s", err)
		}

		d.JWT.publicKey, err = jwt.ParseRSAPublicKeyFromPEM(d.JWT.publicKey.([]byte))
		if err != nil {
			log.Fatalf("unable to parse RSA public key from given jwt public key: %s", err)
		}
	case ECDSA_SIGN_METHOD_TYPE:
		d.JWT.privateKey, err = jwt.ParseECPrivateKeyFromPEM(d.JWT.privateKey.([]byte))
		if err != nil {
			log.Fatalf("unable to parse EC private key from given jwt private key: %s", err)
		}

		d.JWT.publicKey, err = jwt.ParseECPublicKeyFromPEM(d.JWT.publicKey.([]byte))
		if err != nil {
			log.Fatalf("unable to parse EC public key from given jwt public key: %s", err)
		}
	case EdDSA_SIGN_METHOD_TYPE:
		d.JWT.privateKey, err = jwt.ParseEdPrivateKeyFromPEM(d.JWT.privateKey.([]byte))
		if err != nil {
			log.Fatalf("unable to parse EdDSA private key from given jwt private key: %s", err)
		}

		d.JWT.publicKey, err = jwt.ParseEdPublicKeyFromPEM(d.JWT.publicKey.([]byte))
		if err != nil {
			log.Fatalf("unable to parse EdDSA public key from given jwt public key: %s", err)
		}
	}
}

// --- Encryptor ---

func (d *Doorkeeper) GetEncryptorSecretKey() string {
	return d.Encryptor.secretKey
}

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
	"crypto"
	"crypto/sha512"
	"hash"
	"log"
	"os"
	"reflect"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Doorkeeper struct {
	// --- JWT ---
	signMethod      jwt.SigningMethod // This will be used for JWS / JWE
	certPath        string            // Stores the path to the certification keys
	privKey         interface{}       // Stores the private keys parsed from PEM (if asymmetric)
	pubKey          interface{}       // Stores the public keys parsed from PEM (if symmetric)
	issuer          string            // *Claims*
	AccessDuration  time.Duration     // Duration of an access token
	RefreshDuration time.Duration     // Duration of a refresh token

	// --- Password Hasher ---
	hasherFunc func() hash.Hash // This will be used as the hashing method (may on top of PBKD2F)
	hashKeyLen int              // Special case for PBKD2F key length
	hashIter   int              // Special case for PBKD2F iterator
}

var (
	HMAC_SIGN_METHOD_TYPE   = reflect.TypeOf(&jwt.SigningMethodHMAC{})
	RSA_SIGN_METHOD_TYPE    = reflect.TypeOf(&jwt.SigningMethodRSA{})
	RSAPSS_SIGN_METHOD_TYPE = reflect.TypeOf(&jwt.SigningMethodRSAPSS{})
	ECDSA_SIGN_METHOD_TYPE  = reflect.TypeOf(&jwt.SigningMethodECDSA{})
	EdDSA_SIGN_METHOD_TYPE  = reflect.TypeOf(&jwt.SigningMethodEd25519{})
)

var (
	_defaultHasherFunc      = sha512.New384
	_defaultHashKeyLen      = 64
	_defaultHashIter        = 4096
	_defaultSigningMethod   = jwt.SigningMethodHS384
	_defaultAccessDuration  = 15 * time.Minute
	_defaultRefreshDuration = 1 * time.Hour
)

var (
	once                     sync.Once
	doorkeeperSingleInstance *Doorkeeper
)

func GetDoorkeeper(opts ...Option) *Doorkeeper {
	if doorkeeperSingleInstance == nil {
		once.Do(func() {
			doorkeeperSingleInstance = &Doorkeeper{
				AccessDuration:  _defaultAccessDuration,
				RefreshDuration: _defaultRefreshDuration,
				hasherFunc:      _defaultHasherFunc,
				hashKeyLen:      _defaultHashKeyLen,
				hashIter:        _defaultHashIter,
				signMethod:      _defaultSigningMethod,
			}

			for _, opt := range opts {
				opt(doorkeeperSingleInstance)
			}

			doorkeeperSingleInstance.loadSecretKeys()
		})
	}

	return doorkeeperSingleInstance
}

func (d *Doorkeeper) GetHasherFunc() func() hash.Hash {
	return d.hasherFunc
}

func (d *Doorkeeper) GetHashKeyLen() int {
	return d.hashKeyLen
}

func (d *Doorkeeper) GetHashIter() int {
	return d.hashIter
}

func (d *Doorkeeper) GetIssuer() string {
	return d.issuer
}

func (d *Doorkeeper) GetSignMethod() jwt.SigningMethod {
	return d.signMethod
}

func (d *Doorkeeper) GetPubKey() interface{} {
	return d.pubKey
}

func (d *Doorkeeper) GetPrivKey() interface{} {
	return d.privKey
}

func (d *Doorkeeper) GetConcreteSignMethod() reflect.Type {
	return reflect.TypeOf(d.signMethod)
}

func (d *Doorkeeper) loadSecretKeys() {
	switch d.GetConcreteSignMethod() {
	case HMAC_SIGN_METHOD_TYPE:
		d.privKey = d.getSymmetricKeyFromFile("id_secret")
		d.pubKey = d.privKey
	case RSA_SIGN_METHOD_TYPE:
		privKeyByte, pubKeyByte := d.getAsymmetricKeysFromFile("id_rsa")
		d.privKey, d.pubKey = d.parseRSAKeysFromPem(privKeyByte, pubKeyByte)
	case RSAPSS_SIGN_METHOD_TYPE:
		privKeyByte, pubKeyByte := d.getAsymmetricKeysFromFile("id_rsa")
		d.privKey, d.pubKey = d.parseRSAKeysFromPem(privKeyByte, pubKeyByte)
	case ECDSA_SIGN_METHOD_TYPE:
		privKeyByte, pubKeyByte := d.getAsymmetricKeysFromFile("id_ecdsa")
		d.privKey, d.pubKey = d.parseECKeysFromPem(privKeyByte, pubKeyByte)
	case EdDSA_SIGN_METHOD_TYPE:
		privKeyByte, pubKeyByte := d.getAsymmetricKeysFromFile("id_ed2559")
		d.privKey, d.pubKey = d.parseEdKeysFromPem(privKeyByte, pubKeyByte)
	}
}

func (d *Doorkeeper) getSymmetricKeyFromFile(filename string) []byte {
	key, err := os.ReadFile(d.certPath + "/" + filename)
	if err != nil {
		log.Fatalln(err)
	}

	return key
}

func (d *Doorkeeper) getAsymmetricKeysFromFile(filename string) ([]byte, []byte) {
	privKey, err := os.ReadFile(d.certPath + "/" + filename)
	if err != nil {
		log.Fatalln(err)
	}
	pubKey, err := os.ReadFile(d.certPath + "/" + filename + ".pub")
	if err != nil {
		log.Fatalln(err)
	}

	return privKey, pubKey
}

func (d *Doorkeeper) parseECKeysFromPem(privByte, pubByte []byte) (crypto.PrivateKey, crypto.PublicKey) {
	privKey, err := jwt.ParseECPrivateKeyFromPEM(privByte)
	if err != nil {
		log.Fatalf("unable to parse ec private key: %s", err)
	}
	pubKey, err := jwt.ParseECPublicKeyFromPEM(pubByte)
	if err != nil {
		log.Fatalf("unable to parse ec public key: %s", err)
	}

	return privKey, pubKey
}

func (d *Doorkeeper) parseRSAKeysFromPem(privByte, pubByte []byte) (crypto.PrivateKey, crypto.PublicKey) {
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(privByte)
	if err != nil {
		log.Fatalf("unable to parse rsa private key: %s", err)
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubByte)
	if err != nil {
		log.Fatalf("unable to parse rsa public key: %s", err)
	}

	return privKey, pubKey
}

func (d *Doorkeeper) parseEdKeysFromPem(privByte, pubByte []byte) (crypto.PrivateKey, crypto.PublicKey) {
	privKey, err := jwt.ParseEdPrivateKeyFromPEM(privByte)
	if err != nil {
		log.Fatalf("unable to parse ed private key: %s", err)
	}
	pubKey, err := jwt.ParseEdPublicKeyFromPEM(pubByte)
	if err != nil {
		log.Fatalf("unable to parse ed public key: %s", err)
	}

	return privKey, pubKey
}

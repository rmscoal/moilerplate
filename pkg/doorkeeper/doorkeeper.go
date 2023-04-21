// "Boring cryptography" refers to cryptography designs and
// implementations that are obviously secure. This means having
// at least 2128 bits of security (Ed25519) instead of 1024-bit
// RSA (which is estimated to be approximately 280). Boring
// cryptography means being obviously constant-time. When cryptography
// is boring, there's far less room for implementers to make
// cataclysmic mistakes (such as repeating an ECDSA nonce).

// Cryptographers are working hard to bring boring cryptography to the masses. Paragon Initiative Enterprises is similarly working hard to bring boring levels of security to PHP. That is why we're building Airship: The PHP community deserves a CMS/blogging platform that is obviously secure, written from an understanding of how PHP applications are attacked in the real world.

// Remember, "Attacks only get better; they never get worse."

package doorkeeper

import (
	"crypto"
	"io/ioutil"
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Doorkeeper struct {
	signMethod jwt.SigningMethod
	hashMethod crypto.Hash

	path          string
	signMethodStr string

	issuer  string
	secret  string
	salt    string
	privKey any
	pubKey  any

	Duration time.Duration
}

var (
	_defaultHashMethod    = crypto.SHA256
	_defaultSigningMethod = jwt.SigningMethodHS384
	_defaultSecretKey     = "secretKey" // this value should always be replace by passing options
	_defaultSalt          = "saltKey"   // this value should always be replace by passing options
	_defaultDuration      = 15 * time.Minute
)

var (
	once                     sync.Once
	doorkeeperSingleInstance *Doorkeeper
)

func GetDoorkeeper(opts ...Option) *Doorkeeper {
	if doorkeeperSingleInstance == nil {
		once.Do(func() {
			doorkeeperSingleInstance = &Doorkeeper{
				Duration:   _defaultDuration,
				hashMethod: _defaultHashMethod,
				signMethod: _defaultSigningMethod,
				secret:     _defaultSecretKey,
				salt:       _defaultSalt,
			}

			for _, opt := range opts {
				opt(doorkeeperSingleInstance)
			}

			doorkeeperSingleInstance.loadSecretKeys()
		})
	}

	return doorkeeperSingleInstance
}

func (d *Doorkeeper) GetIssuer() string {
	return d.issuer
}

func (d *Doorkeeper) GetSalt() string {
	return d.salt
}

func (d *Doorkeeper) GetSignMethod() jwt.SigningMethod {
	return d.signMethod
}

func (d *Doorkeeper) GetHasMethod() crypto.Hash {
	return d.hashMethod
}

func (d *Doorkeeper) GetPubKey() any {
	return d.pubKey
}

func (d *Doorkeeper) GetPrivKey() any {
	return d.privKey
}

func (d *Doorkeeper) GetConcreteSignMethod() reflect.Type {
	return reflect.TypeOf(d.signMethod)
}

func (d *Doorkeeper) loadSecretKeys() {
	switch d.signMethodStr {
	case "HMAC":
		d.privKey, d.pubKey = []byte(d.secret), []byte(d.secret)
	case "RSA":
		privKeyByte, pubKeyByte := d.getKeyFromFile("id_rsa")
		privKey, err := jwt.ParseRSAPrivateKeyFromPEM(privKeyByte)
		if err != nil {
			log.Fatalf("unable to parse rsa private key: %s", err)
		}
		pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyByte)
		if err != nil {
			log.Fatalf("unable to parse rsa private key: %s", err)
		}
		d.privKey = privKey
		d.pubKey = pubKey
	case "ECDSA":
		privKeyByte, pubKeyByte := d.getKeyFromFile("id_ecdsa")
		privKey, err := jwt.ParseECPrivateKeyFromPEM(privKeyByte)
		if err != nil {
			log.Fatalf("unable to parse ec private key: %s", err)
		}
		pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyByte)
		if err != nil {
			log.Fatalf("unable to parse ec private key: %s", err)
		}
		d.privKey = privKey
		d.pubKey = pubKey
	case "EdDSA":
		privKeyByte, pubKeyByte := d.getKeyFromFile("id_ed2559")
		privKey, err := jwt.ParseEdPrivateKeyFromPEM(privKeyByte)
		if err != nil {
			log.Fatalf("unable to parse ed private key: %s", err)
		}
		pubKey, err := jwt.ParseEdPublicKeyFromPEM(pubKeyByte)
		if err != nil {
			log.Fatalf("unable to parse ed private key: %s", err)
		}
		d.privKey = privKey
		d.pubKey = pubKey
	}
}

func (d *Doorkeeper) getKeyFromFile(fileName string) ([]byte, []byte) {
	privKey, err := ioutil.ReadFile(d.path + "/" + fileName)
	if err != nil {
		log.Fatalln(err)
	}
	pubKey, err := ioutil.ReadFile(d.path + "/" + fileName + ".pub")
	if err != nil {
		log.Fatalln(err)
	}

	return privKey, pubKey
}

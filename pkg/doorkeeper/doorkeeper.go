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
	"crypto/rand"
	"crypto/rsa"
	"log"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Doorkeeper struct {
	signMethod jwt.SigningMethod
	hashMethod crypto.Hash

	signMethodStr string

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
	_defaultDuration      = 300 * time.Second
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

			doorkeeperSingleInstance.generateKeys()
		})
	}

	return doorkeeperSingleInstance
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

func (d *Doorkeeper) generateKeys() {
	switch d.signMethodStr {
	case "HMAC":
		d.privKey, d.pubKey = d.secret, d.secret
	case "RSA":
		privKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			log.Fatalf("doorkeeper unable to generate private key")
		}
		d.privKey = privKey
		d.pubKey = privKey.PublicKey
	}
}

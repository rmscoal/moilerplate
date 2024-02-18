package doorkeeper

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rmscoal/moilerplate/internal/utils"
	"golang.org/x/crypto/sha3"
)

// Option .-.
type Option func(*Doorkeeper)

// --- JWT ---

func RegisterJWTAccessDuration(t time.Duration) Option {
	return func(d *Doorkeeper) {
		if t > 0 {
			d.JWT.accessDuration = t
		}
	}
}

func RegisterJWTRefreshDuration(t time.Duration) Option {
	return func(d *Doorkeeper) {
		if t > 0 {
			d.JWT.refreshDuration = t
		}
	}
}

func RegisterJWTIssuer(iss string) Option {
	return func(d *Doorkeeper) {
		if iss != "" {
			d.JWT.issuer = iss
		}
	}
}

func RegisterJWTPrivateKey(priv string) Option {
	return func(d *Doorkeeper) {
		d.JWT.privateKey = utils.ConvertStringToByteSlice(priv)
	}
}

func RegisterJWTPublicKey(pub string) Option {
	return func(d *Doorkeeper) {
		d.JWT.publicKey = utils.ConvertStringToByteSlice(pub)
	}
}

func RegisterJWTSignMethod(alg, size string) Option {
	return func(d *Doorkeeper) {
		switch alg {
		case "HMAC":
			switch size {
			case "256":
				d.JWT.signMethod = jwt.SigningMethodHS256
			case "384":
				d.JWT.signMethod = jwt.SigningMethodHS384
			case "512":
				d.JWT.signMethod = jwt.SigningMethodHS512
			}
		case "RSA":
			switch size {
			case "256":
				d.JWT.signMethod = jwt.SigningMethodRS256
			case "384":
				d.JWT.signMethod = jwt.SigningMethodRS384
			case "512":
				d.JWT.signMethod = jwt.SigningMethodRS512
			}
		case "ECDSA":
			switch size {
			case "256":
				d.JWT.signMethod = jwt.SigningMethodES256
			case "384":
				d.JWT.signMethod = jwt.SigningMethodES384
			case "512":
				d.JWT.signMethod = jwt.SigningMethodES512
			}
		case "RSA-PSS":
			switch size {
			case "256":
				d.JWT.signMethod = jwt.SigningMethodPS256
			case "384":
				d.JWT.signMethod = jwt.SigningMethodPS384
			case "512":
				d.JWT.signMethod = jwt.SigningMethodPS512
			}
		case "EdDSA":
			d.JWT.signMethod = &jwt.SigningMethodEd25519{}
		}
	}
}

// --- General ---

func RegisterGeneralHasherFunc(alg string) Option {
	return func(d *Doorkeeper) {
		switch alg {
		case "SHA1":
			d.General.hasherFunc = sha1.New
		case "SHA224":
			d.General.hasherFunc = sha256.New224
		case "SHA256":
			d.General.hasherFunc = sha256.New
		case "SHA384":
			d.General.hasherFunc = sha512.New384
		case "SHA512":
			d.General.hasherFunc = sha512.New
		case "SHA3_224":
			d.General.hasherFunc = sha3.New224
		case "SHA3_256":
			d.General.hasherFunc = sha3.New256
		case "SHA3_384":
			d.General.hasherFunc = sha3.New384
		case "SHA3_512":
			d.General.hasherFunc = sha3.New512
		}
	}
}

func ShouldGenerateRandomPassword(b bool) Option {
	return func(d *Doorkeeper) {
		d.General.disableRandomGeneratedPassword = b
	}
}

func RegisterDefaultGeneratedPassword(s string) Option {
	return func(d *Doorkeeper) {
		d.General.defaultGeneratedPassword = s
	}
}

// --- Encryptor ---

func RegisterEncryptorSecretKey(secret string) Option {
	return func(d *Doorkeeper) {
		sha := sha256.Sum256(utils.ConvertStringToByteSlice(secret))
		d.Encryptor.secretKey = string(sha[:])
	}
}

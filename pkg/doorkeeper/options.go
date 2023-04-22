package doorkeeper

import (
	"crypto"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Option -.
type Option func(*Doorkeeper)

func RegisterSecretKey(key string) Option {
	return func(d *Doorkeeper) {
		if key != "" {
			d.secret = key
		}
	}
}

func RegisterSalt(salt string) Option {
	return func(d *Doorkeeper) {
		if salt != "" {
			d.salt = salt
		}
	}
}

func RegisterPath(path string) Option {
	return func(d *Doorkeeper) {
		d.path = path
	}
}

func RegisterDuration(t time.Duration) Option {
	return func(d *Doorkeeper) {
		if t > d.Duration {
			d.Duration = t
		}
	}
}

func RegisterIssuer(iss string) Option {
	return func(d *Doorkeeper) {
		d.issuer = iss
	}
}

func RegisterHashMethod(alg string) Option {
	return func(d *Doorkeeper) {
		switch alg {
		case "MD4":
			d.hashMethod = crypto.MD4
		case "MD5":
			d.hashMethod = crypto.MD5
		case "SHA1":
			d.hashMethod = crypto.SHA1
		case "SHA224":
			d.hashMethod = crypto.SHA224
		case "SHA256":
			d.hashMethod = crypto.SHA256
		case "SHA384":
			d.hashMethod = crypto.SHA384
		case "SHA512":
			d.hashMethod = crypto.SHA512
		case "SHA3_224":
			d.hashMethod = crypto.SHA3_224
		case "SHA3_256":
			d.hashMethod = crypto.SHA3_256
		case "SHA3_384":
			d.hashMethod = crypto.SHA3_384
		case "SHA3_512":
			d.hashMethod = crypto.SHA3_512
		}
	}
}

func RegisterSignMethod(alg, size string) Option {
	return func(d *Doorkeeper) {
		switch alg {
		case "HMAC":
			switch size {
			case "256":
				d.signMethod = jwt.SigningMethodHS256
			case "384":
				d.signMethod = jwt.SigningMethodHS384
			case "512":
				d.signMethod = jwt.SigningMethodHS512
			}
		case "RSA":
			switch size {
			case "256":
				d.signMethod = jwt.SigningMethodRS256
			case "384":
				d.signMethod = jwt.SigningMethodRS384
			case "512":
				d.signMethod = jwt.SigningMethodRS512
			}
		case "ECDSA":
			switch size {
			case "256":
				d.signMethod = jwt.SigningMethodES256
			case "384":
				d.signMethod = jwt.SigningMethodES384
			case "512":
				d.signMethod = jwt.SigningMethodES512
			}
		case "RSA-PSS":
			switch size {
			case "256":
				d.signMethod = jwt.SigningMethodPS256
			case "384":
				d.signMethod = jwt.SigningMethodPS384
			case "512":
				d.signMethod = jwt.SigningMethodPS512
			}
		case "EdDSA":
			d.signMethod = &jwt.SigningMethodEd25519{}
		}
	}
}

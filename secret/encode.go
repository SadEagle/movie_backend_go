package secret

import (
	"bytes"
	"golang.org/x/crypto/argon2"
)

// Random 16 bytes
var ARGON2_SALT = []byte{133, 176, 19, 239, 82, 8, 236, 171, 119, 146, 145, 152, 166, 168, 27, 240}

func EncodePassword(password string) []byte {
	return argon2.IDKey([]byte(password), ARGON2_SALT, 3, 32*1024, 4, 32)
}

func IsEqualPasswd(hashed_passwd []byte, passwd []byte) bool {
	return bytes.Equal(hashed_passwd, passwd)
}

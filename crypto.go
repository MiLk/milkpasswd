package milkpasswd

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha256"
	"errors"
	"log"
)

// Inspired by https://www.socketloop.com/tutorials/golang-how-to-encrypt-with-aes-crypto

// HASH

func Md5sum(data []byte) []byte {
	sum := md5.Sum(data)
	out := make([]byte, len(sum))
	copy(out, sum[:])
	return out
}

func Sha256sum(data []byte) []byte {
	sum := sha256.Sum256(data)
	out := make([]byte, len(sum))
	copy(out, sum[:])
	return out
}

// ENCRYPTION

func generateBlock(key []byte) (block cipher.Block, err error) {
	switch len(key) {
	case 16:
		log.Printf("WARNING: Using AES-128 to encrypt data!")
	case 24:
		log.Printf("WARNING: Using AES-192 to encrypt data!")
	case 32:
		break
	default:
		err = errors.New("milkpasswd crypto: invalid key size")
		return
	}

	block, err = aes.NewCipher(key)
	return
}

func generateInitializationVector(ciphertext []byte) (iv []byte, err error) {
	if len(ciphertext) < aes.BlockSize {
		err = errors.New("milkpasswd crypto: invalid length for the ciphertext")
		return
	}

	if len(ciphertext) > aes.BlockSize {
		log.Printf("WARNING: The ciphertext is too long! The end will be ignored!")
	}

	iv = ciphertext[:aes.BlockSize]
	return
}

func Encrypt(key []byte, ciphertext []byte, str []byte) (encrypted []byte, err error) {
	block, err := generateBlock(key)
	if err != nil {
		return
	}

	iv, err := generateInitializationVector(ciphertext)
	if err != nil {
		return
	}

	encrypter := cipher.NewCFBEncrypter(block, iv)

	encrypted = make([]byte, len(str))
	encrypter.XORKeyStream(encrypted, []byte(str))

	return
}

func Decrypt(key []byte, ciphertext []byte, encrypted []byte) (str string, err error) {
	block, err := generateBlock(key)
	if err != nil {
		return
	}

	iv, err := generateInitializationVector(ciphertext)
	if err != nil {
		return
	}

	decrypter := cipher.NewCFBDecrypter(block, iv)

	decrypted := make([]byte, len(encrypted))
	decrypter.XORKeyStream(decrypted, encrypted)

	str = string(decrypted)

	return
}

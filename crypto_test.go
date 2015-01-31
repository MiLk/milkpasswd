package milkpasswd

import (
	"encoding/hex"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMd5sum(t *testing.T) {
	Convey("Should return the md5sum", t, func() {
		input := []byte("azerty")
		sum, _ := hex.DecodeString("ab4f63f9ac65152575886860dde480a1")
		result := Md5sum(input)
		So(result, ShouldHaveSameTypeAs, sum)
		So(result, ShouldResemble, sum)
	})
}

func TestSha256sum(t *testing.T) {
	Convey("Should return the sha256sum", t, func() {
		input := []byte("azerty")
		sum, _ := hex.DecodeString("f2d81a260dea8a100dd517984e53c56a7523d96942a834b9cdc249bd4e8c7aa9")
		result := Sha256sum(input)
		So(result, ShouldHaveSameTypeAs, sum)
		So(result, ShouldResemble, sum)
	})
}

func TestEncrypt(t *testing.T) {
	Convey("Should encrypt a given string", t, func() {
		key := Sha256sum([]byte("masterkey"))
		cipher := Md5sum([]byte("ciphertext"))
		expected, _ := hex.DecodeString("040928641d81399ab3033c")
		input := []byte("secret data")
		encrypted, err := Encrypt(key, cipher, input)
		So(err, ShouldEqual, nil)
		So(encrypted, ShouldHaveSameTypeAs, expected)
		So(encrypted, ShouldResemble, expected)
	})
}

func TestDecrypt(t *testing.T) {
	Convey("Should encrypt a given string", t, func() {
		key := Sha256sum([]byte("masterkey"))
		cipher := Md5sum([]byte("ciphertext"))
		expected := "secret data"
		input, _ := hex.DecodeString("040928641d81399ab3033c")
		decrypted, err := Decrypt(key, cipher, input)
		So(err, ShouldEqual, nil)
		So(decrypted, ShouldHaveSameTypeAs, expected)
		So(decrypted, ShouldResemble, expected)
	})
}

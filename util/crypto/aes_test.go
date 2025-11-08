package crypto

import "testing"

func TestAesCBC(t *testing.T) {
	plainText := "helloworld"
	key := "this is a secret"
	cryptedByte, err := AESEncryptCBC([]byte(plainText), []byte(key))
	if err != nil {
		t.Errorf("AES encrypt error: %s", err.Error())
	}
	deCrypted, err := AESDecryptCBC(cryptedByte, []byte(key))
	if err != nil {
		t.Errorf("AES decrypt error: %s", err.Error())
	}
	if string(deCrypted) != string(plainText) {
		t.Errorf("the AES decrypt result does not match the plainText, plainText: %s, deCrypted: %s \n", plainText, deCrypted)
	}
}

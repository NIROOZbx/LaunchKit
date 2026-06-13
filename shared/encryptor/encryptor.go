package encryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"github.com/Launchkit-org/LaunchKit/shared/serializer"
)

func EncryptMap(data map[string]string, secretKey string) (string, error) {
	bytes, err := serializer.Marshal(data)
	if err != nil {
		return "", err
	}
	return Encrypt(bytes, secretKey)
}

func DecryptToMap(encrypted string, secretKey string) (map[string]string, error) {
	bytes, err := Decrypt(encrypted, secretKey)
	if err != nil {
		return nil, err
	}
	var data map[string]string
	if err := serializer.UnMarshal(bytes, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func Encrypt(text []byte,key string)(string,error){

	secretkey:=[]byte(key)
	if len(secretkey)!=32{
		return "",errors.New("encryption key must be 32 bytes longg")
	}
	block,err:=aes.NewCipher(secretkey)
	if err!=nil{
		return "",err
	}
	aesGM,err:=cipher.NewGCM(block)
	if err!=nil{
		return "",err
	}
	nonce:=make([]byte,aesGM.NonceSize())
	if _,err:=io.ReadFull(rand.Reader,nonce);err!=nil{
		return "",nil
	}
	cipherText:= aesGM.Seal(nonce,nonce,text,nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}


func Decrypt(cipherText string, Key string) ([]byte, error) {
	secretKey := []byte(Key)
	if len(secretKey) != 32 {
		return nil, errors.New("encryption key must be exactly 32 bytes long")
	}
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
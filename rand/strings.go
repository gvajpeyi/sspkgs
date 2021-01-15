package rand

import (
	"crypto/rand"
	"encoding/base64"
)

const RememberTokenBytes = 32

// Bytes will help generate n random bytes, or will
// return an error if there was one.  The uses the crypto/rand
// package so it is safe to use with things like remember tokens.
func Bytes(n int)([]byte, error){
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil{
		return nil, err
	}
	return b, nil
}
func String(nBytes int)(string, error){
	b, err := Bytes(nBytes)
	
	if err != nil{
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

//RememberToken is a helper that generates remember tokens
// of a predetermined bytes size.

func RememberToken()(string, error){
	return String(RememberTokenBytes)
	
}
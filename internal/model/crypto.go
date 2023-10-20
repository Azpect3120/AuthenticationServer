package model

import (
	"crypto/sha256"
	"encoding/hex"
)

// Hash a string using the sha256 hash method
func HashString (ch chan *StringResult, string string)  {
	hash := sha256.New()

	if _, err := hash.Write([]byte(string)); err != nil {
		// return "", err
		ch <- &StringResult{"", &Error{err.Error(), 500}}
		return
	}

	hashed := hex.EncodeToString(hash.Sum(nil))
	// return hashed, nil
	ch <- &StringResult{hashed, nil}
}

// Compare a string to a hashed string 
func CompareString (unhashed string, hashed string) (bool, *Error) {
	strCh := make(chan *StringResult)

	go HashString(strCh, unhashed)

	result := <- strCh

	if result.Error != nil {
		return false, result.Error
	}

	return result.String == hashed, nil
}

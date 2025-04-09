package utils

import (
	"crypto/rand"
)

func GenerateOTP(length int) string {

	otp := make([]byte, length)
	_, err := rand.Read(otp)

	if err != nil {
		return ""
	}
	for i := 0; i < length; i++ {
		otp[i] = otp[i]%10 + '0'
	}

	return string(otp)
}

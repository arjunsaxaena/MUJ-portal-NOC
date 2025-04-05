package util

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var OTPStore = struct {
	sync.RWMutex
	data map[string]string
}{data: make(map[string]string)}

var validatedEmails = struct {
	sync.RWMutex
	data map[string]bool
}{data: make(map[string]bool)}

func GenerateOTP(email string) string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	otp := rng.Intn(900000) + 100000

	OTPStore.Lock()
	OTPStore.data[email] = fmt.Sprintf("%d", otp)
	OTPStore.Unlock()

	go func(email string) {
		time.Sleep(5 * time.Minute)
		OTPStore.Lock()
		delete(OTPStore.data, email)
		OTPStore.Unlock()
	}(email)

	return fmt.Sprintf("%d", otp)
}

func ValidateOTP(email, otp string) bool {
	OTPStore.RLock()
	storedOTP, exists := OTPStore.data[email]
	OTPStore.RUnlock()

	return exists && storedOTP == otp
}

func MarkEmailValidated(email string) {
	validatedEmails.Lock()
	validatedEmails.data[email] = true
	validatedEmails.Unlock()
}

func IsEmailValidated(email string) bool {
	validatedEmails.RLock()
	_, exists := validatedEmails.data[email]
	validatedEmails.RUnlock()
	return exists
}

func ClearEmailValidation(email string) {
	validatedEmails.Lock()
	delete(validatedEmails.data, email)
	validatedEmails.Unlock()
}

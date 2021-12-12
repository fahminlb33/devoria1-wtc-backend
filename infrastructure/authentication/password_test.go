package authentication_test

import (
	"testing"

	"github.com/fahminlb33/devoria1-wtc-backend/infrastructure/authentication"
	"github.com/stretchr/testify/assert"
)

func TestVerifyPassword(t *testing.T) {
	result := authentication.VerifyPassword("fahmi", "$2a$14$FwEMlETO/XHoB90v/O9zK.KNfm.G5ZxJUcoZWS5IFBOM/Ao4adKiW")
	assert.True(t, result)
}

func TestSafeCompareString(t *testing.T) {
	result := authentication.SafeCompareString("fahmi", "fahmi")
	assert.True(t, result)
}

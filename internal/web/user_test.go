package web

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestEncrypt(t *testing.T) {
	pwd := "123456"
	encryptedPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}

	err = bcrypt.CompareHashAndPassword(encryptedPwd, []byte(pwd))
	// err == nil -> equal
	if err == nil {
		fmt.Println("相等的")
	}
	assert.NoError(t, err)

}

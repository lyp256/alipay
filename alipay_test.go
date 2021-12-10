package alipay

import (
	"crypto"
	"crypto/rsa"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSign(t *testing.T) {
	key, err := rsa.GenerateKey(rand.New(rand.NewSource(time.Now().Unix())), 4096)
	assert.NoError(t, err)
	data := make([]byte, 1<<10)
	signData, err := Sign(data, key, crypto.SHA256)
	assert.NoError(t, err)
	err = Verify(data, signData, &key.PublicKey, crypto.SHA256)
	assert.NoError(t, err)
}

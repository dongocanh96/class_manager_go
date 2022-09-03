package token

import (
	"crypto/rsa"
	"io/ioutil"
	"testing"
	"time"

	"github.com/dongocanh96/class_manager_go/util"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
)

func getKeyPairData(t *testing.T) (privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) {
	privateKeyByte, err := ioutil.ReadFile("../private.pem")
	if err != nil {
		panic(err)
	}

	privatekey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyByte)
	if err != nil {
		panic(err)
	}

	publicKeyByte, err := ioutil.ReadFile("../public.pem")
	if err != nil {
		panic(err)
	}

	publickey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyByte)
	if err != nil {
		panic(err)
	}

	return privatekey, publickey
}

func TestJWTMaker(t *testing.T) {
	privateKey, publicKey := getKeyPairData(t)
	maker, err := NewJWTMaker(privateKey, publicKey)
	require.NoError(t, err)

	userid := util.RandomInt(1, 15)
	username := util.RandomString(8)
	isTeacher := util.RandomBoolean()
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(userid, username, isTeacher, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotEmpty(t, token)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, userid, payload.Userid)
	require.Equal(t, username, payload.Username)
	require.Equal(t, isTeacher, payload.IsTeacher)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	privateKey, publicKey := getKeyPairData(t)
	maker, err := NewJWTMaker(privateKey, publicKey)
	require.NoError(t, err)

	userid := util.RandomInt(1, 15)
	username := util.RandomString(8)
	isTeacher := util.RandomBoolean()
	duration := -time.Minute

	token, payload, err := maker.CreateToken(userid, username, isTeacher, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotEmpty(t, token)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	userid := util.RandomInt(1, 15)
	username := util.RandomString(8)
	isTeacher := util.RandomBoolean()
	duration := time.Minute

	payload, err := NewPayload(userid, username, isTeacher, duration)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	privateKey, publicKey := getKeyPairData(t)
	maker, err := NewJWTMaker(privateKey, publicKey)
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)

}

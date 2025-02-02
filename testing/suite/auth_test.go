package suite_test

import (
	"testing"
	"time"

	ssov1 "github.com/Pinkman-77/Protobuf/gen/go/sso"
	"github.com/Pinkman-77/records-restapi/testing/suite"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	appID             = 1
	appSecret         = "test-secret"
	PasswordDefLength = 10
)

func TestLogin_Success(t *testing.T) {
	ctx, st := suite.New(t)

	// Generate a test user
	email := gofakeit.Email()
	pass := randomPassword()

	// Register the user
	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: pass,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	// Login with the registered user
	respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: pass,
		AppId:    appID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, respLogin.GetToken())

	// Verify JWT token
	token := respLogin.GetToken()
	loginTime := time.Now()

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	require.True(t, ok)

	assert.Equal(t, respReg.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, int(claims["app_id"].(float64)))

	// Check if exp of token is in correct range
	const deltaSeconds = 1
	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTll).Unix(), claims["exp"].(float64), deltaSeconds)
}

func TestRegister_DuplicateEmail(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := randomPassword()

	// Register the first user
	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: pass,
	})
	require.NoError(t, err)
	require.NotEmpty(t, respReg.GetUserId())

	// Try registering the same email again
	respReg, err = st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: pass,
	})
	require.Error(t, err)
	assert.Empty(t, respReg.GetUserId())
	assert.ErrorContains(t, err, "user already exists")
}

func TestRegister_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		email       string
		password    string
		expectedErr string
	}{
		{"Empty Password", gofakeit.Email(), "", "password is required"},
		{"Empty Email", "", randomPassword(), "email is required"},
		{"Both Empty", "", "", "email is required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Email:    tt.email,
				Password: tt.password,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestLogin_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		email       string
		password    string
		appID       int32
		expectedErr string
	}{
		{"Empty Password", gofakeit.Email(), "", appID, "password is required"},
		{"Empty Email", "", randomPassword(), appID, "email is required"},
		{"Both Empty", "", "", appID, "email is required"},
		{"Invalid Password", gofakeit.Email(), randomPassword(), appID, "invalid email or password"},
		{"Missing AppID", gofakeit.Email(), randomPassword(), 0, "app_id is required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Email:    gofakeit.Email(),
				Password: randomPassword(),
			})
			require.NoError(t, err)

			_, err = st.AuthClient.Login(ctx, &ssov1.LoginRequest{
				Email:    tt.email,
				Password: tt.password,
				AppId:    int64(tt.appID),
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

// Helper function to generate secure random passwords
func randomPassword() string {
	return gofakeit.Password(true, true, true, true, false, PasswordDefLength)
}

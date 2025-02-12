package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/davidalecrim1/wolf-workouts/internal/users/adapter"
	"github.com/davidalecrim1/wolf-workouts/internal/users/app"
	"github.com/davidalecrim1/wolf-workouts/internal/users/config"
	"github.com/davidalecrim1/wolf-workouts/internal/users/handler"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

var (
	db *sqlx.DB
	ts *httptest.Server
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, closeDB := GetTestDatabase(ctx)
	db = conn
	defer closeDB()

	repo := adapter.NewPostgresUserRepository(db)
	svc := app.NewUserService(repo, config.NewConfig(os.Getenv("USERS_JWT_SECRET")))
	handler := handler.NewUserHandler(svc)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	handler.RegisterRoutes(router)

	ts = httptest.NewServer(router)
	defer ts.Close()

	os.Exit(m.Run())
}

func TestE2E(t *testing.T) {
	t.Parallel()

	t.Run("should create a user", func(t *testing.T) {
		t.Parallel()

		reqBody := map[string]interface{}{
			"name":     "John Doe",
			"email":    "john.doe@example.com",
			"password": "password",
			"role":     "trainee",
		}

		_ = createUser(t, ts, reqBody)
	})

	t.Run("should login a user and return a valid token", func(t *testing.T) {
		t.Parallel()

		createUserReqBody := map[string]interface{}{
			"name":     "Jane Doe",
			"email":    "jane.doe@example.com",
			"password": "example-password",
			"role":     "trainee",
		}

		_ = createUser(t, ts, createUserReqBody)

		formData := url.Values{}
		formData.Add("email", "jane.doe@example.com")
		formData.Add("password", "example-password")

		res := loginUserRequest(t, ts, formData)
		require.Equal(t, http.StatusOK, res.StatusCode)

		body := res.Body

		var loginResponse map[string]interface{}
		err := json.NewDecoder(body).Decode(&loginResponse)
		require.NoError(t, err)
		require.NoError(t, res.Body.Close())
		require.NotEmpty(t, loginResponse["access_token"])

		assertTokenIsValid(t,
			loginResponse["access_token"].(string),
			createUserReqBody["name"].(string),
			createUserReqBody["role"].(string),
		)
	})
}

func createUser(t *testing.T, ts *httptest.Server, reqBody map[string]interface{}) *http.Response {
	res := createUserRequest(t, ts, reqBody)
	require.NoError(t, res.Body.Close())
	require.Equal(t, http.StatusCreated, res.StatusCode)

	return res
}

func createUserRequest(t *testing.T, ts *httptest.Server, reqBody map[string]interface{}) *http.Response {
	jsonBody, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, ts.URL+"/users", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	res, err := ts.Client().Do(req)
	require.NoError(t, err)

	return res
}

func loginUserRequest(t *testing.T, ts *httptest.Server, formData url.Values) *http.Response {
	req, err := http.NewRequest(
		http.MethodPost,
		ts.URL+"/users/login",
		bytes.NewBufferString(formData.Encode()),
	)
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := ts.Client().Do(req)
	require.NoError(t, err)

	return res
}

func assertTokenIsValid(t *testing.T, token string, name string, role string) {
	config := config.NewConfig(os.Getenv("USERS_JWT_SECRET"))

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", t.Header["alg"])
		}

		_, ok := t.Claims.(jwt.MapClaims)
		if !ok {
			return nil, fmt.Errorf("invalid token claims")
		}

		return config.GetJWTSecret(), nil
	})
	require.NoError(t, err)
	require.True(t, parsedToken.Valid)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	require.True(t, ok)

	userID, ok := claims["sub"].(string)
	require.True(t, ok)
	require.NotEmpty(t, userID)

	username, ok := claims["user_name"].(string)
	require.True(t, ok)
	require.NotEmpty(t, username)
	require.Equal(t, name, username)

	userRole, ok := claims["user_role"].(string)
	require.True(t, ok)
	require.NotEmpty(t, userRole)
	require.Equal(t, role, userRole)
}

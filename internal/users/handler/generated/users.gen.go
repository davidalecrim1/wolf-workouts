// Package generated provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.2.0 DO NOT EDIT.
package generated

const (
	BearerAuthScopes = "bearerAuth.Scopes"
)

// Defines values for CreateUserRequestRole.
const (
	Trainee CreateUserRequestRole = "trainee"
	Trainer CreateUserRequestRole = "trainer"
)

// CreateUserRequest defines model for CreateUserRequest.
type CreateUserRequest struct {
	Email    string                `json:"email"`
	Name     string                `json:"name"`
	Password string                `json:"password"`
	Role     CreateUserRequestRole `json:"role"`
}

// CreateUserRequestRole defines model for CreateUserRequest.Role.
type CreateUserRequestRole string

// LoginUserRequest defines model for LoginUserRequest.
type LoginUserRequest struct {
	Email    string `binding:"required" form:"email" json:"email"`
	Password string `binding:"required" form:"password" json:"password"`
}

// LoginUserResponse defines model for LoginUserResponse.
type LoginUserResponse struct {
	AccessToken *string `json:"access_token,omitempty"`
}

// ResponseError defines model for ResponseError.
type ResponseError struct {
	Message string `json:"message"`
}

// CreateUserJSONRequestBody defines body for CreateUser for application/json ContentType.
type CreateUserJSONRequestBody = CreateUserRequest

// LoginUserFormdataRequestBody defines body for LoginUser for application/x-www-form-urlencoded ContentType.
type LoginUserFormdataRequestBody = LoginUserRequest

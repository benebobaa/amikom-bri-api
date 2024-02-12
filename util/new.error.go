package util

import "errors"

var (
	UsernameAlreadyExists   = errors.New("Username already exists")
	EmailAlreadyExists      = errors.New("Email already exists")
	UsernameOrEmailNotFound = errors.New("Username or Email not found")
	InvalidPassword         = errors.New("Invalid password")
	EmailNotVerified        = errors.New("Email not verified")
	EmailVerifyAlreadyUsed  = errors.New("Email verify already used")
	EmailVerifyExpired      = errors.New("Email verify expired")
	EmailVerifyCodeNotValid = errors.New("Email verify code not valid")
)

package core

import "errors"

type ServiceSignUpUser struct {
	ChatId       int64
	UserName     string
	FirstName    string
	LastName     string
	LanguageCode string
}

var (
	ErrorAuthServiceEmptyInputArg         = errors.New("empty input arguments")
	ErrorAuthServiceAuthUserAlreadyExists = errors.New("error user already exists")
	ErrorAuthServiceInternalError         = errors.New("internal error")
)

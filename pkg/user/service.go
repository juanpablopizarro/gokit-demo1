package user

import "errors"

type Service interface {
	Validate(email, password string) (*User, error)
	Hash(string) (string, error)
}

type service struct{}

func NewServiceInstance() Service {
	return &service{}
}

func (service) Validate(email, password string) (*User, error) {
	return nil, nil
}

func (service) Hash(pass string) (string, error) {
	return "", errors.New("not implementd")
}

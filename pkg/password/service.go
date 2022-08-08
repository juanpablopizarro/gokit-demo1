package password

import "golang.org/x/crypto/bcrypt"

type Service interface {
	Hash(string) (string, error)
	Check(string, string) bool
}

type impl struct {
}

func NewPasswordService() Service {
	return &impl{}
}

func (impl) Hash(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), 14)
	return string(bytes), err
}

func (impl) Check(pass, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	return err == nil
}

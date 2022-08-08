package user

import (
	"fmt"
	"time"

	"github.com/go-kit/log"
)

type loggingMiddleware struct {
	logger log.Logger
	next   Service
}

func (mw loggingMiddleware) Validate(email, password string) (user *User, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "validate",
			"input", fmt.Sprintf("[email:%s password:%s]", email, password),
			"output", fmt.Sprintf("[user:%v]", user),
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	user, err = mw.next.Validate(email, password)
	return
}

func (mw loggingMiddleware) Hash(pass string) (hash string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "hash",
			"input", fmt.Sprintf("[password:%s]", pass),
			"output", fmt.Sprintf("[hash:%v]", hash),
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	hash, err = mw.next.Hash(pass)
	return
}

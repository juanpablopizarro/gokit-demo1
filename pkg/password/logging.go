package password

import (
	"fmt"
	"time"

	"github.com/go-kit/log"
)

type loggingMiddleware struct {
	logger log.Logger
	next   Service
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

func (mw loggingMiddleware) Check(pass, hash string) (check bool) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "hash",
			"input", fmt.Sprintf("[password:%s, hash:%s]", pass, hash),
			"output", fmt.Sprintf("[check:%v]", check),
			"error", nil,
			"took", time.Since(begin),
		)
	}(time.Now())
	check = mw.next.Check(pass, hash)
	return
}

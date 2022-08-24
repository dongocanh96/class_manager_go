package token

import "time"

type Maker interface {
	CreateToken(userid int64, username string, isteacher bool, duration time.Duration) (string, error)

	VerifyToken(token string) (*Payload, error)
}

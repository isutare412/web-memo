package port

import "context"

type Pinger interface {
	Name() string
	Ping(context.Context) error
}

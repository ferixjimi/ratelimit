package ratelimit

type ILimiter interface { // todo add retry after
	Allow(limit *Limit, data interface{}) (bool, interface{}, error)
}

type Result struct {
	Allowed bool
	//RetryAfter time.Duration
}

type Limiter struct {
	ds    IStore
	limit Limit
}

func NewLimiter() *Limiter {
	return &Limiter{}
}

func (l *Limiter) Allow() (*Result, error) {
	return l.AllowN(1)
}

func (l *Limiter) AllowN(n int) (*Result, error) {
	return nil, nil
}

func (l *Limiter) Wait() (*Result, error) {
	return nil, nil
}

func (l *Limiter) Tick() <-chan error {
	return nil
}

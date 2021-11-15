package ratelimit

var TempSession = &Session{make(map[string]*RateLimit)}

type Session struct {
	users map[string]*RateLimit
}

func (s *Session) AddRateLimit(addr string) *RateLimit {
	s.users[addr] = NewRateLimit()
	return s.users[addr]
}

func (s *Session) RemoveRateLimit(addr string) { delete(s.users, addr) }

func (s *Session) ByAddress(addr string) (r *RateLimit) {
	var ok bool
	if r, ok = s.users[addr]; !ok {
		r = s.AddRateLimit(addr)
	}
	return r
}

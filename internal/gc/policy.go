package gc

import "time"

// Policy describes the retention window applied during reclamation.
type Policy struct {
	TTL time.Duration
}

// DefaultPolicy retains blocks for twenty-four hours.
func DefaultPolicy() Policy {
	return Policy{TTL: 24 * time.Hour}
}

// WithTTL returns a copy of the policy with a custom window.
func (p Policy) WithTTL(ttl time.Duration) Policy {
	p.TTL = ttl
	return p
}

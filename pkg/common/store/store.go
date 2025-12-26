package store

import (
	"errors"
	"sync"
	"time"
)

type User struct {
	Phone        string    `json:"phone"`
	RegisteredAt time.Time `json:"registered_at"`
}

// OTP record
type OTPRecord struct {
	Code      string
	ExpiresAt time.Time
}

// In-memory storage with mutexes
type InMemoryStore struct {
	users map[string]User
	otps  map[string]OTPRecord

	// rate limiter: phone -> []time.Time of request timestamps
	otpRequests map[string][]time.Time

	mu sync.RWMutex
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		users:       make(map[string]User),
		otps:        make(map[string]OTPRecord),
		otpRequests: make(map[string][]time.Time),
	}
}

func (s *InMemoryStore) SaveOTP(phone, code string, expiresAt time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.otps[phone] = OTPRecord{Code: code, ExpiresAt: expiresAt}
}

func (s *InMemoryStore) VerifyOTP(phone string, otpCode string) (OTPRecord, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	otp, ok := s.otps[phone]

	if time.Now().After(otp.ExpiresAt) {
		s.DeleteOTP(phone)
		return OTPRecord{}, false
	}

	if otp.Code != otpCode {
		return OTPRecord{}, false
	}

	s.DeleteOTP(phone)

	return otp, ok
}

func (s *InMemoryStore) DeleteOTP(phone string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.otps, phone)
}

func (s *InMemoryStore) SaveUser(u User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[u.Phone] = u
}

func (s *InMemoryStore) GetUser(phone string) (User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[phone]
	return u, ok
}

func (s *InMemoryStore) ListUsers(offset, limit int, q string) ([]User, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var list []User
	for _, u := range s.users {
		if q == "" || contains(u, q) {
			list = append(list, u)
		}
	}
	total := len(list)
	// pagination
	if offset > total {
		return []User{}, total
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return list[offset:end], total
}

func contains(u User, q string) bool {
	// simple search over phone
	return q == "" || stringContains(u.Phone, q)
}

func stringContains(s, q string) bool {
	return len(q) == 0 || (len(s) >= len(q) && (indexOf(s, q) >= 0))
}

// simple indexOf to avoid importing strings in multiple places
func indexOf(s, sep string) int {
	for i := 0; i+len(sep) <= len(s); i++ {
		if s[i:i+len(sep)] == sep {
			return i
		}
	}
	return -1
}

func (s *InMemoryStore) AllowOTPRequest(phone string, max int, window time.Duration) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	arr := s.otpRequests[phone]
	// keep only timestamps within window
	var filtered []time.Time
	cutoff := now.Add(-window)
	for _, t := range arr {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}
	if len(filtered) >= max {
		// not allowed
		s.otpRequests[phone] = filtered
		return false
	}
	// append now and save
	filtered = append(filtered, now)
	s.otpRequests[phone] = filtered
	return true
}

func (s *InMemoryStore) CleanupExpiredOTPs() {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	for phone, otp := range s.otps {
		if otp.ExpiresAt.Before(now) {
			delete(s.otps, phone)
		}
	}
}

var ErrNotFound = errors.New("not found")

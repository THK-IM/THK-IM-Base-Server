package dto

import "golang.org/x/crypto/openpgp/errors"

type ThkClaims map[string]interface{}

const (
	TraceId  = "Trace-Id"
	Language = "Accept-Language"
	JwtToken = "Authorization"
)

func (m ThkClaims) PutValue(key string, value interface{}) {
	m[key] = value
}

func (m ThkClaims) GetTraceId() (*string, error) {
	return m.parseString(TraceId)
}

func (m ThkClaims) GetLanguage() (*string, error) {
	return m.parseString(Language)
}

func (m ThkClaims) GetToken() (*string, error) {
	return m.parseString(JwtToken)
}

func (m ThkClaims) parseString(key string) (*string, error) {
	var cs *string = nil
	switch v := m[key].(type) {
	case *string:
		cs = v
	default:
		return nil, errors.ErrKeyIncorrect
	}
	return cs, nil
}

func (m ThkClaims) parseInt(key string) (*int, error) {
	var cs *int = nil
	switch v := m[key].(type) {
	case *int:
		cs = v
	default:
		return nil, errors.ErrKeyIncorrect
	}
	return cs, nil
}

func (m ThkClaims) parseInt64(key string) (*int64, error) {
	var cs *int64 = nil
	switch v := m[key].(type) {
	case *int64:
		cs = v
	default:
		return nil, errors.ErrKeyIncorrect
	}
	return cs, nil
}

func (m ThkClaims) parseStringArray(key string) ([]string, error) {
	var cs []string = nil
	switch v := m[key].(type) {
	case []string:
		cs = v
	default:
		return nil, errors.ErrKeyIncorrect
	}
	return cs, nil
}

func (m ThkClaims) parseInt64Array(key string) ([]int64, error) {
	var cs []int64 = nil
	switch v := m[key].(type) {
	case []int64:
		cs = v
	default:
		return nil, errors.ErrKeyIncorrect
	}
	return cs, nil
}

func (m ThkClaims) parseIntArray(key string) ([]int, error) {
	var cs []int = nil
	switch v := m[key].(type) {
	case []int:
		cs = v
	default:
		return nil, errors.ErrKeyIncorrect
	}
	return cs, nil
}

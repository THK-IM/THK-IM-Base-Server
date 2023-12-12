package dto

import "golang.org/x/crypto/openpgp/errors"

type MapClaims map[string]interface{}

const (
	TraceId  = "traceId"
	Language = "language"
)

func (m MapClaims) PutValue(key string, value interface{}) {
	m[key] = value
}

func (m MapClaims) GetTraceId() (*string, error) {
	return m.parseString(TraceId)
}

func (m MapClaims) GetLanguage() (*string, error) {
	return m.parseString(Language)
}

func (m MapClaims) parseString(key string) (*string, error) {
	var cs *string = nil
	switch v := m[key].(type) {
	case *string:
		cs = v
	default:
		return nil, errors.ErrKeyIncorrect
	}
	return cs, nil
}

func (m MapClaims) parseInt(key string) (*int, error) {
	var cs *int = nil
	switch v := m[key].(type) {
	case *int:
		cs = v
	default:
		return nil, errors.ErrKeyIncorrect
	}
	return cs, nil
}

func (m MapClaims) parseInt64(key string) (*int64, error) {
	var cs *int64 = nil
	switch v := m[key].(type) {
	case *int64:
		cs = v
	default:
		return nil, errors.ErrKeyIncorrect
	}
	return cs, nil
}

func (m MapClaims) parseStringArray(key string) ([]string, error) {
	var cs []string = nil
	switch v := m[key].(type) {
	case []string:
		cs = v
	default:
		return nil, errors.ErrKeyIncorrect
	}
	return cs, nil
}

func (m MapClaims) parseInt64Array(key string) ([]int64, error) {
	var cs []int64 = nil
	switch v := m[key].(type) {
	case []int64:
		cs = v
	default:
		return nil, errors.ErrKeyIncorrect
	}
	return cs, nil
}

func (m MapClaims) parseIntArray(key string) ([]int, error) {
	var cs []int = nil
	switch v := m[key].(type) {
	case []int:
		cs = v
	default:
		return nil, errors.ErrKeyIncorrect
	}
	return cs, nil
}

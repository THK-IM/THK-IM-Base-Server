package dto

import "golang.org/x/crypto/openpgp/errors"

type ThkClaims map[string]interface{}

const (
	TraceID        = "TraceID"
	ParentSpanID   = "ParentSpanID"
	SpanID         = "SpanID"
	Language       = "Accept-Language"
	JwtToken       = "Authorization"
	ClientPlatform = "Client-Platform" // web/ios/android/centos/windows/apple
	ClientVersion  = "Client-Version"
	ClientIP       = "Client-IP"
)

func (m ThkClaims) PutValue(key string, value interface{}) {
	m[key] = value
}

func (m ThkClaims) GetTraceId() (*string, error) {
	return m.parseString(TraceID)
}

func (m ThkClaims) GetParentSpanID() (*string, error) {
	return m.parseString(ParentSpanID)
}

func (m ThkClaims) GetSpanID() (*string, error) {
	return m.parseString(SpanID)
}

func (m ThkClaims) GetLanguage() (*string, error) {
	return m.parseString(Language)
}

func (m ThkClaims) GetClientPlatform() (*string, error) {
	return m.parseString(ClientPlatform)
}

func (m ThkClaims) GetClientVersion() (*string, error) {
	return m.parseString(ClientVersion)
}

func (m ThkClaims) GetClientIp() (*string, error) {
	return m.parseString(ClientIP)
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

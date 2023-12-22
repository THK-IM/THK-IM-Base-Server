package dto

import "github.com/sirupsen/logrus"

type ThkClaims logrus.Fields

const (
	TraceID        = "TraceID"
	ParentSpanID   = "ParentSpanID"
	SpanID         = "SpanID"
	Language       = "Accept-Language"
	JwtToken       = "Authorization"
	ClientPlatform = "Client-Platform" // web/ios/android/centos/windows/apple
	ClientVersion  = "Client-Version"
	ClientOriginIP = "Client-Origin-IP"
)

func (m ThkClaims) PutValue(key string, value string) {
	m[key] = value
}

func (m ThkClaims) GetTraceId() string {
	return m.getValue(TraceID)
}

func (m ThkClaims) GetParentSpanID() string {
	return m.getValue(ParentSpanID)
}

func (m ThkClaims) GetSpanID() string {
	return m.getValue(SpanID)
}

func (m ThkClaims) GetLanguage() string {
	return m.getValue(Language)
}

func (m ThkClaims) GetClientPlatform() string {
	return m.getValue(ClientPlatform)
}

func (m ThkClaims) GetClientVersion() string {
	return m.getValue(ClientVersion)
}

func (m ThkClaims) GetClientOriginIP() string {
	return m.getValue(ClientOriginIP)
}

func (m ThkClaims) GetToken() string {
	return m.getValue(JwtToken)
}

func (m ThkClaims) getValue(key string) string {
	value, ok := m[key].(string)
	if ok {
		return value
	} else {
		return ""
	}
}

package dto

import "github.com/sirupsen/logrus"

type ThkClaims logrus.Fields

const (
	TraceID      = "TraceID"
	ParentSpanID = "ParentSpanID"
	SpanID       = "SpanID"
	Language     = "Accept-Language"
	JwtToken     = "Authorization"
	Device       = "Device"
	TimeZone     = "TimeZone"
	Platform     = "Platform" // web/ios/android/centos/windows/apple
	Version      = "Version"
	OriginIP     = "Origin-IP"

	ClaimsKey = "Claims"
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

func (m ThkClaims) GetDevice() string {
	return m.getValue(Device)
}

func (m ThkClaims) GetTimeZone() string {
	return m.getValue(TimeZone)
}

func (m ThkClaims) GetPlatform() string {
	return m.getValue(Platform)
}

func (m ThkClaims) GetVersion() string {
	return m.getValue(Version)
}

func (m ThkClaims) GetOriginIp() string {
	return m.getValue(OriginIP)
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

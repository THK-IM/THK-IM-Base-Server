package dto

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

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
	Channel      = "Channel"  // channel
	Version      = "Version"
	OriginIP     = "Origin-IP"
	DeviceId     = "DeviceId"

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

func (m ThkClaims) GetDeviceId() string {
	return m.getValue(DeviceId)
}

func (m ThkClaims) GetTimeZone() string {
	return m.getValue(TimeZone)
}

func (m ThkClaims) GetPlatform() string {
	return m.getValue(Platform)
}

func (m ThkClaims) GetChannel() string {
	return m.getValue(Channel)
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

func (m ThkClaims) ToJsonString() string {
	b, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(b)
}

func ThkClaimsFromJsonString(js string) (ThkClaims, error) {
	claims := ThkClaims{}
	err := json.Unmarshal([]byte(js), &claims)
	return claims, err
}

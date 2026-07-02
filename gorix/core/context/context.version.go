package context

import "strings"

type APIVersion string

const (
	VersionNeutral APIVersion = ""
	V1             APIVersion = "v1"
	V2             APIVersion = "v2"
	V3             APIVersion = "v3"
)

func Version(value string) APIVersion {
	value = strings.TrimSpace(value)
	value = strings.Trim(value, "/")

	if value == "" {
		return VersionNeutral
	}

	return APIVersion(value)
}

func (v APIVersion) String() string {
	return string(v)
}

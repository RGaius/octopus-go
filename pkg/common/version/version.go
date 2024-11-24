package version

var (
	// Version version
	Version string
	// BuildDate build date
	BuildDate string
)

const defaultVersion = "v0.1.0"

// Get 获取版本号
func Get() string {
	if Version == "" {
		return defaultVersion
	}

	return Version
}

// GetRevision 获取完整版本号信息，包括时间戳的
func GetRevision() string {
	if Version == "" || BuildDate == "" {
		return defaultVersion
	}

	return Version + "." + BuildDate
}

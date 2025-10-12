package utils

func CreateFingerprint(ip, userAgent string) string {
	return ip + userAgent
}

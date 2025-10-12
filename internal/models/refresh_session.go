package models

type RefreshSession struct {
	UUID         string
	RefreshToken string
	Fingerprint  string
	ExpiresIn    int
	IP           string
}

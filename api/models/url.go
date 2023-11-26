// url.go
package models

type URL struct {
	ID          string `json:"id"`
	CurrentUrl  string `json:"currentUrl"`
	RedirectUrl string `json:"redirectUrl"`
}

package service

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func extractEmailUsername(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return ""
	}

	username := parts[0]
	caser := cases.Title(language.AmericanEnglish)

	return caser.String(username)
}

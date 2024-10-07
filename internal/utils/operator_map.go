package utils

import (
	"errors"
	"strings"
)

func GetOperator(mno string) (string, error) {
	mno = strings.ToLower(mno)
	switch mno {
	case "grameenphone", "grameenphoneskitto":
		return "GP", nil
	case "robi":
		return "RB", nil
	case "banglalink":
		return "BL", nil
	case "airtel":
		return "RB", nil
	default:
		return "", errors.New("invalid MNO operator")
	}
}

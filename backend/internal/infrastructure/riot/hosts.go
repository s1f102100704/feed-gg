package riot

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidRegion = errors.New("invalid region")
)

func routeHosts(platformRegion string) (regionalHost string, platformHost string, normalizedRegion string, err error) {
	region := strings.ToUpper(strings.TrimSpace(platformRegion))

	platformToRegional := map[string]string{
		"BR1":  "americas",
		"LA1":  "americas",
		"LA2":  "americas",
		"NA1":  "americas",
		"JP1":  "asia",
		"KR":   "asia",
		"EUN1": "europe",
		"EUW1": "europe",
		"ME1":  "europe",
		"RU":   "europe",
		"TR1":  "europe",
		"OC1":  "sea",
		"PH2":  "sea",
		"SG2":  "sea",
		"TH2":  "sea",
		"TW2":  "sea",
		"VN2":  "sea",
	}

	regionalKey, ok := platformToRegional[region]
	if !ok {
		return "", "", "", ErrInvalidRegion
	}

	platformHost = fmt.Sprintf("https://%s.api.riotgames.com", strings.ToLower(region))
	regionalHost = fmt.Sprintf("https://%s.api.riotgames.com", regionalKey)
	return regionalHost, platformHost, region, nil
}

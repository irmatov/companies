package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
)

// Country is a middleware that checks if the request is coming from a specific country.
// Cacheing of requests is needed, and a circuit breaker, but..
type Country struct {
	Next               http.Handler
	LookupURLFormat    string
	AllowedCountryCode string
	Client             *http.Client
}

func (m *Country) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ipport := r.RemoteAddr
	host, _, err := net.SplitHostPort(ipport)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	countryCode, err := m.lookupCountry(host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if countryCode != m.AllowedCountryCode {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	m.Next.ServeHTTP(w, r)
}

func (m *Country) lookupCountry(ip string) (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(m.LookupURLFormat, ip), nil)
	if err != nil {
		return "", err
	}
	resp, err := m.Client.Do(req)
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("internal server error")
	}
	var data struct {
		CountryCode string `json:"country_code"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	return data.CountryCode, nil
}

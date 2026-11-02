package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// reverseGeocode resolves coordinates to an Indian address via OpenStreetMap Nominatim.
func reverseGeocode(ctx context.Context, lat, lon float64) (Address, error) {
	q := url.Values{"format": {"jsonv2"}, "lat": {fmt.Sprint(lat)}, "lon": {fmt.Sprint(lon)}, "addressdetails": {"1"}, "zoom": {"18"}}
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://nominatim.openstreetmap.org/reverse?"+q.Encode(), nil)
	req.Header.Set("User-Agent", "a-saarthi/0.1 (accessibility frontend for Ayushman Bharat)")
	cl := http.Client{Timeout: 8 * time.Second}
	res, err := cl.Do(req)
	if err != nil {
		return Address{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return Address{}, fmt.Errorf("nominatim: %s", res.Status)
	}
	var body struct {
		DisplayName string            `json:"display_name"`
		Address     map[string]string `json:"address"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return Address{}, err
	}
	a := body.Address
	first := func(keys ...string) string {
		for _, k := range keys {
			if a[k] != "" {
				return a[k]
			}
		}
		return ""
	}
	return Address{
		Line:     body.DisplayName,
		Village:  first("village", "hamlet", "suburb", "town", "city", "city_district"),
		District: first("state_district", "county", "district"),
		State:    a["state"],
		Pincode:  a["postcode"],
		Lat:      lat,
		Lon:      lon,
		Detected: true,
	}, nil
}

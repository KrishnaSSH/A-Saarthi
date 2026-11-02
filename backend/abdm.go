package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

// ABDM gateway client. Docs: https://docs.abdm.gov.in/
// Credentials come from ABDM_CLIENT_ID / ABDM_CLIENT_SECRET (issued via the ABDM sandbox/production onboarding).
// ABDM_BASE defaults to the sandbox; set it to the production host once approved.

type ABDM struct {
	base, id, secret string
	http             http.Client

	mu     sync.Mutex
	token  string
	expiry time.Time
}

var errNotConfigured = errors.New("ABDM credentials not configured")

func newABDM() *ABDM {
	base := os.Getenv("ABDM_BASE")
	if base == "" {
		base = "https://dev.abdm.gov.in"
	}
	return &ABDM{base: base, id: os.Getenv("ABDM_CLIENT_ID"), secret: os.Getenv("ABDM_CLIENT_SECRET"), http: http.Client{Timeout: 15 * time.Second}}
}

func (c *ABDM) configured() bool { return c.id != "" && c.secret != "" }

func (c *ABDM) session(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.token != "" && time.Now().Before(c.expiry) {
		return c.token, nil
	}
	if !c.configured() {
		return "", errNotConfigured
	}
	b, _ := json.Marshal(map[string]string{"clientId": c.id, "clientSecret": c.secret, "grantType": "client_credentials"})
	req, _ := http.NewRequestWithContext(ctx, "POST", c.base+"/api/hiecm/gateway/v3/sessions", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("REQUEST-ID", newUUID())
	req.Header.Set("TIMESTAMP", time.Now().UTC().Format("2006-01-02T15:04:05.000Z"))
	res, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", fmt.Errorf("session: %s", res.Status)
	}
	var out struct {
		AccessToken string `json:"accessToken"`
		ExpiresIn   int    `json:"expiresIn"`
	}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return "", err
	}
	c.token = out.AccessToken
	c.expiry = time.Now().Add(time.Duration(out.ExpiresIn-30) * time.Second)
	return c.token, nil
}

// Do sends an authenticated JSON request to the ABDM host and returns status and body.
func (c *ABDM) Do(ctx context.Context, method, path string, payload any) (int, json.RawMessage, error) {
	tok, err := c.session(ctx)
	if err != nil {
		return 0, nil, err
	}
	var r *bytes.Reader
	if payload != nil {
		b, _ := json.Marshal(payload)
		r = bytes.NewReader(b)
	} else {
		r = bytes.NewReader(nil)
	}
	req, _ := http.NewRequestWithContext(ctx, method, c.base+path, r)
	req.Header.Set("Authorization", "Bearer "+tok)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("REQUEST-ID", newUUID())
	req.Header.Set("TIMESTAMP", time.Now().UTC().Format("2006-01-02T15:04:05.000Z"))
	res, err := c.http.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	var raw json.RawMessage
	_ = json.NewDecoder(res.Body).Decode(&raw)
	return res.StatusCode, raw, nil
}

// mapUpstream converts an ABDM failure into one of our reason codes.
func mapUpstream(status int, body json.RawMessage) string {
	switch {
	case status == 404:
		return "ABHA_NOT_FOUND"
	case status >= 500 || status == 0:
		return "GATEWAY_UNAVAILABLE"
	}
	var e struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	_ = json.Unmarshal(body, &e)
	switch e.Code {
	case "ABDM-1114", "ABDM-1115":
		return "OTP_EXPIRED"
	case "ABDM-1117":
		return "OTP_INVALID"
	}
	return "UNKNOWN"
}

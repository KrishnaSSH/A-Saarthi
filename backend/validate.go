package main

import (
	"regexp"
	"strings"
	"time"
)

type Application struct {
	Name    string  `json:"name"`
	DOB     string  `json:"dob"` // DD-MM-YYYY
	Mobile  string  `json:"mobile"`
	Pincode string  `json:"pincode"`
	Address Address `json:"address"`
}

type Address struct {
	Line     string  `json:"line"`
	Village  string  `json:"village"`
	District string  `json:"district"`
	State    string  `json:"state"`
	Pincode  string  `json:"pincode"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	Detected bool    `json:"detected"`
}

var (
	pinRe    = regexp.MustCompile(`^[1-9][0-9]{5}$`)
	mobileRe = regexp.MustCompile(`^[6-9][0-9]{9}$`)
)

// check validates an application before submission and returns reason codes.
func check(a Application) []string {
	var codes []string
	if len(strings.Fields(a.Name)) < 2 {
		codes = append(codes, "NAME_MISMATCH")
	}
	if d, err := time.Parse("02-01-2006", a.DOB); err != nil || d.After(time.Now()) {
		codes = append(codes, "DOB_MISMATCH")
	}
	if !mobileRe.MatchString(a.Mobile) {
		codes = append(codes, "MOBILE_NOT_LINKED")
	}
	if !pinRe.MatchString(a.Pincode) {
		codes = append(codes, "PINCODE_INVALID")
	} else if a.Address.Detected && a.Address.Pincode != "" && a.Address.Pincode != a.Pincode {
		codes = append(codes, "PINCODE_MISMATCH")
	}
	ad := a.Address
	if ad.District == "" || ad.State == "" || (ad.Village == "" && ad.Line == "") {
		codes = append(codes, "ADDRESS_INCOMPLETE")
	}
	if !ad.Detected && (ad.Lat == 0 && ad.Lon == 0) {
		codes = append(codes, "LOCATION_UNVERIFIED")
	}
	return codes
}

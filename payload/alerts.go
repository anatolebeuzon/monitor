package payload

import "time"

// Alerts maps from a website URL to an Alert.
type Alerts map[string]Alert

// Alert represents an alert for a particular website.
type Alert struct {
	Date         time.Time // Date at which the alert was created
	Availability float64   // Average availability of the website

	// BelowThreshold indicates whether the website is considered
	// down (new alert) or up (recovery alert)
	BelowThreshold bool
}

// NewAlert creates and returns a new alert.
func NewAlert(availability float64, belowThreshold bool) Alert {
	return Alert{time.Now(), availability, belowThreshold}
}

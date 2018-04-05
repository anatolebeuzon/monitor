package payload

// Alerts maps from a website URL to an Alert.
type Alerts map[string]Alert

// Alert represents an alert for a particular website.
type Alert struct {
	Timeframe    Timeframe // Time window use to aggregate results
	Availability float64   // Average availability of the website

	// BelowThreshold indicates whether the website is considered
	// down (new alert) or up (recovery alert)
	BelowThreshold bool
}

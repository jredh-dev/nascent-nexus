// Package fees provides delivery fee calculation for Hooper Works.
//
// All calculations are based on configurable rates with sensible defaults
// derived from Washington State minimums and typical vehicle costs.
//
// Reference point for distance: Henry M. Jackson Federal Building,
// 915 2nd Ave, Seattle, WA 98174.
package fees

import "math"

// Rates holds the configurable values used for fee calculation.
type Rates struct {
	// HourlyWage is the labor rate in dollars per hour.
	// Default: WA State minimum wage ($20/hr).
	HourlyWage float64

	// GasPrice is the price per gallon in dollars.
	// Default: $5.00 (rounded WA average for easy calculation).
	GasPrice float64

	// MPG is the vehicle fuel efficiency in miles per gallon.
	// Default: 20 MPG (conservative best-case estimate).
	MPG float64
}

// DefaultRates returns the standard Hooper Works fee rates.
func DefaultRates() Rates {
	return Rates{
		HourlyWage: 20.0,
		MPG:        20.0,
		GasPrice:   5.0,
	}
}

// DeliveryFee is the itemized breakdown of a delivery charge.
type DeliveryFee struct {
	// LaborCost is the time-based charge (both directions).
	LaborCost float64 `json:"labor_cost"`

	// FuelCost is the distance-based charge (both directions).
	FuelCost float64 `json:"fuel_cost"`

	// Total is LaborCost + FuelCost.
	Total float64 `json:"total"`

	// Miles is the one-way distance used in calculation.
	Miles float64 `json:"miles"`

	// DriveMinutes is the one-way drive time used in calculation.
	DriveMinutes int `json:"drive_minutes"`

	// RoundTripMiles is the total distance (both directions).
	RoundTripMiles float64 `json:"round_trip_miles"`

	// RoundTripMinutes is the total drive time (both directions).
	RoundTripMinutes int `json:"round_trip_minutes"`
}

// CalculateDelivery computes the delivery fee for a given one-way distance
// and drive time. Both directions are billed (you drive there and back).
//
// Labor: (driveMinutes * 2) / 60 * hourlyWage
// Fuel:  (miles * 2) / mpg * gasPrice
func CalculateDelivery(miles float64, driveMinutes int, rates Rates) DeliveryFee {
	if miles < 0 {
		miles = 0
	}
	if driveMinutes < 0 {
		driveMinutes = 0
	}

	roundTripMiles := miles * 2
	roundTripMinutes := driveMinutes * 2

	laborHours := float64(roundTripMinutes) / 60.0
	laborCost := laborHours * rates.HourlyWage

	gallons := roundTripMiles / rates.MPG
	fuelCost := gallons * rates.GasPrice

	total := laborCost + fuelCost

	return DeliveryFee{
		LaborCost:        roundCents(laborCost),
		FuelCost:         roundCents(fuelCost),
		Total:            roundCents(total),
		Miles:            miles,
		DriveMinutes:     driveMinutes,
		RoundTripMiles:   roundTripMiles,
		RoundTripMinutes: roundTripMinutes,
	}
}

// CalculateDeliveryDefault computes the delivery fee using default rates.
func CalculateDeliveryDefault(miles float64, driveMinutes int) DeliveryFee {
	return CalculateDelivery(miles, driveMinutes, DefaultRates())
}

// roundCents rounds a dollar amount to the nearest cent.
func roundCents(amount float64) float64 {
	return math.Round(amount*100) / 100
}

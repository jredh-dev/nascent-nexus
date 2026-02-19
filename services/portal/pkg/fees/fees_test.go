package fees

import (
	"testing"
)

func TestDefaultRates(t *testing.T) {
	r := DefaultRates()
	if r.HourlyWage != 20.0 {
		t.Errorf("HourlyWage = %v, want 20.0", r.HourlyWage)
	}
	if r.GasPrice != 5.0 {
		t.Errorf("GasPrice = %v, want 5.0", r.GasPrice)
	}
	if r.MPG != 20.0 {
		t.Errorf("MPG = %v, want 20.0", r.MPG)
	}
}

func TestCalculateDelivery_ExampleFromSpec(t *testing.T) {
	// From the spec: 5mi, 15min one-way
	// Labor: 15*2=30min = 0.5hr * $20 = $10.00
	// Fuel:  5*2=10mi / 20mpg = 0.5gal * $5 = $2.50
	// Total: $12.50
	fee := CalculateDeliveryDefault(5.0, 15)

	if fee.LaborCost != 10.0 {
		t.Errorf("LaborCost = %v, want 10.0", fee.LaborCost)
	}
	if fee.FuelCost != 2.50 {
		t.Errorf("FuelCost = %v, want 2.50", fee.FuelCost)
	}
	if fee.Total != 12.50 {
		t.Errorf("Total = %v, want 12.50", fee.Total)
	}
	if fee.Miles != 5.0 {
		t.Errorf("Miles = %v, want 5.0", fee.Miles)
	}
	if fee.DriveMinutes != 15 {
		t.Errorf("DriveMinutes = %v, want 15", fee.DriveMinutes)
	}
	if fee.RoundTripMiles != 10.0 {
		t.Errorf("RoundTripMiles = %v, want 10.0", fee.RoundTripMiles)
	}
	if fee.RoundTripMinutes != 30 {
		t.Errorf("RoundTripMinutes = %v, want 30", fee.RoundTripMinutes)
	}
}

func TestCalculateDelivery_ZeroDistance(t *testing.T) {
	fee := CalculateDeliveryDefault(0, 0)
	if fee.Total != 0 {
		t.Errorf("Total = %v, want 0 for zero distance", fee.Total)
	}
	if fee.LaborCost != 0 {
		t.Errorf("LaborCost = %v, want 0 for zero time", fee.LaborCost)
	}
	if fee.FuelCost != 0 {
		t.Errorf("FuelCost = %v, want 0 for zero distance", fee.FuelCost)
	}
}

func TestCalculateDelivery_NegativeInputs(t *testing.T) {
	fee := CalculateDeliveryDefault(-5, -10)
	if fee.Total != 0 {
		t.Errorf("Total = %v, want 0 for negative inputs", fee.Total)
	}
	if fee.Miles != 0 {
		t.Errorf("Miles = %v, want 0 for negative input", fee.Miles)
	}
}

func TestCalculateDelivery_CustomRates(t *testing.T) {
	rates := Rates{
		HourlyWage: 30.0,
		GasPrice:   6.0,
		MPG:        15.0,
	}
	// 10mi, 20min one-way
	// Labor: 20*2=40min = 0.6667hr * $30 = $20.00
	// Fuel:  10*2=20mi / 15mpg = 1.3333gal * $6 = $8.00
	// Total: $28.00
	fee := CalculateDelivery(10.0, 20, rates)

	if fee.LaborCost != 20.0 {
		t.Errorf("LaborCost = %v, want 20.0", fee.LaborCost)
	}
	if fee.FuelCost != 8.0 {
		t.Errorf("FuelCost = %v, want 8.0", fee.FuelCost)
	}
	if fee.Total != 28.0 {
		t.Errorf("Total = %v, want 28.0", fee.Total)
	}
}

func TestCalculateDelivery_RoundingCents(t *testing.T) {
	// 3mi, 7min one-way
	// Labor: 7*2=14min = 0.2333hr * $20 = $4.6667 -> $4.67
	// Fuel:  3*2=6mi / 20mpg = 0.3gal * $5 = $1.50
	// Total: $6.17
	fee := CalculateDeliveryDefault(3.0, 7)

	if fee.LaborCost != 4.67 {
		t.Errorf("LaborCost = %v, want 4.67", fee.LaborCost)
	}
	if fee.FuelCost != 1.50 {
		t.Errorf("FuelCost = %v, want 1.50", fee.FuelCost)
	}
	if fee.Total != 6.17 {
		t.Errorf("Total = %v, want 6.17", fee.Total)
	}
}

func TestCalculateDelivery_LargeDistance(t *testing.T) {
	// 50mi, 60min one-way
	// Labor: 60*2=120min = 2hr * $20 = $40.00
	// Fuel:  50*2=100mi / 20mpg = 5gal * $5 = $25.00
	// Total: $65.00
	fee := CalculateDeliveryDefault(50.0, 60)

	if fee.LaborCost != 40.0 {
		t.Errorf("LaborCost = %v, want 40.0", fee.LaborCost)
	}
	if fee.FuelCost != 25.0 {
		t.Errorf("FuelCost = %v, want 25.0", fee.FuelCost)
	}
	if fee.Total != 65.0 {
		t.Errorf("Total = %v, want 65.0", fee.Total)
	}
}

func TestCalculateDelivery_FractionalMiles(t *testing.T) {
	// 2.5mi, 8min one-way
	// Labor: 8*2=16min = 0.2667hr * $20 = $5.3333 -> $5.33
	// Fuel:  2.5*2=5mi / 20mpg = 0.25gal * $5 = $1.25
	// Total: $6.58
	fee := CalculateDeliveryDefault(2.5, 8)

	if fee.LaborCost != 5.33 {
		t.Errorf("LaborCost = %v, want 5.33", fee.LaborCost)
	}
	if fee.FuelCost != 1.25 {
		t.Errorf("FuelCost = %v, want 1.25", fee.FuelCost)
	}
	if fee.Total != 6.58 {
		t.Errorf("Total = %v, want 6.58", fee.Total)
	}
}

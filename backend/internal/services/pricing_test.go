package services

import "testing"

func TestEstimateInvoicePrice(t *testing.T) {
	ps := NewPricingService()

	tests := []struct {
		name     string
		msrp     float64
		make     string
		bodyType string
		want     float64
	}{
		{
			name:     "Toyota sedan with 7-9% margin (should use lower end)",
			msrp:     40000,
			make:     "Toyota",
			bodyType: "Sedan",
			want:     37200, // 40000 * (1 - 0.07) = 37200
		},
		{
			name:     "Toyota truck with 7-9% margin (should use higher end)",
			msrp:     40000,
			make:     "Toyota",
			bodyType: "Pickup Truck",
			want:     36400, // 40000 * (1 - 0.09) = 36400
		},
		{
			name:     "Tesla direct-to-consumer (should return MSRP)",
			msrp:     50000,
			make:     "Tesla",
			bodyType: "Sedan",
			want:     50000,
		},
		{
			name:     "Unknown brand (should use default margins)",
			msrp:     40000,
			make:     "UnknownBrand",
			bodyType: "Sedan",
			want:     38000, // 40000 * (1 - 0.05) = 38000 (lower end of default)
		},
		{
			name:     "Zero MSRP",
			msrp:     0,
			make:     "Toyota",
			bodyType: "Sedan",
			want:     0,
		},
		{
			name:     "Honda with SUV body type",
			msrp:     35000,
			make:     "Honda",
			bodyType: "SUV",
			want:     32025, // 35000 * (1 - 0.085) = 32025 (higher end of 7-8.5%)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ps.EstimateInvoicePrice(tt.msrp, tt.make, tt.bodyType)
			if got != tt.want {
				t.Errorf("EstimateInvoicePrice(%v, %q, %q) = %v, want %v", tt.msrp, tt.make, tt.bodyType, got, tt.want)
			}
		})
	}
}

func TestEstimateDealerHoldback(t *testing.T) {
	ps := NewPricingService()

	tests := []struct {
		name string
		msrp float64
		make string
		want float64
	}{
		{
			name: "Acura with 2% holdback",
			msrp: 40000,
			make: "Acura",
			want: 800, // 40000 * 0.02 = 800
		},
		{
			name: "Audi with no holdback",
			msrp: 50000,
			make: "Audi",
			want: 0,
		},
		{
			name: "BMW with no holdback",
			msrp: 60000,
			make: "BMW",
			want: 0,
		},
		{
			name: "Porsche with no holdback",
			msrp: 100000,
			make: "Porsche",
			want: 0,
		},
		{
			name: "Tesla direct-to-consumer (no holdback)",
			msrp: 50000,
			make: "Tesla",
			want: 0,
		},
		{
			name: "Chevrolet with 3% holdback",
			msrp: 40000,
			make: "Chevrolet",
			want: 1200, // 40000 * 0.03 = 1200
		},
		{
			name: "Unknown brand (default 2% holdback)",
			msrp: 40000,
			make: "UnknownBrand",
			want: 800, // 40000 * 0.02 = 800
		},
		{
			name: "Zero MSRP",
			msrp: 0,
			make: "Toyota",
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ps.EstimateDealerHoldback(tt.msrp, tt.make)
			if got != tt.want {
				t.Errorf("EstimateDealerHoldback(%v, %q) = %v, want %v", tt.msrp, tt.make, got, tt.want)
			}
		})
	}
}

func TestCalculateNetNetPrice(t *testing.T) {
	ps := NewPricingService()

	tests := []struct {
		name     string
		msrp     float64
		make     string
		bodyType string
		want     float64
	}{
		{
			name:     "Toyota sedan (invoice - holdback)",
			msrp:     40000,
			make:     "Toyota",
			bodyType: "Sedan",
			want:     36400, // Invoice: 37200, Holdback: 800, Net-Net: 36400
		},
		{
			name:     "Tesla direct-to-consumer (should return MSRP)",
			msrp:     50000,
			make:     "Tesla",
			bodyType: "Sedan",
			want:     50000,
		},
		{
			name:     "Chevrolet truck",
			msrp:     50000,
			make:     "Chevrolet",
			bodyType: "Pickup Truck",
			want:     44500, // Invoice: 46000 (using 8% max), Holdback: 1500, Net-Net: 44500
		},
		{
			name:     "Zero MSRP",
			msrp:     0,
			make:     "Toyota",
			bodyType: "Sedan",
			want:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ps.CalculateNetNetPrice(tt.msrp, tt.make, tt.bodyType)
			// Allow small floating point differences
			diff := got - tt.want
			if diff < 0 {
				diff = -diff
			}
			if diff > 1.0 {
				t.Errorf("CalculateNetNetPrice(%v, %q, %q) = %v, want %v (diff: %v)", tt.msrp, tt.make, tt.bodyType, got, tt.want, diff)
			}
		})
	}
}


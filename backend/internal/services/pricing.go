package services

import (
	"strings"
)

// BrandMarginData represents margin and holdback data for a brand
type BrandMarginData struct {
	MSRPToInvoiceMarginMinPercent float64 // e.g., 5.0 for 5%
	MSRPToInvoiceMarginMaxPercent float64 // e.g., 7.0 for 7%
	DealerHoldbackPercent         float64 // e.g., 2.0 for 2%
	HoldbackBase                  string   // "Base MSRP", "Total MSRP", or "MSRP"
	IsDirectToConsumer            bool     // true for Tesla, Rivian, Lucid
}

// PricingService handles pricing calculations
type PricingService struct {
	brandMargins map[string]BrandMarginData
}

// NewPricingService creates a new pricing service with brand margin data
func NewPricingService() *PricingService {
	ps := &PricingService{
		brandMargins: make(map[string]BrandMarginData),
	}
	ps.initializeBrandMargins()
	return ps
}

// initializeBrandMargins populates the brand margin data from the provided table
func (s *PricingService) initializeBrandMargins() {
	// Acura: 5% – 7%, 2% of Base MSRP
	s.brandMargins["Acura"] = BrandMarginData{
		MSRPToInvoiceMarginMinPercent: 5.0,
		MSRPToInvoiceMarginMaxPercent: 7.0,
		DealerHoldbackPercent:         2.0,
		HoldbackBase:                  "Base MSRP",
		IsDirectToConsumer:            false,
	}

	// Audi: 6% – 8%, 0%
	s.brandMargins["Audi"] = BrandMarginData{
		MSRPToInvoiceMarginMinPercent: 6.0,
		MSRPToInvoiceMarginMaxPercent: 8.0,
		DealerHoldbackPercent:         0.0,
		HoldbackBase:                  "MSRP",
		IsDirectToConsumer:            false,
	}

	// BMW: 6% – 8%, 0%
	s.brandMargins["BMW"] = BrandMarginData{
		MSRPToInvoiceMarginMinPercent: 6.0,
		MSRPToInvoiceMarginMaxPercent: 8.0,
		DealerHoldbackPercent:         0.0,
		HoldbackBase:                  "MSRP",
		IsDirectToConsumer:            false,
	}

	// Buick: 4% – 6%, 3% of MSRP
	s.brandMargins["Buick"] = BrandMarginData{
		MSRPToInvoiceMarginMinPercent: 4.0,
		MSRPToInvoiceMarginMaxPercent: 6.0,
		DealerHoldbackPercent:         3.0,
		HoldbackBase:                  "MSRP",
		IsDirectToConsumer:            false,
	}

	// Cadillac: 5% – 7%, 3% of MSRP
	s.brandMargins["Cadillac"] = BrandMarginData{
		MSRPToInvoiceMarginMinPercent: 5.0,
		MSRPToInvoiceMarginMaxPercent: 7.0,
		DealerHoldbackPercent:         3.0,
		HoldbackBase:                  "MSRP",
		IsDirectToConsumer:            false,
	}

	// Chevrolet: 5% – 8%, 3% of MSRP
	s.brandMargins["Chevrolet"] = BrandMarginData{
		MSRPToInvoiceMarginMinPercent: 5.0,
		MSRPToInvoiceMarginMaxPercent: 8.0,
		DealerHoldbackPercent:         3.0,
		HoldbackBase:                  "MSRP",
		IsDirectToConsumer:            false,
	}

	// Chrysler/Dodge: 3% – 5%, 3% of MSRP
	s.brandMargins["Chrysler"] = BrandMarginData{
		MSRPToInvoiceMarginMinPercent: 3.0,
		MSRPToInvoiceMarginMaxPercent: 5.0,
		DealerHoldbackPercent:         3.0,
		HoldbackBase:                  "MSRP",
		IsDirectToConsumer:            false,
	}
	s.brandMargins["Dodge"] = BrandMarginData{
		MSRPToInvoiceMarginMinPercent: 3.0,
		MSRPToInvoiceMarginMaxPercent: 5.0,
		DealerHoldbackPercent:         3.0,
		HoldbackBase:                  "MSRP",
		IsDirectToConsumer:            false,
	}

	// Ford: 4.5% – 6.5%, 3% of MSRP
	s.brandMargins["Ford"] = BrandMarginData{
		MSRPToInvoiceMarginMinPercent: 4.5,
		MSRPToInvoiceMarginMaxPercent: 6.5,
		DealerHoldbackPercent:         3.0,
		HoldbackBase:                  "MSRP",
		IsDirectToConsumer:            false,
	}

	// GMC: 5% – 8%, 3% of MSRP
	s.brandMargins["GMC"] = BrandMarginData{
		MSRPToInvoiceMarginMinPercent: 5.0,
		MSRPToInvoiceMarginMaxPercent: 8.0,
		DealerHoldbackPercent:         3.0,
		HoldbackBase:                  "MSRP",
		IsDirectToConsumer:            false,
	}

	// Honda: 7% – 8.5%, 2% of Base MSRP
	s.brandMargins["Honda"] = BrandMarginData{
		MSRPToInvoiceMarginMinPercent: 7.0,
		MSRPToInvoiceMarginMaxPercent: 8.5,
		DealerHoldbackPercent:         2.0,
		HoldbackBase:                  "Base MSRP",
		IsDirectToConsumer:            false,
	}

	// Hyundai/Kia: 2.5% – 4.5%, 2% of MSRP
	s.brandMargins["Hyundai"] = BrandMarginData{
		MSRPToInvoiceMarginMinPercent: 2.5,
		MSRPToInvoiceMarginMaxPercent: 4.5,
		DealerHoldbackPercent:         2.0,
		HoldbackBase:                  "MSRP",
		IsDirectToConsumer:            false,
	}
	s.brandMargins["Kia"] = BrandMarginData{
		MSRPToInvoiceMarginMinPercent: 2.5,
		MSRPToInvoiceMarginMaxPercent: 4.5,
		DealerHoldbackPercent:         2.0,
		HoldbackBase:                  "MSRP",
		IsDirectToConsumer:            false,
	}

	// Jeep/Ram: 4% – 6%, 3% of MSRP
	s.brandMargins["Jeep"] = BrandMarginData{
		MSRPToInvoiceMarginMinPercent: 4.0,
		MSRPToInvoiceMarginMaxPercent: 6.0,
		DealerHoldbackPercent:         3.0,
		HoldbackBase:                  "MSRP",
		IsDirectToConsumer:            false,
	}
	s.brandMargins["Ram"] = BrandMarginData{
		MSRPToInvoiceMarginMinPercent: 4.0,
		MSRPToInvoiceMarginMaxPercent: 6.0,
		DealerHoldbackPercent:         3.0,
		HoldbackBase:                  "MSRP",
		IsDirectToConsumer:            false,
	}

	// Lexus: 6% – 8.5%, 2% of Base MSRP
	s.brandMargins["Lexus"] = BrandMarginData{
		MSRPToInvoiceMarginMinPercent: 6.0,
		MSRPToInvoiceMarginMaxPercent: 8.5,
		DealerHoldbackPercent:         2.0,
		HoldbackBase:                  "Base MSRP",
		IsDirectToConsumer:            false,
	}

	// Mazda: 2.5% – 5%, 2% of MSRP
	s.brandMargins["Mazda"] = BrandMarginData{
		MSRPToInvoiceMarginMinPercent: 2.5,
		MSRPToInvoiceMarginMaxPercent: 5.0,
		DealerHoldbackPercent:         2.0,
		HoldbackBase:                  "MSRP",
		IsDirectToConsumer:            false,
	}

	// Mercedes-Benz: 7% – 8%, 1% – 3%
	s.brandMargins["Mercedes-Benz"] = BrandMarginData{
		MSRPToInvoiceMarginMinPercent: 7.0,
		MSRPToInvoiceMarginMaxPercent: 8.0,
		DealerHoldbackPercent:         2.0, // Average of 1-3%
		HoldbackBase:                  "MSRP",
		IsDirectToConsumer:            false,
	}

	// Nissan: 4% – 6%, 2% of Total MSRP
	s.brandMargins["Nissan"] = BrandMarginData{
		MSRPToInvoiceMarginMinPercent: 4.0,
		MSRPToInvoiceMarginMaxPercent: 6.0,
		DealerHoldbackPercent:         2.0,
		HoldbackBase:                  "Total MSRP",
		IsDirectToConsumer:            false,
	}

	// Porsche: 8% – 10%, 0%
	s.brandMargins["Porsche"] = BrandMarginData{
		MSRPToInvoiceMarginMinPercent: 8.0,
		MSRPToInvoiceMarginMaxPercent: 10.0,
		DealerHoldbackPercent:         0.0,
		HoldbackBase:                  "MSRP",
		IsDirectToConsumer:            false,
	}

	// Subaru: 6% – 7%, 2% of Total MSRP
	s.brandMargins["Subaru"] = BrandMarginData{
		MSRPToInvoiceMarginMinPercent: 6.0,
		MSRPToInvoiceMarginMaxPercent: 7.0,
		DealerHoldbackPercent:         2.0,
		HoldbackBase:                  "Total MSRP",
		IsDirectToConsumer:            false,
	}

	// Toyota: 7% – 9%, 2% of Base MSRP
	s.brandMargins["Toyota"] = BrandMarginData{
		MSRPToInvoiceMarginMinPercent: 7.0,
		MSRPToInvoiceMarginMaxPercent: 9.0,
		DealerHoldbackPercent:         2.0,
		HoldbackBase:                  "Base MSRP",
		IsDirectToConsumer:            false,
	}

	// Volkswagen: 4% – 5.5%, 2% of MSRP
	s.brandMargins["Volkswagen"] = BrandMarginData{
		MSRPToInvoiceMarginMinPercent: 4.0,
		MSRPToInvoiceMarginMaxPercent: 5.5,
		DealerHoldbackPercent:         2.0,
		HoldbackBase:                  "MSRP",
		IsDirectToConsumer:            false,
	}

	// Volvo: 4% – 6%, 1% of MSRP
	s.brandMargins["Volvo"] = BrandMarginData{
		MSRPToInvoiceMarginMinPercent: 4.0,
		MSRPToInvoiceMarginMaxPercent: 6.0,
		DealerHoldbackPercent:         1.0,
		HoldbackBase:                  "MSRP",
		IsDirectToConsumer:            false,
	}

	// Direct-to-consumer brands
	s.brandMargins["Tesla"] = BrandMarginData{
		MSRPToInvoiceMarginMinPercent: 0.0,
		MSRPToInvoiceMarginMaxPercent: 0.0,
		DealerHoldbackPercent:         0.0,
		HoldbackBase:                  "MSRP",
		IsDirectToConsumer:            true,
	}
	s.brandMargins["Rivian"] = BrandMarginData{
		MSRPToInvoiceMarginMinPercent: 0.0,
		MSRPToInvoiceMarginMaxPercent: 0.0,
		DealerHoldbackPercent:         0.0,
		HoldbackBase:                  "MSRP",
		IsDirectToConsumer:            true,
	}
	s.brandMargins["Lucid"] = BrandMarginData{
		MSRPToInvoiceMarginMinPercent: 0.0,
		MSRPToInvoiceMarginMaxPercent: 0.0,
		DealerHoldbackPercent:         0.0,
		HoldbackBase:                  "MSRP",
		IsDirectToConsumer:            true,
	}
}

// getBrandMarginData returns margin data for a brand, or default if not found
func (s *PricingService) getBrandMarginData(make string) BrandMarginData {
	// Try exact match first (for brands like BMW, GMC, etc.)
	if data, exists := s.brandMargins[make]; exists {
		return data
	}
	
	// Normalize make name (capitalize first letter, rest lowercase)
	normalized := strings.Title(strings.ToLower(make))
	if data, exists := s.brandMargins[normalized]; exists {
		return data
	}
	
	// Try uppercase (for brands like BMW)
	upper := strings.ToUpper(make)
	if data, exists := s.brandMargins[upper]; exists {
		return data
	}
	
	// Default margins for unknown brands
	return BrandMarginData{
		MSRPToInvoiceMarginMinPercent: 5.0,
		MSRPToInvoiceMarginMaxPercent: 7.0,
		DealerHoldbackPercent:         2.0,
		HoldbackBase:                  "MSRP",
		IsDirectToConsumer:            false,
	}
}

// isTruckOrLargeSUV determines if a body type is a truck or large SUV
func (s *PricingService) isTruckOrLargeSUV(bodyType string) bool {
	if bodyType == "" {
		return false
	}
	bodyTypeLower := strings.ToLower(bodyType)
	return strings.Contains(bodyTypeLower, "truck") ||
		strings.Contains(bodyTypeLower, "pickup") ||
		strings.Contains(bodyTypeLower, "suv") ||
		strings.Contains(bodyTypeLower, "sport utility")
}

// EstimateInvoicePrice calculates the estimated invoice price from MSRP
func (s *PricingService) EstimateInvoicePrice(msrp float64, make string, bodyType string) float64 {
	if msrp <= 0 {
		return 0
	}

	brandData := s.getBrandMarginData(make)
	
	// Direct-to-consumer brands return MSRP
	if brandData.IsDirectToConsumer {
		return msrp
	}

	// Determine margin percentage based on vehicle type
	marginPercent := (brandData.MSRPToInvoiceMarginMinPercent + brandData.MSRPToInvoiceMarginMaxPercent) / 2.0
	
	// Adjust for truck/large SUV vs sedan/compact
	if s.isTruckOrLargeSUV(bodyType) {
		// Use higher end of margin range for trucks/SUVs
		marginPercent = brandData.MSRPToInvoiceMarginMaxPercent
	} else {
		// Use lower end of margin range for sedans/compacts
		marginPercent = brandData.MSRPToInvoiceMarginMinPercent
	}

	// Convert percent to decimal and calculate invoice
	marginDecimal := marginPercent / 100.0
	invoice := msrp * (1.0 - marginDecimal)
	
	return invoice
}

// EstimateDealerHoldback calculates the estimated dealer holdback amount
func (s *PricingService) EstimateDealerHoldback(msrp float64, make string) float64 {
	if msrp <= 0 {
		return 0
	}

	brandData := s.getBrandMarginData(make)
	
	// Direct-to-consumer brands have no holdback
	if brandData.IsDirectToConsumer {
		return 0
	}

	// Convert percent to decimal and calculate holdback
	holdbackDecimal := brandData.DealerHoldbackPercent / 100.0
	holdback := msrp * holdbackDecimal
	
	return holdback
}

// CalculateNetNetPrice calculates the net net price (invoice - holdback)
func (s *PricingService) CalculateNetNetPrice(msrp float64, make string, bodyType string) float64 {
	if msrp <= 0 {
		return 0
	}

	brandData := s.getBrandMarginData(make)
	
	// Direct-to-consumer brands return MSRP
	if brandData.IsDirectToConsumer {
		return msrp
	}

	invoice := s.EstimateInvoicePrice(msrp, make, bodyType)
	holdback := s.EstimateDealerHoldback(msrp, make)
	
	return invoice - holdback
}


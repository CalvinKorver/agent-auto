// Brand margin data structure
interface BrandMarginData {
  msrpToInvoiceMarginMinPercent: number;
  msrpToInvoiceMarginMaxPercent: number;
  dealerHoldbackPercent: number;
  holdbackBase: string; // "Base MSRP", "Total MSRP", or "MSRP"
  isDirectToConsumer: boolean;
}

// Brand margin data map
const brandMargins: Record<string, BrandMarginData> = {
  Acura: {
    msrpToInvoiceMarginMinPercent: 5.0,
    msrpToInvoiceMarginMaxPercent: 7.0,
    dealerHoldbackPercent: 2.0,
    holdbackBase: "Base MSRP",
    isDirectToConsumer: false,
  },
  Audi: {
    msrpToInvoiceMarginMinPercent: 6.0,
    msrpToInvoiceMarginMaxPercent: 8.0,
    dealerHoldbackPercent: 0.0,
    holdbackBase: "MSRP",
    isDirectToConsumer: false,
  },
  BMW: {
    msrpToInvoiceMarginMinPercent: 6.0,
    msrpToInvoiceMarginMaxPercent: 8.0,
    dealerHoldbackPercent: 0.0,
    holdbackBase: "MSRP",
    isDirectToConsumer: false,
  },
  Buick: {
    msrpToInvoiceMarginMinPercent: 4.0,
    msrpToInvoiceMarginMaxPercent: 6.0,
    dealerHoldbackPercent: 3.0,
    holdbackBase: "MSRP",
    isDirectToConsumer: false,
  },
  Cadillac: {
    msrpToInvoiceMarginMinPercent: 5.0,
    msrpToInvoiceMarginMaxPercent: 7.0,
    dealerHoldbackPercent: 3.0,
    holdbackBase: "MSRP",
    isDirectToConsumer: false,
  },
  Chevrolet: {
    msrpToInvoiceMarginMinPercent: 5.0,
    msrpToInvoiceMarginMaxPercent: 8.0,
    dealerHoldbackPercent: 3.0,
    holdbackBase: "MSRP",
    isDirectToConsumer: false,
  },
  Chrysler: {
    msrpToInvoiceMarginMinPercent: 3.0,
    msrpToInvoiceMarginMaxPercent: 5.0,
    dealerHoldbackPercent: 3.0,
    holdbackBase: "MSRP",
    isDirectToConsumer: false,
  },
  Dodge: {
    msrpToInvoiceMarginMinPercent: 3.0,
    msrpToInvoiceMarginMaxPercent: 5.0,
    dealerHoldbackPercent: 3.0,
    holdbackBase: "MSRP",
    isDirectToConsumer: false,
  },
  Ford: {
    msrpToInvoiceMarginMinPercent: 4.5,
    msrpToInvoiceMarginMaxPercent: 6.5,
    dealerHoldbackPercent: 3.0,
    holdbackBase: "MSRP",
    isDirectToConsumer: false,
  },
  GMC: {
    msrpToInvoiceMarginMinPercent: 5.0,
    msrpToInvoiceMarginMaxPercent: 8.0,
    dealerHoldbackPercent: 3.0,
    holdbackBase: "MSRP",
    isDirectToConsumer: false,
  },
  Honda: {
    msrpToInvoiceMarginMinPercent: 7.0,
    msrpToInvoiceMarginMaxPercent: 8.5,
    dealerHoldbackPercent: 2.0,
    holdbackBase: "Base MSRP",
    isDirectToConsumer: false,
  },
  Hyundai: {
    msrpToInvoiceMarginMinPercent: 2.5,
    msrpToInvoiceMarginMaxPercent: 4.5,
    dealerHoldbackPercent: 2.0,
    holdbackBase: "MSRP",
    isDirectToConsumer: false,
  },
  Kia: {
    msrpToInvoiceMarginMinPercent: 2.5,
    msrpToInvoiceMarginMaxPercent: 4.5,
    dealerHoldbackPercent: 2.0,
    holdbackBase: "MSRP",
    isDirectToConsumer: false,
  },
  Jeep: {
    msrpToInvoiceMarginMinPercent: 4.0,
    msrpToInvoiceMarginMaxPercent: 6.0,
    dealerHoldbackPercent: 3.0,
    holdbackBase: "MSRP",
    isDirectToConsumer: false,
  },
  Ram: {
    msrpToInvoiceMarginMinPercent: 4.0,
    msrpToInvoiceMarginMaxPercent: 6.0,
    dealerHoldbackPercent: 3.0,
    holdbackBase: "MSRP",
    isDirectToConsumer: false,
  },
  Lexus: {
    msrpToInvoiceMarginMinPercent: 6.0,
    msrpToInvoiceMarginMaxPercent: 8.5,
    dealerHoldbackPercent: 2.0,
    holdbackBase: "Base MSRP",
    isDirectToConsumer: false,
  },
  Mazda: {
    msrpToInvoiceMarginMinPercent: 2.5,
    msrpToInvoiceMarginMaxPercent: 5.0,
    dealerHoldbackPercent: 2.0,
    holdbackBase: "MSRP",
    isDirectToConsumer: false,
  },
  "Mercedes-Benz": {
    msrpToInvoiceMarginMinPercent: 7.0,
    msrpToInvoiceMarginMaxPercent: 8.0,
    dealerHoldbackPercent: 2.0, // Average of 1-3%
    holdbackBase: "MSRP",
    isDirectToConsumer: false,
  },
  Nissan: {
    msrpToInvoiceMarginMinPercent: 4.0,
    msrpToInvoiceMarginMaxPercent: 6.0,
    dealerHoldbackPercent: 2.0,
    holdbackBase: "Total MSRP",
    isDirectToConsumer: false,
  },
  Porsche: {
    msrpToInvoiceMarginMinPercent: 8.0,
    msrpToInvoiceMarginMaxPercent: 10.0,
    dealerHoldbackPercent: 0.0,
    holdbackBase: "MSRP",
    isDirectToConsumer: false,
  },
  Subaru: {
    msrpToInvoiceMarginMinPercent: 6.0,
    msrpToInvoiceMarginMaxPercent: 7.0,
    dealerHoldbackPercent: 2.0,
    holdbackBase: "Total MSRP",
    isDirectToConsumer: false,
  },
  Toyota: {
    msrpToInvoiceMarginMinPercent: 7.0,
    msrpToInvoiceMarginMaxPercent: 9.0,
    dealerHoldbackPercent: 2.0,
    holdbackBase: "Base MSRP",
    isDirectToConsumer: false,
  },
  Volkswagen: {
    msrpToInvoiceMarginMinPercent: 4.0,
    msrpToInvoiceMarginMaxPercent: 5.5,
    dealerHoldbackPercent: 2.0,
    holdbackBase: "MSRP",
    isDirectToConsumer: false,
  },
  Volvo: {
    msrpToInvoiceMarginMinPercent: 4.0,
    msrpToInvoiceMarginMaxPercent: 6.0,
    dealerHoldbackPercent: 1.0,
    holdbackBase: "MSRP",
    isDirectToConsumer: false,
  },
  // Direct-to-consumer brands
  Tesla: {
    msrpToInvoiceMarginMinPercent: 0.0,
    msrpToInvoiceMarginMaxPercent: 0.0,
    dealerHoldbackPercent: 0.0,
    holdbackBase: "MSRP",
    isDirectToConsumer: true,
  },
  Rivian: {
    msrpToInvoiceMarginMinPercent: 0.0,
    msrpToInvoiceMarginMaxPercent: 0.0,
    dealerHoldbackPercent: 0.0,
    holdbackBase: "MSRP",
    isDirectToConsumer: true,
  },
  Lucid: {
    msrpToInvoiceMarginMinPercent: 0.0,
    msrpToInvoiceMarginMaxPercent: 0.0,
    dealerHoldbackPercent: 0.0,
    holdbackBase: "MSRP",
    isDirectToConsumer: true,
  },
};

// Get brand margin data, or default if not found
function getBrandMarginData(make: string): BrandMarginData {
  const normalizedMake = make.charAt(0).toUpperCase() + make.slice(1).toLowerCase();
  
  if (brandMargins[normalizedMake]) {
    return brandMargins[normalizedMake];
  }
  
  // Default margins for unknown brands
  return {
    msrpToInvoiceMarginMinPercent: 5.0,
    msrpToInvoiceMarginMaxPercent: 7.0,
    dealerHoldbackPercent: 2.0,
    holdbackBase: "MSRP",
    isDirectToConsumer: false,
  };
}

// Check if body type is truck or large SUV
function isTruckOrLargeSUV(bodyType?: string): boolean {
  if (!bodyType) return false;
  const bodyTypeLower = bodyType.toLowerCase();
  return (
    bodyTypeLower.includes("truck") ||
    bodyTypeLower.includes("pickup") ||
    bodyTypeLower.includes("suv") ||
    bodyTypeLower.includes("sport utility")
  );
}

/**
 * Estimates the invoice price from MSRP based on brand and vehicle type
 */
export function estimateInvoicePrice(msrp: number, make: string, bodyType?: string): number {
  if (msrp <= 0 || !make) {
    return 0;
  }

  const brandData = getBrandMarginData(make);
  
  // Direct-to-consumer brands return MSRP
  if (brandData.isDirectToConsumer) {
    return msrp;
  }

  // Determine margin percentage based on vehicle type
  let marginPercent = (brandData.msrpToInvoiceMarginMinPercent + brandData.msrpToInvoiceMarginMaxPercent) / 2.0;
  
  // Adjust for truck/large SUV vs sedan/compact
  if (isTruckOrLargeSUV(bodyType)) {
    // Use higher end of margin range for trucks/SUVs
    marginPercent = brandData.msrpToInvoiceMarginMaxPercent;
  } else {
    // Use lower end of margin range for sedans/compacts
    marginPercent = brandData.msrpToInvoiceMarginMinPercent;
  }

  // Convert percent to decimal and calculate invoice
  const marginDecimal = marginPercent / 100.0;
  const invoice = msrp * (1.0 - marginDecimal);
  
  return invoice;
}

/**
 * Estimates the dealer holdback amount from MSRP
 */
export function estimateDealerHoldback(msrp: number, make: string): number {
  if (msrp <= 0 || !make) {
    return 0;
  }

  const brandData = getBrandMarginData(make);
  
  // Direct-to-consumer brands have no holdback
  if (brandData.isDirectToConsumer) {
    return 0;
  }

  // Convert percent to decimal and calculate holdback
  const holdbackDecimal = brandData.dealerHoldbackPercent / 100.0;
  const holdback = msrp * holdbackDecimal;
  
  return holdback;
}

/**
 * Calculates the net net price (invoice - holdback), which is the dealer break-even price
 */
export function calculateNetNetPrice(msrp: number, make: string, bodyType?: string): number {
  if (msrp <= 0 || !make) {
    return 0;
  }

  const brandData = getBrandMarginData(make);
  
  // Direct-to-consumer brands return MSRP
  if (brandData.isDirectToConsumer) {
    return msrp;
  }

  const invoice = estimateInvoicePrice(msrp, make, bodyType);
  const holdback = estimateDealerHoldback(msrp, make);
  
  return invoice - holdback;
}

/**
 * Determines the target price based on demand level
 * Returns the target price type: "MSRP", "Invoice", or "Net-Net"
 */
export function getTargetPriceType(make: string, model: string): "MSRP" | "Invoice" | "Net-Net" {
  const makeLower = make.toLowerCase();
  const modelLower = model.toLowerCase();
  
  // High Demand: Toyota Hybrids, Porsche
  if (makeLower === "toyota" && modelLower.includes("hybrid")) {
    return "MSRP";
  }
  if (makeLower === "porsche") {
    return "MSRP";
  }
  
  // Average Demand: Honda, Subaru, Mazda
  if (makeLower === "honda" || makeLower === "subaru" || makeLower === "mazda") {
    return "Invoice";
  }
  
  // Low Demand/High Stock: Jeep, Ram, Ford Trucks
  if (makeLower === "jeep" || makeLower === "ram") {
    return "Net-Net";
  }
  if (makeLower === "ford" && (modelLower.includes("f-150") || modelLower.includes("truck"))) {
    return "Net-Net";
  }
  
  // Default to Invoice for unknown combinations
  return "Invoice";
}


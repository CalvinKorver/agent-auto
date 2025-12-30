import { describe, it, expect } from 'vitest';
import { estimateInvoicePrice, estimateDealerHoldback, calculateNetNetPrice, getTargetPriceType } from './pricing';

describe('estimateInvoicePrice', () => {
  it('should calculate invoice for Toyota sedan (lower margin)', () => {
    const result = estimateInvoicePrice(40000, 'Toyota', 'Sedan');
    expect(result).toBeCloseTo(37200, 0); // 40000 * (1 - 0.07) = 37200
  });

  it('should calculate invoice for Toyota truck (higher margin)', () => {
    const result = estimateInvoicePrice(40000, 'Toyota', 'Pickup Truck');
    expect(result).toBeCloseTo(36400, 0); // 40000 * (1 - 0.09) = 36400
  });

  it('should return MSRP for Tesla (direct-to-consumer)', () => {
    const result = estimateInvoicePrice(50000, 'Tesla', 'Sedan');
    expect(result).toBe(50000);
  });

  it('should use default margins for unknown brand', () => {
    const result = estimateInvoicePrice(40000, 'UnknownBrand', 'Sedan');
    expect(result).toBeCloseTo(38000, 0); // 40000 * (1 - 0.05) = 38000
  });

  it('should return 0 for zero MSRP', () => {
    const result = estimateInvoicePrice(0, 'Toyota', 'Sedan');
    expect(result).toBe(0);
  });

  it('should return 0 for negative MSRP', () => {
    const result = estimateInvoicePrice(-1000, 'Toyota', 'Sedan');
    expect(result).toBe(0);
  });

  it('should handle Honda SUV (higher margin)', () => {
    const result = estimateInvoicePrice(35000, 'Honda', 'SUV');
    expect(result).toBeCloseTo(31925, 0); // 35000 * (1 - 0.085) = 31925
  });
});

describe('estimateDealerHoldback', () => {
  it('should calculate holdback for Acura (2%)', () => {
    const result = estimateDealerHoldback(40000, 'Acura');
    expect(result).toBe(800); // 40000 * 0.02 = 800
  });

  it('should return 0 for Audi (no holdback)', () => {
    const result = estimateDealerHoldback(50000, 'Audi');
    expect(result).toBe(0);
  });

  it('should return 0 for BMW (no holdback)', () => {
    const result = estimateDealerHoldback(60000, 'BMW');
    expect(result).toBe(0);
  });

  it('should return 0 for Porsche (no holdback)', () => {
    const result = estimateDealerHoldback(100000, 'Porsche');
    expect(result).toBe(0);
  });

  it('should return 0 for Tesla (direct-to-consumer)', () => {
    const result = estimateDealerHoldback(50000, 'Tesla');
    expect(result).toBe(0);
  });

  it('should calculate holdback for Chevrolet (3%)', () => {
    const result = estimateDealerHoldback(40000, 'Chevrolet');
    expect(result).toBe(1200); // 40000 * 0.03 = 1200
  });

  it('should use default holdback for unknown brand', () => {
    const result = estimateDealerHoldback(40000, 'UnknownBrand');
    expect(result).toBe(800); // 40000 * 0.02 = 800
  });

  it('should return 0 for zero MSRP', () => {
    const result = estimateDealerHoldback(0, 'Toyota');
    expect(result).toBe(0);
  });
});

describe('calculateNetNetPrice', () => {
  it('should calculate net net for Toyota sedan', () => {
    const result = calculateNetNetPrice(40000, 'Toyota', 'Sedan');
    // Invoice: 37200, Holdback: 800, Net-Net: 36400
    expect(result).toBeCloseTo(36400, 0);
  });

  it('should return MSRP for Tesla (direct-to-consumer)', () => {
    const result = calculateNetNetPrice(50000, 'Tesla', 'Sedan');
    expect(result).toBe(50000);
  });

  it('should calculate net net for Chevrolet truck', () => {
    const result = calculateNetNetPrice(50000, 'Chevrolet', 'Pickup Truck');
    // Invoice: 46000 (8% max), Holdback: 1500, Net-Net: 44500
    expect(result).toBeCloseTo(44500, 0);
  });

  it('should return 0 for zero MSRP', () => {
    const result = calculateNetNetPrice(0, 'Toyota', 'Sedan');
    expect(result).toBe(0);
  });
});

describe('getTargetPriceType', () => {
  it('should return MSRP for Toyota Hybrid', () => {
    const result = getTargetPriceType('Toyota', 'Camry Hybrid');
    expect(result).toBe('MSRP');
  });

  it('should return MSRP for Porsche', () => {
    const result = getTargetPriceType('Porsche', '911');
    expect(result).toBe('MSRP');
  });

  it('should return Invoice for Honda', () => {
    const result = getTargetPriceType('Honda', 'Civic');
    expect(result).toBe('Invoice');
  });

  it('should return Invoice for Subaru', () => {
    const result = getTargetPriceType('Subaru', 'Outback');
    expect(result).toBe('Invoice');
  });

  it('should return Invoice for Mazda', () => {
    const result = getTargetPriceType('Mazda', 'CX-5');
    expect(result).toBe('Invoice');
  });

  it('should return Net-Net for Jeep', () => {
    const result = getTargetPriceType('Jeep', 'Wrangler');
    expect(result).toBe('Net-Net');
  });

  it('should return Net-Net for Ram', () => {
    const result = getTargetPriceType('Ram', '1500');
    expect(result).toBe('Net-Net');
  });

  it('should return Net-Net for Ford F-150', () => {
    const result = getTargetPriceType('Ford', 'F-150');
    expect(result).toBe('Net-Net');
  });

  it('should return Invoice as default for unknown combinations', () => {
    const result = getTargetPriceType('Unknown', 'Model');
    expect(result).toBe('Invoice');
  });
});


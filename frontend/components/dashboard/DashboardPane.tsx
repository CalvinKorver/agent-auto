'use client';

import { useAuth } from '@/contexts/AuthContext';
import { TrackedOffer, Thread } from '@/lib/api';
import { IconSearch, IconUpload, IconSettings } from '@tabler/icons-react';
import { SidebarTrigger } from '@/components/ui/sidebar';
import { Separator } from '@/components/ui/separator';

interface DashboardPaneProps {
  offers: TrackedOffer[];
  threads: Thread[];
  onNavigateToThread?: (threadId: string) => void;
}

// Utility function to parse price from offer text
function parsePriceFromOfferText(offerText: string): number | null {
  // Try to find dollar amounts in the text
  const pricePatterns = [
    /\$[\d,]+(?:\.\d{2})?/g, // $50,000 or $50,000.00
    /[\d,]+(?:\.\d{2})?\s*(?:dollars?|USD)/gi, // 50,000 dollars
  ];

  for (const pattern of pricePatterns) {
    const matches = offerText.match(pattern);
    if (matches && matches.length > 0) {
      // Take the first match and extract the number
      const priceStr = matches[0].replace(/[^0-9.]/g, '');
      const price = parseFloat(priceStr);
      if (!isNaN(price) && price > 0) {
        return price;
      }
    }
  }

  return null;
}

// Calculate Lolo Score (placeholder - based on price vs market average)
function calculateLoloScore(price: number, marketAverage: number): { score: number; label: string; color: string } {
  const diff = marketAverage - price;
  const percentDiff = (diff / marketAverage) * 100;

  if (percentDiff >= 5) {
    return { score: 95, label: 'Great Deal!', color: 'text-green-500' };
  } else if (percentDiff >= 2) {
    return { score: 80, label: 'Good Deal', color: 'text-green-400' };
  } else if (percentDiff >= -2) {
    return { score: 60, label: 'Fair Price', color: 'text-yellow-500' };
  } else if (percentDiff >= -5) {
    return { score: 45, label: 'Overpriced', color: 'text-orange-500' };
  } else {
    return { score: 30, label: 'Overpriced', color: 'text-red-500' };
  }
}

// Gradient bar component for Target Fair Price
function GradientBar({ targetPrice, marketAverage }: { targetPrice: number; marketAverage: number }) {
  const minPrice = Math.min(targetPrice, marketAverage) * 0.9;
  const maxPrice = Math.max(targetPrice, marketAverage) * 1.1;
  const range = maxPrice - minPrice;
  
  const targetPosition = ((targetPrice - minPrice) / range) * 100;
  const marketPosition = ((marketAverage - minPrice) / range) * 100;

  return (
    <div className="relative w-full h-8 rounded-md overflow-hidden bg-gradient-to-r from-red-500 via-yellow-500 to-green-500">
      <div className="absolute inset-0 flex items-center">
        <div
          className="absolute w-1 h-full bg-white border-2 border-gray-900 z-10"
          style={{ left: `${targetPosition}%` }}
        />
        <div
          className="absolute w-1 h-full bg-blue-500 border-2 border-white z-10"
          style={{ left: `${marketPosition}%` }}
        />
      </div>
    </div>
  );
}

export default function DashboardPane({ offers, threads, onNavigateToThread }: DashboardPaneProps) {
  const { user } = useAuth();
  const preferences = user?.preferences;

  // Placeholder market data
  const targetFairPrice = 40800;
  const marketAverage = 41200;
  const marketVelocity = 'High Demand';
  const daysToMove = 8;

  // Find historic low from offers
  let historicLow: { price: number; dealer: string; time: string } | null = null;
  if (offers.length > 0) {
    const offerPrices = offers
      .map(offer => ({ price: parsePriceFromOfferText(offer.offerText), offer }))
      .filter(item => item.price !== null) as Array<{ price: number; offer: TrackedOffer }>;
    
    if (offerPrices.length > 0) {
      const lowest = offerPrices.reduce((min, current) => 
        current.price < min.price ? current : min
      );
      const daysAgo = Math.floor((Date.now() - new Date(lowest.offer.trackedAt).getTime()) / (1000 * 60 * 60 * 24));
      historicLow = {
        price: lowest.price,
        dealer: lowest.offer.sellerName || 'Unknown Dealer',
        time: daysAgo === 0 ? 'today' : daysAgo === 1 ? 'yesterday' : `${daysAgo} days ago`
      };
    }
  }

  // Process offers for comparison table
  const processedOffers = offers.map(offer => {
    const price = parsePriceFromOfferText(offer.offerText);
    const loloScore = price ? calculateLoloScore(price, marketAverage) : null;
    const gapToMarket = price ? price - targetFairPrice : null;

    return {
      ...offer,
      parsedPrice: price,
      loloScore,
      gapToMarket,
    };
  });

  return (
    <>
      <header className="flex h-16 shrink-0 items-center gap-2 border-b transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-16">
        <div className="flex w-full items-center gap-1 px-4 lg:gap-2 lg:px-6">
          <SidebarTrigger className="-ml-1" />
          <Separator
            orientation="vertical"
            className="mx-2 data-[orientation=vertical]:h-4"
          />
          <h1 className="text-base font-medium">Search Dashboard</h1>
          <div className="flex-1" />
          <div className="flex items-center gap-2">
            <button className="p-2 hover:bg-accent rounded-md transition-colors">
              <IconSearch className="h-5 w-5 text-muted-foreground" />
            </button>
            <button className="p-2 hover:bg-accent rounded-md transition-colors">
              <IconUpload className="h-5 w-5 text-muted-foreground" />
            </button>
            <button className="p-2 hover:bg-accent rounded-md transition-colors">
              <IconSettings className="h-5 w-5 text-muted-foreground" />
            </button>
          </div>
        </div>
      </header>
      <div className="flex flex-1 flex-col overflow-hidden">
        <div className="flex-1 overflow-y-auto bg-background px-6 py-6">
          <div className="max-w-7xl mx-auto space-y-6">
            {/* My Vehicle Search Section */}
            <div>
              <div className="bg-card border border-border rounded-lg p-6">
                <div className="flex gap-6">
                  {/* Vehicle Details */}
                  <div className="flex-1 space-y-3">
                    <div>
                      <h3 className="text-xl font-semibold text-card-foreground mb-2">
                        {preferences?.year} {preferences?.make} {preferences?.model}
                      </h3>
                    </div>
                    <div className="space-y-2 text-sm">
                      <div className="flex items-center gap-2">
                        <span className="text-muted-foreground">Trim:</span>
                        <span className="text-card-foreground">Turbo Signature</span>
                      </div>
                      <div className="flex items-center gap-2">
                        <span className="text-muted-foreground">Exterior:</span>
                        <span className="text-card-foreground">Crystal Red</span>
                      </div>
                      <div className="flex items-center gap-2">
                        <span className="text-muted-foreground">Status:</span>
                        <span className="text-card-foreground">
                          Actively searching {threads.length} {threads.length === 1 ? 'dealership' : 'dealerships'} in a 50-mile radius
                        </span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* Market Intelligence Section */}
            <div>
              <h2 className="text-2xl font-semibold text-foreground mb-4">
                Market Intelligence
              </h2>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                {/* Target Fair Price Card */}
                <div className="bg-card border border-border rounded-lg p-4">
                  <h3 className="text-sm font-medium text-card-foreground mb-3">
                    Target Fair Price
                  </h3>
                  <GradientBar targetPrice={targetFairPrice} marketAverage={marketAverage} />
                  <div className="mt-3 space-y-1">
                    <div className="text-sm">
                      <span className="text-muted-foreground">Target Fair Price: </span>
                      <span className="font-semibold text-card-foreground">
                        ${targetFairPrice.toLocaleString()}
                      </span>
                    </div>
                    <div className="text-sm">
                      <span className="text-muted-foreground">Current Market Average: </span>
                      <span className="text-card-foreground">
                        ${marketAverage.toLocaleString()}
                      </span>
                    </div>
                  </div>
                </div>

                {/* Market Velocity Card */}
                <div className="bg-card border border-border rounded-lg p-4">
                  <h3 className="text-sm font-medium text-card-foreground mb-3">
                    Market Velocity
                  </h3>
                  <div className="space-y-2">
                    <div className="text-lg font-semibold text-card-foreground">
                      {marketVelocity}
                    </div>
                    <div className="text-sm text-muted-foreground">
                      Cars move in ~{daysToMove} days
                    </div>
                  </div>
                </div>

                {/* Historic Low Card */}
                <div className="bg-card border border-border rounded-lg p-4">
                  <h3 className="text-sm font-medium text-card-foreground mb-3">
                    Historic Low
                  </h3>
                  {historicLow ? (
                    <div className="space-y-1">
                      <div className="text-lg font-semibold text-card-foreground">
                        ${historicLow.price.toLocaleString()}
                      </div>
                      <div className="text-sm text-muted-foreground">
                        at {historicLow.dealer} ({historicLow.time})
                      </div>
                    </div>
                  ) : (
                    <div className="text-sm text-muted-foreground">
                      No tracked offers yet
                    </div>
                  )}
                </div>
              </div>
            </div>

            {/* Offer Comparison Table */}
            <div>
              <h2 className="text-2xl font-semibold text-foreground mb-4">
                Offer Comparison Table
              </h2>
              <div className="bg-card border border-border rounded-lg overflow-hidden">
                <div className="overflow-x-auto">
                  <table className="w-full">
                    <thead className="bg-muted/50 border-b border-border">
                      <tr>
                        <th className="px-4 py-3 text-left text-sm font-medium text-foreground">
                          Dealer
                        </th>
                        <th className="px-4 py-3 text-left text-sm font-medium text-foreground">
                          Offer
                        </th>
                        <th className="px-4 py-3 text-left text-sm font-medium text-foreground">
                          Lolo Score
                        </th>
                        <th className="px-4 py-3 text-left text-sm font-medium text-foreground">
                          Gap to Market
                        </th>
                      </tr>
                    </thead>
                    <tbody>
                      {processedOffers.length === 0 ? (
                        <tr>
                          <td colSpan={4} className="px-4 py-8 text-center text-muted-foreground">
                            No tracked offers yet
                          </td>
                        </tr>
                      ) : (
                        processedOffers.map((offer) => (
                          <tr
                            key={offer.id}
                            onClick={() => onNavigateToThread?.(offer.threadId)}
                            className="border-b border-border hover:bg-accent/50 cursor-pointer transition-colors"
                          >
                            <td className="px-4 py-3 text-sm text-card-foreground">
                              {offer.sellerName || 'Unknown Seller'}
                            </td>
                            <td className="px-4 py-3 text-sm text-card-foreground">
                              {offer.parsedPrice ? (
                                `$${offer.parsedPrice.toLocaleString()}`
                              ) : (
                                <span className="text-muted-foreground italic">
                                  {offer.offerText.substring(0, 50)}...
                                </span>
                              )}
                            </td>
                            <td className="px-4 py-3 text-sm">
                              {offer.loloScore ? (
                                <div className="flex items-center gap-2">
                                  <div
                                    className={`w-10 h-10 rounded-full flex items-center justify-center font-semibold ${
                                      offer.loloScore.score >= 80
                                        ? 'bg-green-500/20 text-green-500'
                                        : offer.loloScore.score >= 60
                                        ? 'bg-yellow-500/20 text-yellow-500'
                                        : 'bg-red-500/20 text-red-500'
                                    }`}
                                  >
                                    {offer.loloScore.score}
                                  </div>
                                  <span className={offer.loloScore.color}>
                                    ({offer.loloScore.label})
                                  </span>
                                </div>
                              ) : (
                                <span className="text-muted-foreground">—</span>
                              )}
                            </td>
                            <td className="px-4 py-3 text-sm">
                              {offer.gapToMarket !== null ? (
                                <span
                                  className={
                                    offer.gapToMarket < 0
                                      ? 'text-green-500'
                                      : offer.gapToMarket > 1000
                                      ? 'text-red-500'
                                      : 'text-yellow-500'
                                  }
                                >
                                  {offer.gapToMarket < 0 ? '' : '+'}
                                  ${Math.abs(offer.gapToMarket).toLocaleString()}{' '}
                                  {offer.gapToMarket < 0 ? '(Below Target)' : '(Above Target)'}
                                </span>
                              ) : (
                                <span className="text-muted-foreground">—</span>
                              )}
                            </td>
                          </tr>
                        ))
                      )}
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </>
  );
}


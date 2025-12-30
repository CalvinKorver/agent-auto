'use client';

import { useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { TrackedOffer, Thread, offerAPI } from '@/lib/api';
import { IconSearch, IconUpload, IconSettings, IconTrash } from '@tabler/icons-react';
import { SidebarTrigger } from '@/components/ui/sidebar';
import { Separator } from '@/components/ui/separator';
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip';
import { toast } from 'sonner';
import { estimateInvoicePrice, estimateDealerHoldback, calculateNetNetPrice, getTargetPriceType } from '@/lib/pricing';

interface DashboardPaneProps {
  offers: TrackedOffer[];
  threads: Thread[];
  onNavigateToThread?: (threadId: string) => void;
  onOfferDeleted?: () => void;
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

// Calculate Otto Score (placeholder - based on price vs market average)
function calculateOttoScore(price: number, marketAverage: number): { score: number; label: string; color: string } {
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


export default function DashboardPane({ offers, threads, onNavigateToThread, onOfferDeleted }: DashboardPaneProps) {
  const [deletingOfferId, setDeletingOfferId] = useState<string | null>(null);

  const handleDeleteOffer = async (offerId: string, e: React.MouseEvent) => {
    e.stopPropagation(); // Prevent row click navigation
    
    if (!confirm('Are you sure you want to delete this offer?')) {
      return;
    }

    setDeletingOfferId(offerId);
    try {
      await offerAPI.deleteOffer(offerId);
      toast.success('Offer deleted successfully');
      if (onOfferDeleted) {
        onOfferDeleted();
      }
    } catch (error) {
      console.error('Failed to delete offer:', error);
      toast.error('Failed to delete offer. Please try again.');
    } finally {
      setDeletingOfferId(null);
    }
  };
  const { user } = useAuth();
  const preferences = user?.preferences;

  // Calculate pricing data from preferences
  const msrp = preferences?.baseMsrp || 0;
  const make = preferences?.make || '';
  const model = preferences?.model || '';
  const estimatedInvoice = msrp > 0 && make ? estimateInvoicePrice(msrp, make) : 0;
  const estimatedHoldback = msrp > 0 && make ? estimateDealerHoldback(msrp, make) : 0;
  const netNetPrice = msrp > 0 && make ? calculateNetNetPrice(msrp, make) : 0;
  const targetPriceType = make && model ? getTargetPriceType(make, model) : 'Invoice';

  // Placeholder market data
  const marketAverage = msrp > 0 ? msrp * 1.01 : 41200; // 1% above MSRP as placeholder

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
    const ottoScore = price ? calculateOttoScore(price, marketAverage) : null;
    const gapToMarket = price ? price - msrp : null;

    return {
      ...offer,
      parsedPrice: price,
      ottoScore,
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
                        <span className="text-card-foreground">{preferences?.trim || 'Unspecified'}</span>
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
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                {/* MSRP (Suggested retail) Card */}
                <div className="bg-card border border-border rounded-lg p-4">
                  <h3 className="text-sm font-medium text-card-foreground mb-3">
                    MSRP (Suggested retail)
                  </h3>
                  <div className="space-y-2">
                    <div className="text-2xl font-semibold text-card-foreground">
                      {msrp > 0 ? `$${msrp.toLocaleString(undefined, { maximumFractionDigits: 0 })}` : 'N/A'}
                    </div>
                    <div className="text-xs text-muted-foreground">
                      Manufacturer&apos;s suggested retail price
                    </div>
                  </div>
                </div>

                {/* Estimated Invoice Price Card */}
                <div className="bg-card border border-border rounded-lg p-4">
                  <h3 className="text-sm font-medium text-card-foreground mb-3">
                    Estimated Invoice Price
                  </h3>
                  <div className="space-y-2">
                    <div className="text-2xl font-semibold text-card-foreground">
                      {estimatedInvoice > 0 ? `$${estimatedInvoice.toLocaleString(undefined, { maximumFractionDigits: 0 })}` : 'N/A'}
                    </div>
                    <div className="text-xs text-muted-foreground">
                      Dealer cost before holdback
                    </div>
                  </div>
                </div>

                {/* Estimated Dealer Holdback Card */}
                <div className="bg-card border border-border rounded-lg p-4">
                  <h3 className="text-sm font-medium text-card-foreground mb-3">
                    Estimated Dealer Holdback
                  </h3>
                  <div className="space-y-2">
                    <div className="text-2xl font-semibold text-card-foreground">
                      {estimatedHoldback > 0 ? `$${estimatedHoldback.toLocaleString(undefined, { maximumFractionDigits: 0 })}` : 'N/A'}
                    </div>
                    <div className="text-xs text-muted-foreground">
                      Manufacturer rebate to dealer
                    </div>
                  </div>
                </div>

                {/* Net Net Price Card */}
                <div className="bg-card border border-border rounded-lg p-4">
                  <h3 className="text-sm font-medium text-card-foreground mb-3">
                    <Tooltip>
                      <TooltipTrigger asChild>
                        <span className="cursor-help underline decoration-dotted">
                          Net Net Price
                        </span>
                      </TooltipTrigger>
                      <TooltipContent>
                        <p>This is the price that dealers &apos;break even&apos; on after the holdbacks and invoice pricing.</p>
                      </TooltipContent>
                    </Tooltip>
                  </h3>
                  <div className="space-y-2">
                    <div className="text-2xl font-semibold text-card-foreground">
                      {netNetPrice > 0 ? `$${netNetPrice.toLocaleString(undefined, { maximumFractionDigits: 0 })}` : 'N/A'}
                    </div>
                    <div className="text-xs text-muted-foreground">
                      Dealer break-even price
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* What a Good Deal Is Section */}
            <div>
              <h2 className="text-2xl font-semibold text-foreground mb-4">
                What a Good Deal Is
              </h2>
              <div className="bg-card border border-border rounded-lg p-6">
                <div className="overflow-x-auto">
                  <table className="w-full">
                    <thead className="bg-muted/50 border-b border-border">
                      <tr>
                        <th className="px-4 py-3 text-left text-sm font-medium text-foreground">
                          If the car is...
                        </th>
                        <th className="px-4 py-3 text-left text-sm font-medium text-foreground">
                          Your Target Price
                        </th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr className="border-b border-border">
                        <td className="px-4 py-3 text-sm text-card-foreground">
                          <span className="font-medium">High Demand</span> (Toyota Hybrids, Porsche)
                        </td>
                        <td className="px-4 py-3 text-sm text-card-foreground">
                          MSRP is often the floor.
                        </td>
                      </tr>
                      <tr className="border-b border-border">
                        <td className="px-4 py-3 text-sm text-card-foreground">
                          <span className="font-medium">Average Demand</span> (Honda, Subaru, Mazda)
                        </td>
                        <td className="px-4 py-3 text-sm text-card-foreground">
                          Invoice is a very fair deal.
                        </td>
                      </tr>
                      <tr>
                        <td className="px-4 py-3 text-sm text-card-foreground">
                          <span className="font-medium">Low Demand/High Stock</span> (Jeep, Ram, Ford Trucks)
                        </td>
                        <td className="px-4 py-3 text-sm text-card-foreground">
                          Net-Net or below (using hidden incentives).
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </div>
                {make && model && (
                  <div className="mt-4 p-4 bg-muted/30 rounded-lg">
                    <div className="text-sm">
                      <span className="text-muted-foreground">For your </span>
                      <span className="font-medium text-card-foreground">{make} {model}</span>
                      <span className="text-muted-foreground">, target: </span>
                      <span className="font-semibold text-card-foreground">
                        {targetPriceType === 'MSRP' && `$${msrp.toLocaleString(undefined, { maximumFractionDigits: 0 })} (MSRP)`}
                        {targetPriceType === 'Invoice' && `$${estimatedInvoice.toLocaleString(undefined, { maximumFractionDigits: 0 })} (Invoice)`}
                        {targetPriceType === 'Net-Net' && `$${netNetPrice.toLocaleString(undefined, { maximumFractionDigits: 0 })} (Net-Net)`}
                      </span>
                    </div>
                  </div>
                )}
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
                          Otto Score
                        </th>
                        <th className="px-4 py-3 text-left text-sm font-medium text-foreground">
                          Gap to Market
                        </th>
                        <th className="px-4 py-3 text-left text-sm font-medium text-foreground w-16">
                          Actions
                        </th>
                      </tr>
                    </thead>
                    <tbody>
                      {processedOffers.length === 0 ? (
                        <tr>
                          <td colSpan={5} className="px-4 py-8 text-center text-muted-foreground">
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
                              {offer.ottoScore ? (
                                <div className="flex items-center gap-2">
                                  <div
                                    className={`w-10 h-10 rounded-full flex items-center justify-center font-semibold ${
                                      offer.ottoScore.score >= 80
                                        ? 'bg-green-500/20 text-green-500'
                                        : offer.ottoScore.score >= 60
                                        ? 'bg-yellow-500/20 text-yellow-500'
                                        : 'bg-red-500/20 text-red-500'
                                    }`}
                                  >
                                    {offer.ottoScore.score}
                                  </div>
                                  <span className={offer.ottoScore.color}>
                                    ({offer.ottoScore.label})
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
                                  {offer.gapToMarket < 0 ? '(Below MSRP)' : '(Above MSRP)'}
                                </span>
                              ) : (
                                <span className="text-muted-foreground">—</span>
                              )}
                            </td>
                            <td className="px-4 py-3 text-sm">
                              <button
                                onClick={(e) => handleDeleteOffer(offer.id, e)}
                                disabled={deletingOfferId === offer.id}
                                className="p-2 hover:bg-destructive/10 rounded-md transition-colors disabled:opacity-50 disabled:cursor-not-allowed text-muted-foreground hover:text-destructive"
                                title="Delete offer"
                              >
                                <IconTrash className="h-4 w-4" />
                              </button>
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


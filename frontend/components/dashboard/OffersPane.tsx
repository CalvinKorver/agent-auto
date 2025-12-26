'use client';

import { useState, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { TrackedOffer, offerAPI } from '@/lib/api';

interface OffersPaneProps {
  onNavigateToThread?: (threadId: string) => void;
}

export default function OffersPane({ onNavigateToThread }: OffersPaneProps) {
  const [offers, setOffers] = useState<TrackedOffer[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadOffers();
  }, []);

  const loadOffers = async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await offerAPI.getAllOffers();
      setOffers(data);
    } catch (err) {
      console.error('Failed to load offers:', err);
      setError('Failed to load offers. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const handleOfferClick = (threadId: string) => {
    if (onNavigateToThread) {
      onNavigateToThread(threadId);
    }
  };

  if (loading) {
    return (
      <div className="flex-1 flex items-center justify-center bg-background">
        <div className="text-center">
          <div className="text-muted-foreground">Loading offers...</div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex-1 flex items-center justify-center bg-background">
        <div className="text-center max-w-md">
          <div className="text-red-500 mb-4">{error}</div>
          <Button onClick={loadOffers}>
            Retry
          </Button>
        </div>
      </div>
    );
  }

  if (offers.length === 0) {
    return (
      <div className="flex-1 flex items-center justify-center bg-background">
        <div className="text-center max-w-md">
          <div className="mb-6 flex justify-center">
            <svg
              className="w-24 h-24 text-muted"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={1.5}
                d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
          </div>
          <h2 className="text-2xl font-semibold text-foreground mb-2">
            No Tracked Offers
          </h2>
          <p className="text-muted-foreground">
            Start tracking offers from seller messages to use them as leverage in negotiations.
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex-1 flex flex-col bg-background h-full">
      {/* Header */}
      <div className="border-b border-border px-6 py-4 bg-card">
        <h2 className="text-xl font-semibold text-card-foreground">Tracked Offers</h2>
        <p className="text-sm text-muted-foreground mt-1">
          {offers.length} {offers.length === 1 ? 'offer' : 'offers'} tracked across all dealers
        </p>
      </div>

      {/* Offers List */}
      <div className="flex-1 overflow-y-auto px-6 py-4">
        <div className="space-y-4">
          {offers.map((offer) => (
            <div
              key={offer.id}
              onClick={() => handleOfferClick(offer.threadId)}
              className="bg-card border border-border rounded-lg p-4 hover:bg-accent hover:border-accent-foreground transition-colors cursor-pointer"
            >
              <div className="flex items-start justify-between mb-2">
                <div>
                  <h3 className="font-semibold text-foreground">
                    {offer.sellerName || 'Unknown Seller'}
                  </h3>
                  {offer.threadType && (
                    <span className="text-xs text-muted-foreground capitalize">
                      {offer.threadType}
                    </span>
                  )}
                </div>
                <div className="text-xs text-muted-foreground">
                  {new Date(offer.trackedAt).toLocaleDateString()}
                </div>
              </div>
              <p className="text-card-foreground whitespace-pre-wrap break-words">
                {offer.offerText}
              </p>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

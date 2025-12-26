'use client';

import { useState } from 'react';
import { offerAPI } from '@/lib/api';
import { Button } from '@/components/ui/button';

interface TrackOfferButtonProps {
  threadId: string;
  messageId: string;
  messageContent: string;
}

export default function TrackOfferButton({ threadId, messageId, messageContent }: TrackOfferButtonProps) {
  const [showDialog, setShowDialog] = useState(false);
  const [offerText, setOfferText] = useState(messageContent);
  const [tracking, setTracking] = useState(false);

  const handleTrackOffer = async () => {
    if (!offerText.trim()) {
      return;
    }

    setTracking(true);
    try {
      await offerAPI.createOffer(threadId, offerText.trim(), messageId);
      setShowDialog(false);
      alert('Offer tracked successfully!');
    } catch (error) {
      console.error('Failed to track offer:', error);
      alert('Failed to track offer. Please try again.');
    } finally {
      setTracking(false);
    }
  };

  return (
    <>
      <Button
        onClick={() => setShowDialog(true)}
        variant="default"
        size="sm"
        className="mt-1.5 mb-1.5"
      >
        Track Offer
      </Button>

      {showDialog && (
        <div
          className="fixed inset-0 bg-black/50 flex items-center justify-center z-50"
          onClick={() => !tracking && setShowDialog(false)}
        >
          <div
            className="bg-popover rounded-lg p-6 max-w-md w-full mx-4 border border-border"
            onClick={(e) => e.stopPropagation()}
          >
            <h2 className="text-xl font-bold text-popover-foreground mb-4">Track Offer</h2>
            <p className="text-sm text-muted-foreground mb-4">
              Save this offer to use as leverage when negotiating with other sellers.
            </p>

            <div className="mb-4">
              <label className="block text-sm font-medium mb-2 text-foreground">
                Offer Details
              </label>
              <textarea
                value={offerText}
                onChange={(e) => setOfferText(e.target.value)}
                rows={4}
                className="w-full px-3 py-2 bg-background border border-input rounded-md resize-none focus:outline-none focus:ring-2 focus:ring-ring text-foreground"
                placeholder="e.g., $25,000 OTD with free floor mats"
              />
            </div>

            <div className="flex gap-3">
              <Button
                onClick={handleTrackOffer}
                disabled={tracking || !offerText.trim()}
                className="flex-1"
              >
                {tracking ? 'Tracking...' : 'Track Offer'}
              </Button>
              <Button
                onClick={() => setShowDialog(false)}
                disabled={tracking}
                variant="secondary"
                className="flex-1"
              >
                Cancel
              </Button>
            </div>
          </div>
        </div>
      )}
    </>
  );
}

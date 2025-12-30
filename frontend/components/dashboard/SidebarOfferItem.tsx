'use client';

import { TrackedOffer } from '@/lib/api';

interface SidebarOfferItemProps {
  offer: TrackedOffer;
  isActive: boolean;
  onSelect: (offer: TrackedOffer) => void;
}

export function SidebarOfferItem({
  offer,
  isActive,
  onSelect,
}: SidebarOfferItemProps) {
  return (
    <div
      className={`cursor-pointer rounded-md p-1.5 hover:bg-accent transition-colors ${
        isActive ? 'bg-accent' : ''
      }`}
      onClick={() => onSelect(offer)}
    >
      <div className="font-medium text-sm truncate leading-tight">
        {offer.sellerName || 'Unknown Seller'}
      </div>
      <div className="text-xs text-muted-foreground truncate leading-tight">
        {offer.offerText}
      </div>
    </div>
  );
}


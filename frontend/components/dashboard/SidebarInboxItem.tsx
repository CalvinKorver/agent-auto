'use client';

import { InboxMessage } from '@/lib/api';

interface SidebarInboxItemProps {
  message: InboxMessage;
  isActive: boolean;
  onSelect: (message: InboxMessage) => void;
}

export function SidebarInboxItem({
  message,
  isActive,
  onSelect,
}: SidebarInboxItemProps) {

  return (
    <div
      className={`cursor-pointer rounded-md p-2 hover:bg-accent transition-colors ${
        isActive ? 'bg-accent' : ''
      }`}
      onClick={() => onSelect(message)}
    >
      <div className="font-medium text-sm truncate">
        {message.subject || 'No Subject'}
      </div>
      <div className="text-xs text-muted-foreground truncate">
        {message.senderEmail}
      </div>
    </div>
  );
}

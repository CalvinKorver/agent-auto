import { ReactNode } from 'react';

interface ChatInterfaceProps {
  children: ReactNode;
}

export default function ChatInterface({ children }: ChatInterfaceProps) {
  return (
    <div className="flex-1 bg-muted/30 overflow-y-auto">
      <div className="max-w-3xl mx-auto px-4 py-8 lg:px-8 lg:py-12">
        {children}
      </div>
    </div>
  );
}

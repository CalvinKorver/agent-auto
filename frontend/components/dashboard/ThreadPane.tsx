'use client';

import { useState, useEffect } from 'react';
import { Thread, threadAPI, InboxMessage, messageAPI } from '@/lib/api';
import { useAuth } from '@/contexts/AuthContext';

interface ThreadPaneProps {
  threads: Thread[];
  selectedThreadId: string | null;
  selectedInboxMessageId: string | null;
  onThreadSelect: (threadId: string) => void;
  onThreadCreated: (thread: Thread) => void;
  onInboxMessageSelect: (message: InboxMessage) => void;
}

export default function ThreadPane({
  threads,
  selectedThreadId,
  selectedInboxMessageId,
  onThreadSelect,
  onThreadCreated,
  onInboxMessageSelect,
}: ThreadPaneProps) {
  const { user } = useAuth();
  const [isCreatingThread, setIsCreatingThread] = useState(false);
  const [newSellerName, setNewSellerName] = useState('');
  const [newSellerType, setNewSellerType] = useState<'private' | 'dealership' | 'other'>('dealership');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [inboxMessages, setInboxMessages] = useState<InboxMessage[]>([]);
  const [loadingInbox, setLoadingInbox] = useState(false);
  const [copied, setCopied] = useState(false);
  const [showEmailDialog, setShowEmailDialog] = useState(false);
  const [refreshKey, setRefreshKey] = useState(0);

  useEffect(() => {
    loadInboxMessages();
  }, [refreshKey]);

  // Expose refresh method via custom event
  useEffect(() => {
    const handleRefresh = () => {
      setRefreshKey(prev => prev + 1);
    };
    window.addEventListener('refreshInboxMessages', handleRefresh);
    return () => window.removeEventListener('refreshInboxMessages', handleRefresh);
  }, []);

  const loadInboxMessages = async () => {
    setLoadingInbox(true);
    try {
      const response = await messageAPI.getInboxMessages();
      setInboxMessages(response.messages);
    } catch (error) {
      console.error('Failed to load inbox messages:', error);
    } finally {
      setLoadingInbox(false);
    }
  };

  const handleCopyEmail = async () => {
    if (user?.inboxEmail) {
      try {
        await navigator.clipboard.writeText(user.inboxEmail);
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
      } catch (err) {
        console.error('Failed to copy email:', err);
      }
    }
  };

  const handleCreateThread = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!newSellerName.trim()) {
      setError('Seller name is required');
      return;
    }

    setLoading(true);
    setError('');

    try {
      const newThread = await threadAPI.create({
        sellerName: newSellerName.trim(),
        sellerType: newSellerType,
      });

      onThreadCreated(newThread);
      setNewSellerName('');
      setNewSellerType('dealership');
      setIsCreatingThread(false);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to create thread');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="w-80 bg-sidebar text-sidebar-foreground flex flex-col h-screen border-r border-sidebar-border">
      {/* Header */}
      <div className="p-6 border-b border-sidebar-border">
        <div className="flex items-center gap-3">
          <svg className="w-8 h-8 text-muted-foreground" fill="currentColor" viewBox="0 0 24 24">
            <path d="M18.92 6.01C18.72 5.42 18.16 5 17.5 5h-11c-.66 0-1.21.42-1.42 1.01L3 12v8c0 .55.45 1 1 1h1c.55 0 1-.45 1-1v-1h12v1c0 .55.45 1 1 1h1c.55 0 1-.45 1-1v-8l-2.08-5.99zM6.5 16c-.83 0-1.5-.67-1.5-1.5S5.67 13 6.5 13s1.5.67 1.5 1.5S7.33 16 6.5 16zm11 0c-.83 0-1.5-.67-1.5-1.5s.67-1.5 1.5-1.5 1.5.67 1.5 1.5-.67 1.5-1.5 1.5zM5 11l1.5-4.5h11L19 11H5z" />
          </svg>
          <h1 className="text-xl font-bold">Lolo AI</h1>
        </div>
      </div>

      {/* Inbox Section */}
      <div className="border-b border-sidebar-border">
        <div className="p-4">
          <div className="text-xs text-muted-foreground uppercase tracking-wider mb-2">
            Inbox
          </div>

          {/* Add Emails Button */}
          {user?.inboxEmail && (
            <button
              onClick={() => setShowEmailDialog(true)}
              className="w-full bg-primary hover:bg-primary/90 text-primary-foreground text-xs font-medium py-2 px-3 rounded-md mb-3 transition-colors cursor-pointer"
            >
              Add Emails
            </button>
          )}

          {loadingInbox ? (
            <div className="text-muted-foreground text-sm">Loading...</div>
          ) : inboxMessages.length === 0 ? (
            <div className="text-muted-foreground text-sm italic">
              No new messages
            </div>
          ) : (
            <div className="space-y-2">
              {inboxMessages.map((message) => (
                <div
                  key={message.id}
                  onClick={() => onInboxMessageSelect(message)}
                  className={`rounded-md px-3 py-2 text-sm hover:bg-accent transition-colors cursor-pointer ${
                    selectedInboxMessageId === message.id
                      ? 'bg-accent text-accent-foreground'
                      : 'bg-accent/60'
                  }`}
                >
                  <div className="font-medium truncate">{message.subject || 'No Subject'}</div>
                  <div className="text-xs text-muted-foreground mt-0.5 truncate">{message.senderEmail}</div>
                  <div className="text-xs text-muted-foreground mt-1 line-clamp-2">{message.content}</div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>

      {/* Threads Section */}
      <div className="flex-1 overflow-y-auto">
        <div className="p-4">
          <div className="text-xs text-muted-foreground uppercase tracking-wider mb-3">
            Threads
          </div>

          {threads.length === 0 ? (
            <div className="text-muted-foreground text-sm italic">
              (No active negotiations)
            </div>
          ) : (
            <div className="space-y-2">
              {threads.map((thread) => (
                <button
                  key={thread.id}
                  onClick={() => onThreadSelect(thread.id)}
                  className={`w-full text-left px-3 py-2 rounded-md transition-colors ${
                    selectedThreadId === thread.id
                      ? 'bg-accent text-accent-foreground'
                      : 'text-foreground hover:bg-accent/70'
                  }`}
                >
                  <div className="font-medium text-sm">{thread.sellerName}</div>
                </button>
              ))}
            </div>
          )}
        </div>
      </div>

      {/* New Thread Button/Form */}
      <div className="p-4 border-t border-sidebar-border">
        {!isCreatingThread ? (
          <button
            onClick={() => setIsCreatingThread(true)}
            className="w-full bg-primary hover:bg-primary/90 text-primary-foreground font-medium py-3 px-4 rounded-md flex items-center justify-center gap-2 transition-colors cursor-pointer"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
            </svg>
            New Seller Thread
          </button>
        ) : (
          <form onSubmit={handleCreateThread} className="space-y-3">
            <div>
              <select
                value={newSellerType}
                onChange={(e) => setNewSellerType(e.target.value as any)}
                className="w-full bg-accent text-accent-foreground px-3 py-2 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-ring"
                disabled={loading}
              >
                <option value="dealership">Dealership</option>
                <option value="private">Private Seller</option>
                <option value="other">Other</option>
              </select>
            </div>

            <div>
              <input
                type="text"
                value={newSellerName}
                onChange={(e) => setNewSellerName(e.target.value)}
                placeholder="Seller name (e.g., Subaru of Renton)"
                className="w-full bg-accent text-accent-foreground px-3 py-2 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-ring"
                disabled={loading}
              />
            </div>

            {error && (
              <div className="text-red-400 text-xs">{error}</div>
            )}

            <div className="flex gap-2">
              <button
                type="submit"
                disabled={loading}
                className="flex-1 bg-primary hover:bg-primary/90 text-primary-foreground font-medium py-2 px-3 rounded-md text-sm transition-colors disabled:opacity-50"
              >
                {loading ? 'Creating...' : 'Create'}
              </button>
              <button
                type="button"
                onClick={() => {
                  setIsCreatingThread(false);
                  setNewSellerName('');
                  setError('');
                }}
                disabled={loading}
                className="flex-1 bg-accent hover:bg-accent/80 text-accent-foreground font-medium py-2 px-3 rounded-md text-sm transition-colors disabled:opacity-50"
              >
                Cancel
              </button>
            </div>
          </form>
        )}
      </div>

      {/* Email Dialog */}
      {showEmailDialog && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" onClick={() => setShowEmailDialog(false)}>
          <div className="bg-popover rounded-lg p-6 max-w-md w-full mx-4 border border-border" onClick={(e) => e.stopPropagation()}>
            <h2 className="text-xl font-bold text-popover-foreground mb-4">Forward Emails Here</h2>
            <p className="text-sm text-muted-foreground mb-4">
              Forward or BCC emails from sellers to this address and they&apos;ll appear in your inbox:
            </p>
            <div className="bg-accent rounded-md p-3 mb-4">
              <div className="flex items-center justify-between">
                <code className="text-sm text-primary break-all">{user?.inboxEmail}</code>
                <button
                  onClick={handleCopyEmail}
                  className="ml-3 flex-shrink-0 px-3 py-1.5 bg-primary hover:bg-primary/90 text-primary-foreground text-xs rounded transition-colors cursor-pointer"
                >
                  {copied ? 'Copied!' : 'Copy'}
                </button>
              </div>
            </div>
            <button
              onClick={() => setShowEmailDialog(false)}
              className="w-full bg-accent hover:bg-accent/80 text-accent-foreground font-medium py-2 px-4 rounded-md transition-colors cursor-pointer"
            >
              Close
            </button>
          </div>
        </div>
      )}
    </div>
  );
}

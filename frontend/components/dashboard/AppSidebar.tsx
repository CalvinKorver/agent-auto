'use client';

import { useState, useEffect } from 'react';
import { IconMail, IconMessageCircle, IconPlus, IconCopy, IconCheck } from '@tabler/icons-react';
import { Thread, InboxMessage, messageAPI, threadAPI } from '@/lib/api';
import { useAuth } from '@/contexts/AuthContext';
import { SidebarInboxItem } from './SidebarInboxItem';
import { NavUser } from './NavUser';
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarFooter,
  SidebarHeader,
} from '@/components/ui/sidebar';
import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import * as React from 'react';

interface AppSidebarProps extends React.ComponentProps<typeof Sidebar> {
  threads: Thread[];
  selectedThreadId: string | null;
  selectedInboxMessageId: string | null;
  onThreadSelect: (threadId: string) => void;
  onThreadCreated: (thread: Thread) => void;
  onInboxMessageSelect: (message: InboxMessage) => void;
}

export function AppSidebar({
  threads,
  selectedThreadId,
  selectedInboxMessageId,
  onThreadSelect,
  onThreadCreated,
  onInboxMessageSelect,
  ...props
}: AppSidebarProps) {
  const { user } = useAuth();
  const [inboxMessages, setInboxMessages] = useState<InboxMessage[]>([]);
  const [loadingInbox, setLoadingInbox] = useState(false);
  const [showEmailDialog, setShowEmailDialog] = useState(false);
  const [showNewThreadDialog, setShowNewThreadDialog] = useState(false);
  const [copied, setCopied] = useState(false);
  const [refreshKey, setRefreshKey] = useState(0);

  // New thread form state
  const [newSellerName, setNewSellerName] = useState('');
  const [newSellerType, setNewSellerType] = useState<'private' | 'dealership' | 'other'>('dealership');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    loadInboxMessages();
  }, [refreshKey]);

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
      setShowNewThreadDialog(false);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to create thread');
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <Sidebar {...props}>

        <SidebarHeader className="border-sidebar-border h-16 border-b">
          <NavUser user={{
            name: user?.email?.split('@')[0] || 'User',
            email: user?.email || ''
          }} />
        </SidebarHeader>
        <SidebarContent>

          {/* Inbox Section */}
          <SidebarGroup>
            <SidebarGroupLabel className="flex items-center justify-between">
              <span className="flex items-center gap-2">
                <IconMail className="h-4 w-4" />
                Inbox
              </span>
            </SidebarGroupLabel>
            <SidebarGroupContent>
              {user?.inboxEmail && (
                <div className="px-2 pb-2">
                  <Button
                    onClick={() => setShowEmailDialog(true)}
                    variant="outline"
                    size="sm"
                    className="w-full"
                  >
                    <IconPlus className="h-4 w-4 mr-2" />
                    Add Emails
                  </Button>
                </div>
              )}

              <div className="flex flex-col gap-2 px-2">
                {loadingInbox ? (
                  <div className="px-2 py-2 text-sm text-muted-foreground">Loading...</div>
                ) : (
                  inboxMessages.map((message) => (
                    <SidebarInboxItem
                      key={message.id}
                      message={message}
                      isActive={selectedInboxMessageId === message.id}
                      onSelect={onInboxMessageSelect}
                    />
                  ))
                )}
              </div>
            </SidebarGroupContent>
          </SidebarGroup>

          {/* Threads Section */}
          <SidebarGroup>
            <SidebarGroupLabel className="flex items-center gap-2">
              <IconMessageCircle className="h-4 w-4" />
              Threads
            </SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {threads.length === 0 ? (
                  <div className="px-4 py-2 text-sm text-muted-foreground italic">
                    No active negotiations
                  </div>
                ) : (
                  threads.map((thread) => (
                    <SidebarMenuItem key={thread.id}>
                      <SidebarMenuButton
                        onClick={() => onThreadSelect(thread.id)}
                        isActive={selectedThreadId === thread.id}
                        className="flex flex-col items-start h-auto py-2"
                      >
                        <div className="font-medium">{thread.sellerName}</div>
                      </SidebarMenuButton>
                    </SidebarMenuItem>
                  ))
                )}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>

        <SidebarFooter className="border-t p-4">
          <Button
            onClick={() => setShowNewThreadDialog(true)}
            className="w-full"
          >
            <IconPlus className="h-4 w-4 mr-2" />
            New Seller Thread
          </Button>
        </SidebarFooter>
      </Sidebar>

      {/* Email Dialog */}
      <Dialog open={showEmailDialog} onOpenChange={setShowEmailDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Forward Emails Here</DialogTitle>
            <DialogDescription>
              Forward or BCC emails from sellers to this address and they'll appear in your inbox:
            </DialogDescription>
          </DialogHeader>
          <div className="flex items-center gap-2 p-3 bg-muted rounded-md">
            <code className="flex-1 text-sm break-all">{user?.inboxEmail}</code>
            <Button
              onClick={handleCopyEmail}
              variant="outline"
              size="sm"
            >
              {copied ? (
                <>
                  <IconCheck className="h-4 w-4 mr-2" />
                  Copied
                </>
              ) : (
                <>
                  <IconCopy className="h-4 w-4 mr-2" />
                  Copy
                </>
              )}
            </Button>
          </div>
        </DialogContent>
      </Dialog>

      {/* New Thread Dialog */}
      <Dialog open={showNewThreadDialog} onOpenChange={setShowNewThreadDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Create New Seller Thread</DialogTitle>
            <DialogDescription>
              Start a thread to organize and track your communication with this seller.
            </DialogDescription>
          </DialogHeader>
          <form onSubmit={handleCreateThread} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="sellerType">Seller Type</Label>
              <Select
                value={newSellerType}
                onValueChange={(value) => setNewSellerType(value as any)}
                disabled={loading}
              >
                <SelectTrigger id="sellerType">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="dealership">Dealership</SelectItem>
                  <SelectItem value="private">Private Seller</SelectItem>
                  <SelectItem value="other">Other</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2">
              <Label htmlFor="sellerName">Seller Name</Label>
              <Input
                id="sellerName"
                value={newSellerName}
                onChange={(e) => setNewSellerName(e.target.value)}
                placeholder="e.g., Subaru of Renton"
                disabled={loading}
              />
            </div>

            {error && (
              <div className="text-sm text-destructive">{error}</div>
            )}

            <div className="flex gap-2">
              <Button type="submit" disabled={loading} className="flex-1">
                {loading ? 'Creating...' : 'Create'}
              </Button>
              <Button
                type="button"
                variant="outline"
                onClick={() => {
                  setShowNewThreadDialog(false);
                  setNewSellerName('');
                  setError('');
                }}
                disabled={loading}
                className="flex-1"
              >
                Cancel
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>
    </>
  );
}

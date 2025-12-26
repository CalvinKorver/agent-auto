'use client';

import { useEffect, useState, Suspense } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';
import { AppSidebar } from '@/components/dashboard/AppSidebar';
import ChatPane from '@/components/dashboard/ChatPane';
import TopNavBar from '@/components/dashboard/TopNavBar';
import OffersPane from '@/components/dashboard/OffersPane';
import { Thread, threadAPI, InboxMessage } from '@/lib/api';
import { SidebarProvider, SidebarInset } from '@/components/ui/sidebar';
import { toast } from 'sonner';

type ViewMode = 'chat' | 'inbox' | 'offers';

function DashboardContent() {
  const { user, loading, refreshUser } = useAuth();
  const router = useRouter();
  const searchParams = useSearchParams();
  const [threads, setThreads] = useState<Thread[]>([]);
  const [selectedThreadId, setSelectedThreadId] = useState<string | null>(null);
  const [selectedInboxMessage, setSelectedInboxMessage] = useState<InboxMessage | null>(null);
  const [loadingThreads, setLoadingThreads] = useState(false);
  const [viewMode, setViewMode] = useState<ViewMode>('chat');

  // Check if user just connected Gmail
  useEffect(() => {
    const gmailConnected = searchParams.get('gmail_connected');
    if (gmailConnected === 'true') {
      // Refresh user data to get updated Gmail connection status
      refreshUser();
      // Show success toast
      toast.success('Gmail connected successfully!');
      // Clean up URL
      router.replace('/dashboard');
    }
  }, [searchParams, router, refreshUser]);

  useEffect(() => {
    if (!loading) {
      if (!user) {
        router.push('/login');
      } else if (!user.preferences) {
        router.push('/onboarding');
      }
    }
  }, [user, loading, router]);

  useEffect(() => {
    if (user && user.preferences) {
      loadThreads();
    }
  }, [user]);

  useEffect(() => {
    const handleRefresh = () => {
      loadThreads();
      setSelectedThreadId(null);
    };
    window.addEventListener('refreshThreads', handleRefresh);
    return () => window.removeEventListener('refreshThreads', handleRefresh);
  }, []);

  const loadThreads = async () => {
    setLoadingThreads(true);
    try {
      const fetchedThreads = await threadAPI.getAll();
      setThreads(fetchedThreads);
    } catch (error) {
      console.error('Failed to load threads:', error);
    } finally {
      setLoadingThreads(false);
    }
  };

  const handleInboxMessageAssigned = () => {
    setSelectedInboxMessage(null);
    setSelectedThreadId(null);
  };

  const handleThreadCreated = (newThread: Thread) => {
    setThreads([...threads, newThread]);
    setSelectedThreadId(newThread.id);
    setSelectedInboxMessage(null);
  };

  const handleThreadSelect = (threadId: string) => {
    setSelectedThreadId(threadId);
    setSelectedInboxMessage(null);
    setViewMode('chat');
  };

  const handleInboxMessageSelect = (message: InboxMessage) => {
    setSelectedInboxMessage(message);
    setSelectedThreadId(null);
    setViewMode('inbox');
  };

  const handleViewOffers = () => {
    setViewMode('offers');
    setSelectedThreadId(null);
    setSelectedInboxMessage(null);
  };

  const handleNavigateToThread = (threadId: string) => {
    setSelectedThreadId(threadId);
    setSelectedInboxMessage(null);
    setViewMode('chat');
  };

  if (loading || loadingThreads) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-gray-600">Loading...</div>
      </div>
    );
  }

  if (!user || !user.preferences) {
    return null;
  }

  return (
    <SidebarProvider
      defaultOpen
      style={
        {
          '--sidebar-width': '22rem',
          '--header-height': 'calc(var(--spacing) * 12)',
        } as React.CSSProperties
      }
    >
      <AppSidebar
        variant="inset"
        threads={threads}
        selectedThreadId={selectedThreadId}
        selectedInboxMessageId={selectedInboxMessage?.id || null}
        onThreadSelect={handleThreadSelect}
        onThreadCreated={handleThreadCreated}
        onInboxMessageSelect={handleInboxMessageSelect}
      />
      <SidebarInset className="overflow-x-hidden">
        {/* <TopNavBar onViewOffers={handleViewOffers} /> */}
        <div className="flex flex-1 flex-col overflow-x-hidden">
          {viewMode === 'offers' ? (
            <OffersPane onNavigateToThread={handleNavigateToThread} />
          ) : (
            <ChatPane
              selectedThreadId={selectedThreadId}
              selectedInboxMessage={selectedInboxMessage}
              threads={threads}
              onInboxMessageAssigned={handleInboxMessageAssigned}
            />
          )}
        </div>
      </SidebarInset>
    </SidebarProvider>
  );
}

export default function DashboardPage() {
  return (
    <Suspense fallback={
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-gray-600">Loading...</div>
      </div>
    }>
      <DashboardContent />
    </Suspense>
  );
}

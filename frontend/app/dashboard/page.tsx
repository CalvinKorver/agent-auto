'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';
import { AppSidebar } from '@/components/dashboard/AppSidebar';
import ChatPane from '@/components/dashboard/ChatPane';
import TopNavBar from '@/components/dashboard/TopNavBar';
import OffersPane from '@/components/dashboard/OffersPane';
import { Thread, threadAPI, InboxMessage } from '@/lib/api';
import { SidebarProvider, SidebarInset } from '@/components/ui/sidebar';

type ViewMode = 'chat' | 'inbox' | 'offers';

export default function DashboardPage() {
  const { user, loading } = useAuth();
  const router = useRouter();
  const [threads, setThreads] = useState<Thread[]>([]);
  const [selectedThreadId, setSelectedThreadId] = useState<string | null>(null);
  const [selectedInboxMessage, setSelectedInboxMessage] = useState<InboxMessage | null>(null);
  const [loadingThreads, setLoadingThreads] = useState(false);
  const [viewMode, setViewMode] = useState<ViewMode>('chat');

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
      <SidebarInset>
        {/* <TopNavBar onViewOffers={handleViewOffers} /> */}
        <div className="flex flex-1 flex-col">
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

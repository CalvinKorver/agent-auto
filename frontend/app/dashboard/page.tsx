'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';
import { AppSidebar } from '@/components/dashboard/AppSidebar';
import ChatPane from '@/components/dashboard/ChatPane';
import { Thread, threadAPI, InboxMessage } from '@/lib/api';
import { SidebarProvider, SidebarInset } from '@/components/ui/sidebar';

export default function DashboardPage() {
  const { user, loading } = useAuth();
  const router = useRouter();
  const [threads, setThreads] = useState<Thread[]>([]);
  const [selectedThreadId, setSelectedThreadId] = useState<string | null>(null);
  const [selectedInboxMessage, setSelectedInboxMessage] = useState<InboxMessage | null>(null);
  const [loadingThreads, setLoadingThreads] = useState(false);

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
  };

  const handleInboxMessageSelect = (message: InboxMessage) => {
    setSelectedInboxMessage(message);
    setSelectedThreadId(null);
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
    <SidebarProvider>
      <div className="flex h-screen w-full overflow-hidden">
        <AppSidebar
          threads={threads}
          selectedThreadId={selectedThreadId}
          selectedInboxMessageId={selectedInboxMessage?.id || null}
          onThreadSelect={handleThreadSelect}
          onThreadCreated={handleThreadCreated}
          onInboxMessageSelect={handleInboxMessageSelect}
        />
        <SidebarInset className="flex-1 flex flex-col overflow-hidden">
          <ChatPane
            selectedThreadId={selectedThreadId}
            selectedInboxMessage={selectedInboxMessage}
            threads={threads}
            onInboxMessageAssigned={handleInboxMessageAssigned}
          />
        </SidebarInset>
      </div>
    </SidebarProvider>
  );
}

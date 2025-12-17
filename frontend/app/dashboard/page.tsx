'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';
import ThreadPane from '@/components/dashboard/ThreadPane';
import ChatPane from '@/components/dashboard/ChatPane';
import { Thread, threadAPI } from '@/lib/api';

export default function DashboardPage() {
  const { user, loading } = useAuth();
  const router = useRouter();
  const [threads, setThreads] = useState<Thread[]>([]);
  const [selectedThreadId, setSelectedThreadId] = useState<string | null>(null);
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

  const handleThreadCreated = (newThread: Thread) => {
    setThreads([...threads, newThread]);
    setSelectedThreadId(newThread.id);
  };

  const handleThreadSelect = (threadId: string) => {
    setSelectedThreadId(threadId);
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

  const targetVehicle = `${user.preferences.year} ${user.preferences.make} ${user.preferences.model}`;

  return (
    <div className="flex h-screen overflow-hidden">
      <ThreadPane
        targetVehicle={targetVehicle}
        threads={threads}
        selectedThreadId={selectedThreadId}
        onThreadSelect={handleThreadSelect}
        onThreadCreated={handleThreadCreated}
      />
      <ChatPane selectedThreadId={selectedThreadId} />
    </div>
  );
}

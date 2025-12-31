'use client';

import { useEffect, useState, Suspense } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';
import { AppSidebar } from '@/components/dashboard/AppSidebar';
import ChatPane from '@/components/dashboard/ChatPane';
import TopNavBar from '@/components/dashboard/TopNavBar';
import OffersPane from '@/components/dashboard/OffersPane';
import { Thread, InboxMessage, TrackedOffer, Dealer, dashboardAPI, dealersAPI } from '@/lib/api';
import { SidebarProvider, SidebarInset } from '@/components/ui/sidebar';
import { toast } from 'sonner';

type ViewMode = 'chat' | 'inbox' | 'offers';

function DashboardContent() {
  const { user, loading, refreshUser } = useAuth();
  const router = useRouter();
  const searchParams = useSearchParams();
  const [threads, setThreads] = useState<Thread[]>([]);
  const [inboxMessages, setInboxMessages] = useState<InboxMessage[]>([]);
  const [offers, setOffers] = useState<TrackedOffer[]>([]);
  const [dealers, setDealers] = useState<Dealer[]>([]);
  const [selectedThreadId, setSelectedThreadId] = useState<string | null>(null);
  const [selectedInboxMessage, setSelectedInboxMessage] = useState<InboxMessage | null>(null);
  const [selectedOfferId, setSelectedOfferId] = useState<string | null>(null);
  const [loadingDashboard, setLoadingDashboard] = useState(false);
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
      loadDashboard();
    }
  }, [user]);

  useEffect(() => {
    const handleRefresh = () => {
      loadDashboard();
      setSelectedThreadId(null);
    };
    window.addEventListener('refreshThreads', handleRefresh);
    return () => window.removeEventListener('refreshThreads', handleRefresh);
  }, []);

  const loadDashboard = async () => {
    setLoadingDashboard(true);
    try {
      const [dashboardData, dealersData] = await Promise.all([
        dashboardAPI.getDashboard(),
        dealersAPI.getDealers().catch(() => []), // Don't fail if dealers endpoint fails
      ]);
      setThreads(dashboardData.threads);
      setInboxMessages(dashboardData.inboxMessages);
      setOffers(dashboardData.offers);
      setDealers(dealersData);
    } catch (error) {
      console.error('Failed to load dashboard data:', error);
    } finally {
      setLoadingDashboard(false);
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
    setSelectedOfferId(null);
  };

  const handleThreadSelect = (threadId: string) => {
    setSelectedThreadId(threadId);
    setSelectedInboxMessage(null);
    setSelectedOfferId(null);
    setViewMode('chat');
  };

  const handleInboxMessageSelect = (message: InboxMessage) => {
    setSelectedInboxMessage(message);
    setSelectedThreadId(null);
    setSelectedOfferId(null);
    setViewMode('inbox');
  };

  const handleOfferSelect = (offer: TrackedOffer) => {
    setSelectedOfferId(offer.id);
    setSelectedThreadId(offer.threadId);
    setSelectedInboxMessage(null);
    setViewMode('chat');
  };

  const handleViewOffers = () => {
    setViewMode('offers');
    setSelectedThreadId(null);
    setSelectedInboxMessage(null);
    setSelectedOfferId(null);
  };

  const handleNavigateToThread = (threadId: string) => {
    setSelectedThreadId(threadId);
    setSelectedInboxMessage(null);
    setSelectedOfferId(null);
    setViewMode('chat');
  };

  const handleGoToDashboard = () => {
    setSelectedThreadId(null);
    setSelectedInboxMessage(null);
    setSelectedOfferId(null);
    setViewMode('chat');
  };

  if (loading || loadingDashboard) {
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
        inboxMessages={inboxMessages}
        selectedThreadId={selectedThreadId}
        selectedInboxMessageId={selectedInboxMessage?.id || null}
        offers={offers}
        selectedOfferId={selectedOfferId}
        onThreadSelect={handleThreadSelect}
        onThreadCreated={handleThreadCreated}
        onInboxMessageSelect={handleInboxMessageSelect}
        onOfferSelect={handleOfferSelect}
        onGoToDashboard={handleGoToDashboard}
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
              offers={offers}
              dealers={dealers}
              onInboxMessageAssigned={handleInboxMessageAssigned}
              onNavigateToThread={handleNavigateToThread}
              onOfferDeleted={loadDashboard}
              onDealersUpdated={loadDashboard}
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

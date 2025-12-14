'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';
import Header from '@/components/Header';

export default function DashboardPage() {
  const { user, loading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!loading) {
      if (!user) {
        router.push('/login');
      } else if (!user.preferences) {
        router.push('/onboarding');
      }
    }
  }, [user, loading, router]);

  if (loading) {
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
    <>
      <Header />
      <div className="min-h-screen bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 py-8">
          <div className="bg-white rounded-lg shadow-md p-8">
            <h1 className="text-3xl font-bold text-gray-900 mb-4">Dashboard</h1>

            <div className="mb-8">
              <h2 className="text-xl font-semibold text-gray-800 mb-2">Target Vehicle</h2>
              <div className="bg-blue-50 border border-blue-200 rounded-md p-4">
                <p className="text-lg">
                  <span className="font-medium">{user.preferences.year} {user.preferences.make} {user.preferences.model}</span>
                </p>
              </div>
            </div>

            <div className="border-t border-gray-200 pt-8">
              <h2 className="text-xl font-semibold text-gray-800 mb-4">Active Negotiations</h2>
              <div className="text-center py-12 bg-gray-50 rounded-lg">
                <p className="text-gray-600">No active negotiations</p>
                <p className="text-sm text-gray-500 mt-2">
                  Thread management coming in Phase 3!
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </>
  );
}

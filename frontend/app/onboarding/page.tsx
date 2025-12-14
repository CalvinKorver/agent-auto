'use client';

import { useRouter } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';
import PreferencesForm from '@/components/PreferencesForm';
import Header from '@/components/Header';

export default function OnboardingPage() {
  const router = useRouter();
  const { refreshUser } = useAuth();

  const handleSuccess = async () => {
    await refreshUser();
    router.push('/dashboard');
  };

  return (
    <>
      <Header />
      <PreferencesForm onSuccess={handleSuccess} />
    </>
  );
}

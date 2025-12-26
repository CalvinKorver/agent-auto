'use client';

import { useRouter } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';
import OnboardingContainer from '@/components/onboarding/OnboardingContainer';

export default function OnboardingPage() {
  const router = useRouter();
  const { refreshUser } = useAuth();

  const handleSuccess = async () => {
    await refreshUser();
    router.push('/dashboard');
  };

  return <OnboardingContainer onSuccess={handleSuccess} />;
}

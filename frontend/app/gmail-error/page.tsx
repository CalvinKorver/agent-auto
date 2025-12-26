'use client';

import { useEffect, useState, Suspense } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { Button } from '@/components/ui/button';

function GmailErrorContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [errorMessage, setErrorMessage] = useState('');

  useEffect(() => {
    const error = searchParams.get('error');

    switch (error) {
      case 'no_code':
        setErrorMessage('Authorization failed. Please try again.');
        break;
      case 'invalid_state':
        setErrorMessage('Invalid session. Please try connecting again.');
        break;
      case 'token_exchange_failed':
        setErrorMessage('Failed to connect Gmail. Please try again.');
        break;
      case 'store_failed':
        setErrorMessage('Failed to save Gmail connection. Please try again.');
        break;
      case 'access_denied':
        setErrorMessage('You need to grant permission to send emails via Gmail.');
        break;
      default:
        setErrorMessage('Something went wrong. Please try again.');
    }
  }, [searchParams]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-background">
      <div className="max-w-md w-full space-y-6 p-8">
        <div className="text-center">
          <div className="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-red-100 dark:bg-red-950">
            <svg
              className="h-6 w-6 text-red-600 dark:text-red-400"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </div>
          <h2 className="mt-4 text-2xl font-bold">Gmail Connection Failed</h2>
          <p className="mt-2 text-muted-foreground">{errorMessage}</p>
        </div>

        <div className="space-y-3">
          <Button
            onClick={() => router.push('/dashboard')}
            className="w-full"
          >
            Return to Dashboard
          </Button>
          <Button
            onClick={() => router.push('/dashboard')}
            variant="outline"
            className="w-full"
          >
            Try Again
          </Button>
        </div>

        <p className="text-center text-sm text-muted-foreground">
          If this problem persists, please contact support.
        </p>
      </div>
    </div>
  );
}

export default function GmailErrorPage() {
  return (
    <Suspense fallback={
      <div className="min-h-screen flex items-center justify-center bg-background">
        <div className="text-gray-600">Loading...</div>
      </div>
    }>
      <GmailErrorContent />
    </Suspense>
  );
}

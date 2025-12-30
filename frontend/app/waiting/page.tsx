'use client';

import Image from 'next/image';
import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';

export default function WaitingPage() {
  const [email, setEmail] = useState('');
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setMessage(null);

    try {
      const response = await fetch('/api/waitlist', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email }),
      });

      const data = await response.json();

      if (response.ok) {
        setMessage({ type: 'success', text: data.message || 'Successfully added to waitlist!' });
        setEmail('');
      } else {
        setMessage({ type: 'error', text: data.error || 'Failed to add to waitlist. Please try again.' });
      }
    } catch {
      setMessage({ type: 'error', text: 'An error occurred. Please try again later.' });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-50 to-white dark:from-gray-900 dark:to-gray-800 flex items-center justify-center px-4">
      <div className="max-w-2xl w-full text-center">
        <div className="mb-8 flex justify-center">
          <Image
            src="/otto-ai-round.png"
            alt="Otto AI"
            width={120}
            height={120}
            className="w-24 h-24 md:w-32 md:h-32 rounded-full"
          />
        </div>
        
        <h1 className="text-xl md:text-xl font-bold tracking-tight mb-6 text-foreground">
          We&apos;re Currently in Private Beta
        </h1>
        
        <p className="text-lg text-muted-foreground mb-8 max-w-xl mx-auto">
          Thank you for your interest in Otto AI! We&apos;re currently accepting a limited number of users. 
          Join our waitlist to be notified when we launch.
        </p>

        <form onSubmit={handleSubmit} className="max-w-md mx-auto mb-8">
          <div className="flex gap-2 flex-col sm:flex-row">
            <Input
              type="email"
              placeholder="Enter your email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              disabled={loading}
              className="flex-1"
            />
            <Button 
              type="submit" 
              disabled={loading || !email}
              size="lg"
            >
              {loading ? 'Joining...' : 'Join Waitlist'}
            </Button>
          </div>
          
          {message && (
            <div className={`mt-4 p-3 rounded-md text-sm ${
              message.type === 'success' 
                ? 'bg-green-50 dark:bg-green-900/20 text-green-800 dark:text-green-200 border border-green-200 dark:border-green-800' 
                : 'bg-red-50 dark:bg-red-900/20 text-red-800 dark:text-red-200 border border-red-200 dark:border-red-800'
            }`}>
              {message.text}
            </div>
          )}
        </form>
        
        <div className="mt-12 text-sm text-muted-foreground">
          <p>&copy; 2025 OttoAI. All rights reserved.</p>
        </div>
      </div>
    </div>
  );
}


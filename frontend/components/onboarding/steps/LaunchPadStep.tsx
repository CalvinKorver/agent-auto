import { useState } from 'react';
import { Button } from '@/components/ui/button';
import PollyAvatar from '../PollyAvatar';

interface LaunchPadStepProps {
  vehicleData: {
    make: string;
    model: string;
    year: number;
  };
  onComplete: () => Promise<void>;
}

export default function LaunchPadStep({ vehicleData, onComplete }: LaunchPadStepProps) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleGetStarted = async () => {
    try {
      setLoading(true);
      setError('');
      await onComplete();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to save preferences. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-6 animate-in fade-in duration-500">
      {/* Lolo AI's message - no bubble, just text like ChatPane */}
      <div className="flex items-start gap-1">
        <PollyAvatar />
        <div className="space-y-2 text-lg leading-relaxed text-foreground">
          <div className="font-semibold text-xl">Ready to disrupt the status quo?</div>
          <div>
            Our AI is ready to help you save an average of $1,280 or more while shielding you from &quot;Communication Fatigue&quot;.
          </div>
          <div className="bg-muted/50 rounded-md p-3 mt-3">
            <div className="text-xs font-medium text-muted-foreground mb-1">Your Target Vehicle:</div>
            <div className="font-semibold">
              {vehicleData.year} {vehicleData.make} {vehicleData.model}
            </div>
          </div>
          <div className="text-sm text-muted-foreground italic mt-2">
            By clicking below, you&apos;re one step closer to your new car without the dread of the dealership floor.
          </div>
        </div>
      </div>

      {error && (
        <div className="bg-destructive/10 border border-destructive/20 text-destructive rounded-md p-3 text-sm">
          {error}
        </div>
      )}

      <div className="flex justify-center pt-4">
        <Button onClick={handleGetStarted} disabled={loading} size="lg" className="px-8">
          {loading ? 'Setting Up...' : 'Get Started'}
        </Button>
      </div>
    </div>
  );
}

import { useState, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import PollyAvatar from '../PollyAvatar';

interface JourneyStateStepProps {
  onSelectJourney: (type: 'exploring' | 'high-intent') => void;
  isActive: boolean;
}

export default function JourneyStateStep({ onSelectJourney, isActive }: JourneyStateStepProps) {
  const [selected, setSelected] = useState<'exploring' | 'high-intent' | null>(null);
  const [showExploringMessage, setShowExploringMessage] = useState(false);
  const [showContent, setShowContent] = useState(false);

  useEffect(() => {
    if (isActive) {
      const timer = setTimeout(() => {
        setShowContent(true);
      }, 500); // 500ms delay after the step becomes active
      return () => clearTimeout(timer);
    }
  }, [isActive]);

  const handleExploring = () => {
    setSelected('exploring');
    setShowExploringMessage(true);
  };

  const handleHighIntent = () => {
    setSelected('high-intent');
    onSelectJourney('high-intent');
  };

  return (
    <div className="space-y-4">
      {/* Lolo AI's message and buttons - animate in together */}
      {showContent && (
        <div className="space-y-3 animate-in fade-in slide-in-from-bottom-2 duration-1000">
          <div className="flex gap-3 items-start">
            <PollyAvatar />
            <div className="text-md leading-relaxed text-foreground">
              <div>Where are you in your journey?</div>
            </div>
          </div>

          {!selected && isActive && (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
              <button
                onClick={handleHighIntent}
                className="p-4 border-2 border-border rounded-lg hover:border-primary hover:bg-accent/50 transition-all text-center group bg-accent text-sm">

                I know which car I want!

              </button>

              <button
                onClick={handleExploring}
                className="p-4 border-2 border-border rounded-lg hover:border-primary hover:bg-accent/50 transition-all text-center group bg-accent text-sm">
                  I don&apos;t know which car I want yet.

              </button>
            </div>
          )}
        </div>
      )}

      {selected === 'exploring' && showExploringMessage && (
        <div className="space-y-3 animate-in fade-in duration-500">
          <div className="flex justify-end">
            <div className="max-w-[70%]">
              <div className="text-xs text-muted-foreground mb-1 text-right">You</div>
              <div className="bg-card text-card-foreground border border-border rounded-lg px-4 py-2.5">
                <div className="text-sm">Just exploring!</div>
              </div>
            </div>
          </div>

          <div className="flex items-start gap-3">
            <PollyAvatar />
            <div className="text-md leading-relaxed text-foreground">
              <div>
                Great! Our coaching mode is coming soon. For now, let me help you if you know what car you want. Would you like to continue with vehicle selection?
              </div>
            </div>
          </div>

          <div className="flex justify-center pt-2">
            <Button onClick={handleHighIntent} size="lg">
              Yes, Let&apos;s Continue
            </Button>
          </div>
        </div>
      )}

      {selected === 'high-intent' && (
        <div className="flex justify-end animate-in fade-in duration-500">
          <div className="max-w-[70%]">
            <div className="text-xs text-muted-foreground mb-1 text-right">You</div>
            <div className="bg-card text-card-foreground border border-border rounded-lg px-4 py-2.5">
              <div className="text-sm">I know which car I want!</div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

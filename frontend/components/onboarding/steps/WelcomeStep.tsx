import { Button } from '@/components/ui/button';
import PollyAvatar from '../PollyAvatar';

interface WelcomeStepProps {
  onNext: () => void;
  isActive: boolean;
}

export default function WelcomeStep({ onNext, isActive }: WelcomeStepProps) {
  return (
    <div className="space-y-4 animate-in fade-in duration-500">
      {/* Lolo AI's message - no bubble, just text like ChatPane */}
      <div className="flex items-start gap-3">
        <PollyAvatar />
        <div className="space-y-1 text-md leading-relaxed text-foreground">
          <div className="font-semibold text-md">Welcome to Lolo AI</div>
          <p className="text-base">
            We&apos;re here to guide you through your car-buying journey. Our AI agent handles the
            &quot;haggling&quot; so you can focus on the driving.
          </p>
        </div>
      </div>

      {isActive && (
        <div className="flex justify-center pt-2">
          <Button onClick={onNext} size="lg" className="px-8">
            Let&apos;s Get Started
          </Button>
        </div>
      )}

      {!isActive && (
        <div className="flex justify-end animate-in fade-in duration-300">
          <div className="max-w-[70%]">
            <div className="text-xs text-muted-foreground mb-1 text-right">You</div>
            <div className="bg-card text-card-foreground border border-border rounded-lg px-4 py-2.5">
              <div className="text-sm">Let&apos;s Get Started</div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

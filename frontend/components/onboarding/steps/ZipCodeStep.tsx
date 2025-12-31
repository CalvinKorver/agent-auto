'use client';

import { useState, useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Button } from '@/components/ui/button';
import { Field, FieldLabel, FieldError } from '@/components/ui/field';
import { Input } from '@/components/ui/input';
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip';
import { Info } from 'lucide-react';
import PollyAvatar from '../PollyAvatar';

const zipCodeSchema = z.object({
  zipCode: z.string().regex(/^\d{5}$/, 'Zip code must be 5 digits'),
});

type ZipCodeFormData = z.infer<typeof zipCodeSchema>;

interface ZipCodeStepProps {
  onComplete: (zipCode: string) => void;
  isActive: boolean;
}

export default function ZipCodeStep({ onComplete, isActive }: ZipCodeStepProps) {
  const [showContent, setShowContent] = useState(false);
  const [userResponse, setUserResponse] = useState<string | null>(null);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<ZipCodeFormData>({
    resolver: zodResolver(zipCodeSchema),
    defaultValues: {
      zipCode: '',
    },
  });

  useEffect(() => {
    if (isActive) {
      const timer = setTimeout(() => {
        setShowContent(true);
      }, 500);
      return () => clearTimeout(timer);
    }
  }, [isActive]);

  const onSubmit = (data: ZipCodeFormData) => {
    setUserResponse(data.zipCode);
    onComplete(data.zipCode);
  };

  // If completed (not active), show condensed summary
  if (!isActive && userResponse) {
    return (
      <div className="space-y-6 animate-in fade-in duration-500">
        <div className="flex justify-end">
          <div className="max-w-[70%]">
            <div className="text-xs text-muted-foreground mb-1 text-right">You</div>
            <div className="bg-card text-card-foreground border border-border rounded-lg px-4 py-3">
              <div className="text-sm">{userResponse}</div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-4 animate-in fade-in duration-500">
      {showContent && (
        <div className="animate-in fade-in slide-in-from-bottom-4 duration-500">
          <div className="flex items-start gap-3 mb-2">
            <PollyAvatar />
            <div className="text-sm leading-relaxed text-foreground">
              <p className="text-sm">One last question: What Zip Code are you looking in?</p>
            </div>
          </div>

          <div className="pl-13">
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
              <Field>
                <div className="flex items-center gap-2">
                  <FieldLabel htmlFor="zipCode">Zip Code</FieldLabel>
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <Info className="h-4 w-4 text-muted-foreground cursor-help" />
                    </TooltipTrigger>
                    <TooltipContent>
                      <p>Zip Code will help us show you results that are local to you</p>
                    </TooltipContent>
                  </Tooltip>
                </div>
                <Input
                  id="zipCode"
                  type="text"
                  placeholder="12345"
                  maxLength={5}
                  {...register('zipCode')}
                  className="w-full"
                />
                <FieldError>{errors.zipCode?.message}</FieldError>
              </Field>

              <div className="flex justify-center pt-2">
                <Button type="submit" size="lg" className="px-8">
                  Continue
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}


'use client';

import { useState, useEffect } from 'react';
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Button } from '@/components/ui/button';
import { Field, FieldLabel, FieldError } from '@/components/ui/field';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Input } from '@/components/ui/input';
import PollyAvatar from '../PollyAvatar';
import { VEHICLE_MAKES, VEHICLE_YEARS } from '@/lib/data/vehicleData';

const vehicleSchema = z.object({
  make: z.string().min(1, 'Make is required'),
  model: z.string().min(1, 'Model is required'),
  year: z.number().min(2000).max(2030),
});

type VehicleFormData = z.infer<typeof vehicleSchema>;

interface VehicleSpecStepProps {
  onComplete: (data: VehicleFormData) => void;
  isActive: boolean;
}

export default function VehicleSpecStep({ onComplete, isActive }: VehicleSpecStepProps) {
  const [showMake, setShowMake] = useState(true);
  const [showModel, setShowModel] = useState(false);
  const [showYear, setShowYear] = useState(false);
  const [userResponses, setUserResponses] = useState<{
    make?: string;
    model?: string;
    year?: number;
  }>({});

  const {
    register,
    handleSubmit,
    watch,
    control,
    formState: { errors },
  } = useForm<VehicleFormData>({
    resolver: zodResolver(vehicleSchema),
    defaultValues: {
      make: '',
      model: '',
      year: 2024,
    },
  });

  const makeValue = watch('make');
  const modelValue = watch('model');
  const yearValue = watch('year');
  const modelHasValue = (modelValue ?? '').trim().length > 0;

  // Progressive disclosure logic - only when active
  useEffect(() => {
    if (isActive && makeValue && !showModel) {
      setUserResponses((prev) => ({ ...prev, make: makeValue }));
      setTimeout(() => setShowModel(true), 300);
    }
  }, [makeValue, showModel, isActive]);

  useEffect(() => {
    if (isActive && modelHasValue && showModel && !showYear) {
      setUserResponses((prev) => ({ ...prev, model: modelValue }));
      setTimeout(() => setShowYear(true), 300);
    }
  }, [modelHasValue, modelValue, showModel, showYear, isActive]);

  useEffect(() => {
    if (isActive && yearValue && showYear) {
      setUserResponses((prev) => ({ ...prev, year: yearValue }));
    }
  }, [yearValue, showYear, isActive]);

  const onSubmit = (data: VehicleFormData) => {
    onComplete(data);
  };

  // If completed (not active), show condensed summary
  if (!isActive && userResponses.make && userResponses.model && userResponses.year) {
    return (
      <div className="space-y-6 animate-in fade-in duration-500">

        <div className="flex justify-end">
          <div className="max-w-[70%]">
            <div className="text-xs text-muted-foreground mb-1 text-right">You</div>
            <div className="bg-card text-card-foreground border border-border rounded-lg px-4 py-3">
              <div className="text-sm">
                {userResponses.year} {userResponses.make} {userResponses.model}
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-4 animate-in fade-in duration-500">
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        {/* Make Selection */}
        {showMake && (
          <div className="animate-in fade-in slide-in-from-bottom-4 duration-500">
            <div className="flex items-start gap-3 mb-2">
              <PollyAvatar />
              <div className="text-md leading-relaxed text-foreground">
                <p className="text-base">What make are you looking for?</p>
              </div>
            </div>

            <div className="pl-13">
              <Field>
                <FieldLabel htmlFor="make">Select Make</FieldLabel>
                <Controller
                  name="make"
                  control={control}
                  render={({ field }) => (
                    <Select onValueChange={field.onChange} value={field.value}>
                      <SelectTrigger id="make" className="w-full">
                        <SelectValue placeholder="Choose a make..." />
                      </SelectTrigger>
                      <SelectContent>
                        {VEHICLE_MAKES.map((make) => (
                          <SelectItem key={make} value={make}>
                            {make}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  )}
                />
                <FieldError>{errors.make?.message}</FieldError>
              </Field>
            </div>

            {showModel && (
              <div className="pl-13 mt-4 animate-in fade-in slide-in-from-bottom-4 duration-500">
                <Field>
                  <FieldLabel htmlFor="model">Model Name</FieldLabel>
                  <Input
                    id="model"
                    type="text"
                    {...register('model')}
                    placeholder="e.g., CX-90, Camry, Model 3"
                  />
                  <FieldError>{errors.model?.message}</FieldError>
                </Field>

                {showYear && (
                  <div className="mt-4">
                    <Field>
                      <FieldLabel htmlFor="year">Select Year</FieldLabel>
                      <Controller
                        name="year"
                        control={control}
                        render={({ field }) => (
                          <Select
                            onValueChange={(value) => field.onChange(Number(value))}
                            value={field.value?.toString()}
                          >
                            <SelectTrigger id="year" className="w-full">
                              <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                              {VEHICLE_YEARS.map((year) => (
                                <SelectItem key={year} value={year.toString()}>
                                  {year}
                                </SelectItem>
                              ))}
                            </SelectContent>
                          </Select>
                        )}
                      />
                      <FieldError>{errors.year?.message}</FieldError>
                    </Field>
                  </div>
                )}
              </div>
            )}

            {userResponses.make && !isActive && (
              <div className="flex justify-end mt-3 animate-in fade-in duration-300">
                <div className="max-w-[70%]">
                  <div className="text-xs text-muted-foreground mb-1 text-right">You</div>
                  <div className="bg-card text-card-foreground border border-border rounded-lg px-4 py-2.5">
                    <div className="text-sm">{userResponses.make}</div>
                  </div>
                </div>
              </div>
            )}
          </div>
        )}

        {showModel && userResponses.model && !isActive && (
          <div className="flex justify-end mt-3 animate-in fade-in duration-300">
            <div className="max-w-[70%]">
              <div className="text-xs text-muted-foreground mb-1 text-right">You</div>
              <div className="bg-card text-card-foreground border border-border rounded-lg px-4 py-2.5">
                <div className="text-sm">{userResponses.model}</div>
              </div>
            </div>
          </div>
        )}

        {showYear && userResponses.year && !isActive && (
          <div className="flex justify-end mt-3 animate-in fade-in duration-300">
            <div className="max-w-[70%]">
              <div className="text-xs text-muted-foreground mb-1 text-right">You</div>
              <div className="bg-card text-card-foreground border border-border rounded-lg px-4 py-2.5">
                <div className="text-sm">{userResponses.year}</div>
              </div>
            </div>
          </div>
        )}

        {/* Submit Button */}
        {showYear && yearValue && (
          <div className="flex justify-center pt-2 animate-in fade-in duration-500">
            <Button type="submit" size="lg" className="px-8">
              Continue
            </Button>
          </div>
        )}
      </form>
    </div>
  );
}



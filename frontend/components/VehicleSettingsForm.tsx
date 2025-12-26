'use client';

import { useState, useEffect } from 'react';
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Button } from '@/components/ui/button';
import { Field, FieldLabel, FieldError } from '@/components/ui/field';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Input } from '@/components/ui/input';
import { VEHICLE_MAKES } from '@/lib/data/vehicleData';
import { preferencesAPI } from '@/lib/api';
import { toast } from 'sonner';

const vehicleSchema = z.object({
  make: z.string().min(1, 'Make is required'),
  model: z.string().min(1, 'Model is required'),
  year: z.number().min(1960).max(2026),
});

type VehicleFormData = z.infer<typeof vehicleSchema>;

// Generate years from 1960 to 2026
const VEHICLE_YEARS = Array.from(
  { length: 2026 - 1960 + 1 },
  (_, i) => 2026 - i
);

export default function VehicleSettingsForm() {
  const [loading, setLoading] = useState(false);
  const [initialLoading, setInitialLoading] = useState(true);

  const {
    register,
    handleSubmit,
    control,
    reset,
    formState: { errors },
  } = useForm<VehicleFormData>({
    resolver: zodResolver(vehicleSchema),
    defaultValues: {
      make: '',
      model: '',
      year: 2024,
    },
  });

  // Load existing preferences
  useEffect(() => {
    const loadPreferences = async () => {
      try {
        const prefs = await preferencesAPI.get();
        reset({
          make: prefs.make,
          model: prefs.model,
          year: prefs.year,
        });
      } catch (error: any) {
        // If preferences don't exist, that's okay - form will start empty
        if (error.response?.status !== 404) {
          console.error('Failed to load preferences:', error);
        }
      } finally {
        setInitialLoading(false);
      }
    };

    loadPreferences();
  }, [reset]);

  const onSubmit = async (data: VehicleFormData) => {
    setLoading(true);
    try {
      await preferencesAPI.update(data.year, data.make, data.model);
      toast.success('Vehicle settings saved successfully');
    } catch (error: any) {
      console.error('Failed to save preferences:', error);
      if (error.response?.status === 404) {
        // If preferences don't exist, create them instead
        try {
          await preferencesAPI.create(data.year, data.make, data.model);
          toast.success('Vehicle settings saved successfully');
        } catch (createError) {
          console.error('Failed to create preferences:', createError);
          toast.error('Failed to save vehicle settings');
        }
      } else {
        toast.error('Failed to save vehicle settings');
      }
    } finally {
      setLoading(false);
    }
  };

  const onDiscard = () => {
    reset({
      make: '',
      model: '',
      year: 2024,
    });
    toast.info('Form cleared');
  };

  if (initialLoading) {
    return (
      <div className="flex items-center justify-center py-8">
        <div className="text-muted-foreground">Loading settings...</div>
      </div>
    );
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
      <div className="space-y-4">
        <Field>
          <FieldLabel htmlFor="make">Make</FieldLabel>
          <Controller
            name="make"
            control={control}
            render={({ field }) => (
              <Select onValueChange={field.onChange} value={field.value}>
                <SelectTrigger id="make" className="w-full">
                  <SelectValue placeholder="Select a make..." />
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

        <Field>
          <FieldLabel htmlFor="model">Model</FieldLabel>
          <Input
            id="model"
            type="text"
            {...register('model')}
            placeholder="e.g., CX-90, Camry, Model 3"
          />
          <FieldError>{errors.model?.message}</FieldError>
        </Field>

        <Field>
          <FieldLabel htmlFor="year">Year</FieldLabel>
          <Controller
            name="year"
            control={control}
            render={({ field }) => (
              <Select
                onValueChange={(value) => field.onChange(Number(value))}
                value={field.value?.toString()}
              >
                <SelectTrigger id="year" className="w-full">
                  <SelectValue placeholder="Select a year..." />
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

      <div className="flex gap-4 pt-4">
        <Button type="submit" disabled={loading}>
          {loading ? 'Saving...' : 'Submit'}
        </Button>
        <Button type="button" variant="outline" onClick={onDiscard}>
          Discard
        </Button>
      </div>
    </form>
  );
}








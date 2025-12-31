'use client';

import { useState, useEffect } from 'react';
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Button } from '@/components/ui/button';
import { Field, FieldLabel, FieldError } from '@/components/ui/field';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Input } from '@/components/ui/input';
import { preferencesAPI, modelsAPI, VehicleModelsResponse, trimsAPI, Trim } from '@/lib/api';
import { useAuth } from '@/contexts/AuthContext';
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip';
import { Info } from 'lucide-react';
import { toast } from 'sonner';

const vehicleSchema = z.object({
  make: z.string().min(1, 'Make is required'),
  model: z.string().min(1, 'Model is required'),
  year: z.number().min(1960).max(2026),
  trim: z.string().min(1, 'Trim is required'),
  zipCode: z.string().regex(/^\d{5}$/, 'Zip code must be 5 digits').optional().or(z.literal('')),
});

type VehicleFormData = z.infer<typeof vehicleSchema>;

// Generate years from 1960 to 2026
const VEHICLE_YEARS = Array.from(
  { length: 2026 - 1960 + 1 },
  (_, i) => 2026 - i
);

export default function VehicleSettingsForm() {
  const { user } = useAuth();
  const [loading, setLoading] = useState(false);
  const [initialLoading, setInitialLoading] = useState(true);
  const [vehicleModels, setVehicleModels] = useState<VehicleModelsResponse>({});
  const [makes, setMakes] = useState<string[]>([]);
  const [isLoadingModels, setIsLoadingModels] = useState(true);
  const [modelsError, setModelsError] = useState<string | null>(null);
  const [trims, setTrims] = useState<Trim[]>([]);
  const [isLoadingTrims, setIsLoadingTrims] = useState(false);

  const {
    register,
    handleSubmit,
    control,
    reset,
    resetField,
    watch,
    formState: { errors },
  } = useForm<VehicleFormData>({
    resolver: zodResolver(vehicleSchema),
    defaultValues: {
      make: '',
      model: '',
      year: 2024,
      trim: '',
      zipCode: '',
    },
  });

  const makeValue = watch('make');
  const modelValue = watch('model');
  const yearValue = watch('year');
  const trimValue = watch('trim');

  // Reset model when make changes if current model is not available for new make
  useEffect(() => {
    // Only run if we have models data loaded and both make and model are set
    if (Object.keys(vehicleModels).length > 0 && makeValue && modelValue) {
      const availableModels = vehicleModels[makeValue] || [];
      if (availableModels.length > 0 && !availableModels.includes(modelValue)) {
        reset({
          make: makeValue,
          model: '',
          year: yearValue,
          trim: '',
        }, { keepDefaultValues: false });
      }
    }
  }, [makeValue, vehicleModels, modelValue, yearValue, reset]);

  // Fetch trims when make, model, and year are all selected
  useEffect(() => {
    if (makeValue && modelValue && yearValue) {
      const fetchTrims = async () => {
        try {
          setIsLoadingTrims(true);
          const fetchedTrims = await trimsAPI.getTrims(makeValue, modelValue, yearValue);
          setTrims(fetchedTrims);
        } catch (error) {
          console.error('Failed to fetch trims:', error);
          setTrims([]);
        } finally {
          setIsLoadingTrims(false);
        }
      };
      fetchTrims();
    } else {
      // Reset trim when make, model, or year changes
      resetField('trim');
      setTrims([]);
    }
  }, [makeValue, modelValue, yearValue, resetField]);

  // Fetch vehicle models on mount
  useEffect(() => {
    const fetchModels = async () => {
      try {
        setIsLoadingModels(true);
        setModelsError(null);
        const data = await modelsAPI.getModels();
        setVehicleModels(data);
        // Extract and sort makes
        const makesList = Object.keys(data).sort();
        setMakes(makesList);
      } catch (error: any) {
        console.error('Failed to fetch vehicle models:', error);
        const errorMessage = error.response?.data?.error || error.message || 'Unknown error';
        setModelsError(`Failed to load vehicle makes: ${errorMessage}. Please refresh the page.`);
        // Fallback to hardcoded makes if API fails
        setMakes(['Acura', 'Audi', 'BMW', 'Buick', 'Cadillac', 'Chevrolet', 'Chrysler', 'Dodge', 'Ford', 'Genesis', 'GMC', 'Honda', 'Hyundai', 'Infiniti', 'Jaguar', 'Jeep', 'Kia', 'Land Rover', 'Lexus', 'Lincoln', 'Mazda', 'Mercedes-Benz', 'Mini', 'Nissan', 'Porsche', 'Ram', 'Subaru', 'Tesla', 'Toyota', 'Volkswagen', 'Volvo']);
      } finally {
        setIsLoadingModels(false);
      }
    };

    fetchModels();
  }, []);

  // Load existing preferences
  useEffect(() => {
    const loadPreferences = async () => {
      try {
        const prefs = await preferencesAPI.get();
        reset({
          make: prefs.make,
          model: prefs.model,
          year: prefs.year,
          trim: prefs.trimId || 'unspecified', // Use 'unspecified' if no trimId
          zipCode: user?.zipCode || '',
        });
        // If we have make/model/year, fetch trims to populate the dropdown
        if (prefs.make && prefs.model && prefs.year) {
          try {
            setIsLoadingTrims(true);
            const fetchedTrims = await trimsAPI.getTrims(prefs.make, prefs.model, prefs.year);
            setTrims(fetchedTrims);
          } catch (error) {
            console.error('Failed to fetch trims:', error);
            setTrims([]);
          } finally {
            setIsLoadingTrims(false);
          }
        }
      } catch (error: any) {
        // If preferences don't exist, that's okay - form will start empty
        if (error.response?.status !== 404) {
          console.error('Failed to load preferences:', error);
        }
        // Still set zip code from user if available
        if (user?.zipCode) {
          reset({
            make: '',
            model: '',
            year: 2024,
            trim: '',
            zipCode: user.zipCode,
          });
        }
      } finally {
        setInitialLoading(false);
      }
    };

    loadPreferences();
  }, [reset, user]);

  // Get available models for selected make
  const availableModels = makeValue ? vehicleModels[makeValue] || [] : [];

  const onSubmit = async (data: VehicleFormData) => {
    setLoading(true);
    try {
      // Convert trim value to trimId (empty string or "unspecified" means "Unspecified")
      const trimId = data.trim === '' || data.trim === 'unspecified' ? null : data.trim;
      const zipCode = data.zipCode || '';
      await preferencesAPI.update(data.year, data.make, data.model, trimId, zipCode);
      toast.success('Vehicle settings saved successfully');
    } catch (error: any) {
      console.error('Failed to save preferences:', error);
      if (error.response?.status === 404) {
        // If preferences don't exist, create them instead
        try {
          const trimId = data.trim === '' || data.trim === 'unspecified' ? null : data.trim;
          const zipCode = data.zipCode || '';
          await preferencesAPI.create(data.year, data.make, data.model, trimId, zipCode);
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
      trim: '',
      zipCode: user?.zipCode || '',
    });
    setTrims([]);
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
              <Select onValueChange={field.onChange} value={field.value} disabled={isLoadingModels}>
                <SelectTrigger id="make" className="w-full">
                  <SelectValue placeholder={isLoadingModels ? 'Loading makes...' : 'Select a make...'} />
                </SelectTrigger>
                <SelectContent>
                  {makes.map((make) => (
                    <SelectItem key={make} value={make}>
                      {make}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            )}
          />
          {modelsError && <p className="text-sm text-destructive mt-1">{modelsError}</p>}
          <FieldError>{errors.make?.message}</FieldError>
        </Field>

        <Field>
          <FieldLabel htmlFor="model">Model</FieldLabel>
          {availableModels.length > 0 ? (
            <Controller
              name="model"
              control={control}
              render={({ field }) => (
                <Select onValueChange={field.onChange} value={field.value}>
                  <SelectTrigger id="model" className="w-full">
                    <SelectValue placeholder="Select a model..." />
                  </SelectTrigger>
                  <SelectContent>
                    {availableModels.map((model) => (
                      <SelectItem key={model} value={model}>
                        {model}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              )}
            />
          ) : (
            <Input
              id="model"
              type="text"
              {...register('model')}
              placeholder="e.g., CX-90, Camry, Model 3"
            />
          )}
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

        <Field>
          <FieldLabel htmlFor="trim">Trim</FieldLabel>
          <Controller
            name="trim"
            control={control}
            render={({ field }) => (
              <Select
                onValueChange={field.onChange}
                value={field.value}
                disabled={!makeValue || !modelValue || !yearValue || isLoadingTrims}
              >
                <SelectTrigger id="trim" className="w-full">
                  <SelectValue placeholder={
                    !makeValue || !modelValue || !yearValue 
                      ? 'Select make, model, and year first' 
                      : isLoadingTrims 
                        ? 'Loading trims...' 
                        : 'Choose a trim...'
                  } />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="unspecified">Unspecified</SelectItem>
                  {trims.map((trim) => (
                    <SelectItem key={trim.id} value={trim.id}>
                      {trim.trimName}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            )}
          />
          <FieldError>{errors.trim?.message}</FieldError>
        </Field>

        <Field>
          <div className="flex items-center gap-2">
            <FieldLabel htmlFor="zipCode">Location (Zip Code)</FieldLabel>
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








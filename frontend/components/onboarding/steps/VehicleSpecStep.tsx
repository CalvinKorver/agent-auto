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
import { VEHICLE_YEARS } from '@/lib/data/vehicleData';
import { modelsAPI, VehicleModelsResponse, trimsAPI, Trim } from '@/lib/api';

const vehicleSchema = z.object({
  make: z.string().min(1, 'Make is required'),
  model: z.string().min(1, 'Model is required'),
  year: z.number().min(2000).max(2030),
  trim: z.string().min(1, 'Trim is required'),
});

type VehicleFormData = z.infer<typeof vehicleSchema>;

interface VehicleSpecStepProps {
  onComplete: (data: { make: string; model: string; year: number; trimId?: string | null; trimName?: string | null }) => void;
  isActive: boolean;
}

export default function VehicleSpecStep({ onComplete, isActive }: VehicleSpecStepProps) {
  const [showMake] = useState(true);
  const [showModel, setShowModel] = useState(false);
  const [showYear, setShowYear] = useState(false);
  const [showTrim, setShowTrim] = useState(false);
  const [userResponses, setUserResponses] = useState<{
    make?: string;
    model?: string;
    year?: number;
    trim?: string;
  }>({});
  const [vehicleModels, setVehicleModels] = useState<VehicleModelsResponse>({});
  const [makes, setMakes] = useState<string[]>([]);
  const [isLoadingModels, setIsLoadingModels] = useState(true);
  const [modelsError, setModelsError] = useState<string | null>(null);
  const [trims, setTrims] = useState<Trim[]>([]);
  const [isLoadingTrims, setIsLoadingTrims] = useState(false);

  const {
    register,
    handleSubmit,
    watch,
    control,
    resetField,
    formState: { errors },
  } = useForm<VehicleFormData>({
    resolver: zodResolver(vehicleSchema),
    defaultValues: {
      make: '',
      model: '',
      year: 2024,
      trim: '',
    },
  });

  const makeValue = watch('make');
  const modelValue = watch('model');
  const yearValue = watch('year');
  const trimValue = watch('trim');
  const modelHasValue = (modelValue ?? '').trim().length > 0;

  // Fetch vehicle models on mount
  useEffect(() => {
    const fetchModels = async () => {
      try {
        setIsLoadingModels(true);
        setModelsError(null);
        console.log('[VehicleSpecStep] Fetching models from API...');
        const data = await modelsAPI.getModels();
        console.log('[VehicleSpecStep] Fetched data:', data);
        console.log('[VehicleSpecStep] Data type:', typeof data);
        console.log('[VehicleSpecStep] Is array?', Array.isArray(data));
        console.log('[VehicleSpecStep] Data keys:', Object.keys(data));
        console.log('[VehicleSpecStep] Data keys count:', Object.keys(data).length);
        
        if (Object.keys(data).length > 0) {
          const firstMake = Object.keys(data)[0];
          console.log('[VehicleSpecStep] First make:', firstMake);
          console.log('[VehicleSpecStep] Models for first make:', data[firstMake]);
        }
        
        setVehicleModels(data);
        // Extract and sort makes
        const makesList = Object.keys(data).sort();
        console.log('[VehicleSpecStep] Makes list:', makesList);
        console.log('[VehicleSpecStep] Makes count:', makesList.length);
        setMakes(makesList);
      } catch (error: unknown) {
        console.error('[VehicleSpecStep] Failed to fetch vehicle models:', error);
        const errorMessage = (error as { response?: { data?: { error?: string }; message?: string }; message?: string })?.response?.data?.error || 
                            (error as { message?: string })?.message || 
                            'Unknown error';
        setModelsError(`Failed to load vehicle makes: ${errorMessage}. Please refresh the page.`);
        // Fallback to hardcoded makes if API fails
        setMakes(['Acura', 'Audi', 'BMW', 'Buick', 'Cadillac', 'Chevrolet', 'Chrysler', 'Dodge', 'Ford', 'Genesis', 'GMC', 'Honda', 'Hyundai', 'Infiniti', 'Jaguar', 'Jeep', 'Kia', 'Land Rover', 'Lexus', 'Lincoln', 'Mazda', 'Mercedes-Benz', 'Mini', 'Nissan', 'Porsche', 'Ram', 'Subaru', 'Tesla', 'Toyota', 'Volkswagen', 'Volvo']);
      } finally {
        setIsLoadingModels(false);
        console.log('[VehicleSpecStep] Loading complete');
      }
    };

    fetchModels();
  }, []);

  // Get available models for selected make
  const availableModels = makeValue ? vehicleModels[makeValue] || [] : [];
  
  // Debug logging for makes state
  useEffect(() => {
    console.log('[VehicleSpecStep] Makes state updated:', {
      makesCount: makes.length,
      makes: makes,
      isLoading: isLoadingModels,
      vehicleModelsKeys: Object.keys(vehicleModels).length,
      selectedMake: makeValue,
      availableModelsCount: availableModels.length,
    });
  }, [makes, isLoadingModels, vehicleModels, makeValue, availableModels.length]);

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

  // Fetch trims when make, model, and year are all selected
  useEffect(() => {
    if (isActive && makeValue && modelValue && yearValue) {
      const fetchTrims = async () => {
        try {
          setIsLoadingTrims(true);
          const fetchedTrims = await trimsAPI.getTrims(makeValue, modelValue, yearValue);
          setTrims(fetchedTrims);
          // Show trim field after trims are loaded
          if (!showTrim) {
            setTimeout(() => setShowTrim(true), 300);
          }
        } catch (error) {
          console.error('[VehicleSpecStep] Failed to fetch trims:', error);
          setTrims([]);
        } finally {
          setIsLoadingTrims(false);
        }
      };
      fetchTrims();
    } else {
      // Reset trim when make, model, or year changes
      resetField('trim');
      setShowTrim(false);
      setTrims([]);
    }
  }, [makeValue, modelValue, yearValue, isActive, showTrim, resetField]);

  useEffect(() => {
    if (isActive && yearValue && showYear) {
      setUserResponses((prev) => ({ ...prev, year: yearValue }));
    }
  }, [yearValue, showYear, isActive]);

  useEffect(() => {
    if (isActive && trimValue && showTrim) {
      setUserResponses((prev) => ({ ...prev, trim: trimValue }));
    }
  }, [trimValue, showTrim, isActive]);

  const onSubmit = (data: VehicleFormData) => {
    // Convert trim value to trimId (empty string or "unspecified" means "Unspecified")
    const trimId = data.trim === '' || data.trim === 'unspecified' ? null : data.trim;
    // Find trim name from trims array
    const selectedTrim = trims.find(t => t.id === trimId);
    const trimName = trimId === null ? 'Unspecified' : (selectedTrim?.trimName || null);
    onComplete({
      make: data.make,
      model: data.model,
      year: data.year,
      trimId,
      trimName,
    });
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
              <div className="text-sm leading-relaxed text-foreground">
                <p className="text-sm">Great! What are you looking for?</p>
              </div>
            </div>

            <div className="pl-13">
              <Field>
                <FieldLabel htmlFor="make">Select Make</FieldLabel>
                <Controller
                  name="make"
                  control={control}
                  render={({ field }) => (
                    <Select onValueChange={field.onChange} value={field.value} disabled={isLoadingModels}>
                      <SelectTrigger id="make" className="w-full">
                        <SelectValue placeholder={isLoadingModels ? 'Loading makes...' : 'Choose a make...'} />
                      </SelectTrigger>
                      <SelectContent>
                        {makes.length === 0 && !isLoadingModels && (
                          <div className="p-2 text-sm text-muted-foreground">
                            No makes available. Debug: makes={makes.length}, vehicleModels={Object.keys(vehicleModels).length}
                          </div>
                        )}
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
            </div>

            {showModel && (
              <div className="pl-13 mt-4 animate-in fade-in slide-in-from-bottom-4 duration-500">
                <Field>
                  <FieldLabel htmlFor="model">Select Model</FieldLabel>
                  {availableModels.length > 0 ? (
                    <Controller
                      name="model"
                      control={control}
                      render={({ field }) => (
                        <Select onValueChange={field.onChange} value={field.value}>
                          <SelectTrigger id="model" className="w-full">
                            <SelectValue placeholder="Choose a model..." />
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

                {showTrim && (
                  <div className="mt-4 animate-in fade-in slide-in-from-bottom-4 duration-500">
                    <Field>
                      <FieldLabel htmlFor="trim">Select Trim</FieldLabel>
                      <Controller
                        name="trim"
                        control={control}
                        render={({ field }) => (
                          <Select 
                            onValueChange={field.onChange} 
                            value={field.value}
                            disabled={isLoadingTrims}
                          >
                            <SelectTrigger id="trim" className="w-full">
                              <SelectValue placeholder={isLoadingTrims ? 'Loading trims...' : 'Choose a trim...'} />
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
        {showTrim && trimValue && (
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



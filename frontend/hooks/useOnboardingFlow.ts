import { useState } from 'react';

export interface OnboardingState {
  currentStep: number;
  journeyType: 'exploring' | 'high-intent' | null;
  vehicleData: {
    make: string;
    model: string;
    year: number | null;
  };
}

export function useOnboardingFlow() {
  const [state, setState] = useState<OnboardingState>({
    currentStep: 0,
    journeyType: null,
    vehicleData: {
      make: '',
      model: '',
      year: null,
    },
  });

  const nextStep = () => {
    setState((prev) => ({ ...prev, currentStep: prev.currentStep + 1 }));
  };

  const previousStep = () => {
    setState((prev) => ({ ...prev, currentStep: Math.max(0, prev.currentStep - 1) }));
  };

  const setJourneyType = (type: 'exploring' | 'high-intent') => {
    setState((prev) => ({ ...prev, journeyType: type }));
  };

  const setVehicleData = (data: Partial<OnboardingState['vehicleData']>) => {
    setState((prev) => ({
      ...prev,
      vehicleData: { ...prev.vehicleData, ...data },
    }));
  };

  return {
    state,
    nextStep,
    previousStep,
    setJourneyType,
    setVehicleData,
  };
}

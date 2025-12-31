'use client';

import { useState, useEffect, useRef } from 'react';
import { preferencesAPI } from '@/lib/api';
import OnboardingSidebar from './OnboardingSidebar';
import ChatInterface from './ChatInterface';
import WelcomeStep from './steps/WelcomeStep';
import JourneyStateStep from './steps/JourneyStateStep';
import VehicleSpecStep from './steps/VehicleSpecStep';
import ZipCodeStep from './steps/ZipCodeStep';
import LaunchPadStep from './steps/LaunchPadStep';

interface OnboardingContainerProps {
  onSuccess: () => Promise<void>;
}

export default function OnboardingContainer({ onSuccess }: OnboardingContainerProps) {
  const [currentStep, setCurrentStep] = useState(0);
  const [showSteps, setShowSteps] = useState({
    welcome: true,
    journey: false,
    vehicle: false,
    zipCode: false,
    launchpad: false,
  });
  const [vehicleData, setVehicleData] = useState<{
    make: string;
    model: string;
    year: number;
    trimId?: string | null;
    trimName?: string | null;
  } | null>(null);
  const [zipCode, setZipCode] = useState<string | null>(null);
  const chatEndRef = useRef<HTMLDivElement>(null);

  // Auto-scroll to bottom when new steps appear
  useEffect(() => {
    chatEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [showSteps]);

  const handleWelcomeNext = () => {
    setCurrentStep(1);
    setTimeout(() => {
      setShowSteps((prev) => ({ ...prev, journey: true }));
    }, 300);
  };

  const handleJourneySelect = () => {
    setCurrentStep(2);
    setTimeout(() => {
      setShowSteps((prev) => ({ ...prev, vehicle: true }));
    }, 300);
  };

  const handleVehicleComplete = (data: { make: string; model: string; year: number; trimId?: string | null; trimName?: string | null }) => {
    setVehicleData(data);
    setCurrentStep(3);
    setTimeout(() => {
      setShowSteps((prev) => ({ ...prev, zipCode: true }));
    }, 300);
  };

  const handleZipCodeComplete = (zip: string) => {
    setZipCode(zip);
    setCurrentStep(4);
    setTimeout(() => {
      setShowSteps((prev) => ({ ...prev, launchpad: true }));
    }, 300);
  };

  const handleLaunchPadComplete = async () => {
    if (!vehicleData) return;

    await preferencesAPI.create(vehicleData.year, vehicleData.make, vehicleData.model, vehicleData.trimId, zipCode || '');
    await onSuccess();
  };

  // Map step to sidebar progress (0-2)
  const getSidebarStep = () => {
    if (currentStep === 0 || currentStep === 1) return 0; // About you
    if (currentStep === 2 || currentStep === 3) return 1; // Car
    return 2; // Checkout
  };

  return (
    <div className="min-h-screen flex flex-col lg:flex-row bg-background">
      <OnboardingSidebar currentStep={getSidebarStep()} />
      <ChatInterface>
        <div className="space-y-8">
          {showSteps.welcome && (
            <WelcomeStep onNext={handleWelcomeNext} isActive={currentStep === 0} />
          )}
          {showSteps.journey && (
            <JourneyStateStep onSelectJourney={handleJourneySelect} isActive={currentStep === 1} />
          )}
          {showSteps.vehicle && (
            <VehicleSpecStep onComplete={handleVehicleComplete} isActive={currentStep === 2} />
          )}
          {showSteps.zipCode && (
            <ZipCodeStep onComplete={handleZipCodeComplete} isActive={currentStep === 3} />
          )}
          {showSteps.launchpad && vehicleData && (
            <LaunchPadStep vehicleData={vehicleData} onComplete={handleLaunchPadComplete} />
          )}
          <div ref={chatEndRef} />
        </div>
      </ChatInterface>
    </div>
  );
}

interface OnboardingSidebarProps {
  currentStep: number;
}

export default function OnboardingSidebar({ currentStep }: OnboardingSidebarProps) {
  const steps = [
    { label: 'Welcome!', step: 0 },
    { label: 'Quick Questions', step: 1 },
    { label: 'Lets go', step: 2 },
  ];

  return (
    <div className="w-full lg:w-64 bg-card border-r border-border p-6 flex lg:flex-col gap-6">
      <div className="hidden lg:block mb-8">
        <h2 className="text-xl font-bold text-card-foreground">Lolo AI</h2>
      </div>

      <nav className="flex lg:flex-col gap-4 lg:gap-6 flex-1 lg:flex-none">
        {steps.map((step, index) => {
          const isActive = currentStep >= step.step;
          const isCurrent = currentStep === step.step;

          return (
            <div key={step.label} className="flex items-center gap-3">
              <div
                className={`w-8 h-8 rounded-full flex items-center justify-center text-sm font-semibold transition-all ${
                  isActive
                    ? 'bg-primary text-primary-foreground'
                    : 'border-2 border-border text-muted-foreground'
                } ${isCurrent ? 'ring-4 ring-primary/20' : ''}`}
              >
                {isActive ? '✓' : index + 1}
              </div>
              <div
                className={`text-sm font-medium transition-colors ${
                  isActive ? 'text-foreground' : 'text-muted-foreground'
                }`}
              >
                {step.label}
              </div>
            </div>
          );
        })}
      </nav>
    </div>
  );
}

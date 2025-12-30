'use client';

import Link from 'next/link';
import Image from 'next/image';
import { useTheme } from 'next-themes';
import { useEffect, useState } from 'react';
import { Button } from '@/components/ui/button';
import { Car, MessageSquare, Bell, TrendingUp } from 'lucide-react';

export default function LandingPage() {
  const { resolvedTheme } = useTheme();
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  // Use dark logo in light mode, light logo in dark mode
  const logoSrc = mounted && resolvedTheme === 'light' 
    ? '/logo-dark-v2.png' 
    : '/logo-light-v2.png';

  // Check if redirect to waiting page is enabled
  const redirectToWaiting = process.env.NEXT_PUBLIC_REDIRECT_TO_WAITING === 'True';
  const loginHref = redirectToWaiting ? '/waiting' : '/login';
  const registerHref = redirectToWaiting ? '/waiting' : '/register';

  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-50 to-white dark:from-gray-900 dark:to-gray-800">
      {/* Hero Section */}
      <header className="container mx-auto px-4 py-6">
        <nav className="flex justify-between items-center">
          <div className="flex items-center">
            <Image
              src={logoSrc}
              alt="Otto"
              width={200}
              height={60}
              className="h-12 w-auto"
            />
          </div>
          <div className="flex gap-3">
            <Link href={loginHref}>
              <Button variant="ghost">Sign In</Button>
            </Link>
            <Link href={registerHref}>
              <Button>Get Started</Button>
            </Link>
          </div>
        </nav>
      </header>

      {/* Hero Content */}
      <section className="container mx-auto px-4 py-20 text-center">
        <h1 className="text-5xl font-bold tracking-tight mb-6 max-w-3xl mx-auto text-foreground">
          Your AI-Powered Car Buying Assistant
        </h1>
        <p className="text-xl text-muted-foreground mb-8 max-w-2xl mx-auto">
          Stay in control while AI helps you communicate with dealers and negotiate the best price.
          Get the car you want without the stress.
        </p>
        <div className="flex gap-4 justify-center">
          <Link href={registerHref}>
            <Button size="lg">Start Your Search</Button>
          </Link>
          <Link href="#features">
            <Button size="lg" variant="outline">Learn More</Button>
          </Link>
        </div>
      </section>

      {/* Features Section */}
      <section id="features" className="container mx-auto px-4 py-20">
        <h2 className="text-3xl font-bold text-center mb-12 text-foreground">How It Works</h2>
        <div className="grid md:grid-cols-3 gap-8 max-w-5xl mx-auto">
          <FeatureCard
            icon={<MessageSquare className="size-10 text-primary" />}
            title="Set Your Preferences"
            description="Tell us what you're looking for - make, model, budget, and must-haves. You're always in control."
          />
          <FeatureCard
            icon={<Car className="size-10 text-primary" />}
            title="AI-Assisted Negotiation"
            description="When you're ready, use our AI agent to help craft messages to dealers and negotiate the lowest price possible."
          />
          <FeatureCard
            icon={<Bell className="size-10 text-primary" />}
            title="Track & Compare"
            description="Manage all dealer conversations in one place. Compare offers side-by-side - Otto will keep track of all offers."
          />
        </div>
      </section>

      {/* Benefits Section */}
      <section className="bg-gray-100 dark:bg-gray-800 py-20">
        <div className="container mx-auto px-4">
          <h2 className="text-3xl font-bold text-center mb-12 text-foreground">Why Choose AgentAuto?</h2>
          <div className="grid md:grid-cols-2 gap-8 max-w-4xl mx-auto">
            <BenefitCard
              icon={<TrendingUp className="size-8 text-primary" />}
              title="Get the Best Price"
              description="Harness AI to help negotiate effectively with dealers. Let technology level the playing field and get you competitive pricing."
            />
            <BenefitCard
              icon={<MessageSquare className="size-8 text-primary" />}
              title="You Stay in Control"
              description="AgentAuto assists you, never replaces you. Review AI suggestions, approve messages, and make all final decisions."
            />
            <BenefitCard
              icon={<Car className="size-8 text-primary" />}
              title="Centralized Hub"
              description="Manage all your dealer conversations in one organized inbox. No more juggling emails and texts."
            />
            <BenefitCard
              icon={<Bell className="size-8 text-primary" />}
              title="No Pressure"
              description="Communicate on your schedule. Take your time reviewing offers without high-pressure sales tactics."
            />
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="container mx-auto px-4 py-20 text-center">
        <h2 className="text-4xl font-bold mb-6 text-foreground">Ready to Find Your Next Car?</h2>
        <p className="text-xl text-muted-foreground mb-8 max-w-2xl mx-auto">
          Join thousands of satisfied customers who found their perfect vehicle through AgentAuto.
        </p>
        <Link href={registerHref}>
          <Button size="lg">Get Started Today</Button>
        </Link>
      </section>

      {/* Footer */}
      <footer className="border-t py-8">
        <div className="container mx-auto px-4 text-center text-muted-foreground">
          <p>&copy; 2025 AgentAuto. All rights reserved.</p>
        </div>
      </footer>
    </div>
  );
}

function FeatureCard({ icon, title, description }: { icon: React.ReactNode; title: string; description: string }) {
  return (
    <div className="text-center p-6">
      <div className="flex justify-center mb-4">{icon}</div>
      <h3 className="text-xl font-semibold mb-3 text-foreground">{title}</h3>
      <p className="text-muted-foreground">{description}</p>
    </div>
  );
}

function BenefitCard({ icon, title, description }: { icon: React.ReactNode; title: string; description: string }) {
  return (
    <div className="flex gap-4 p-6 bg-white dark:bg-gray-900 rounded-lg shadow-sm">
      <div className="flex-shrink-0">{icon}</div>
      <div>
        <h3 className="text-lg font-semibold mb-2 text-foreground">{title}</h3>
        <p className="text-muted-foreground">{description}</p>
      </div>
    </div>
  );
}

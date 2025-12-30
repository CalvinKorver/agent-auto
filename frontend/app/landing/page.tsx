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
      <section className="container mx-auto px-4 pt-20 pb-12 text-center space-y-8">
        <h1 className="text-5xl font-bold tracking-tight max-w-3xl mx-auto text-foreground">
          Haggling for a car sucks.
        </h1>
        <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
          Dealerships have sophisticated tools. You don&apos;t - and it&apos;s costing hundreds if not thousands of dollars when you buy a new car.
        </p>
        <p className="text-2xl font-semibold max-w-2xl mx-auto text-foreground pt-8 mb-4 md:mb-16">
          That&apos;s why we built Otto.
        </p>
        <div className="flex justify-center pt-4">
          <Image
            src="/car_art.png"
            alt="Car illustration"
            width={200}
            height={100}
            className="w-36 h-auto md:w-auto md:h-auto max-w-full"
          />
        </div>
      </section>

      {/* Features Section */}
      <section id="features" className="container mx-auto px-4 py-20">
        <h2 className="text-3xl font-bold text-center mb-6 max-w-3xl mx-auto text-foreground">
          Otto is your personal negotiation advocate
        </h2>
        <p className="text-xl text-muted-foreground text-center mb-12 max-w-3xl mx-auto">
          It levels the playing field by harnessing decades of professional car-buying expertise and market data to ensure you aren&apos;t overpaying.
        </p>
        <div className="grid md:grid-cols-3 gap-8 max-w-5xl mx-auto">
          <FeatureCard
            icon={<MessageSquare className="size-10 text-primary" />}
            title="Expert Playbooks"
            description="We&apos;ve distilled expert playbooks into one platform. Otto knows exactly how to navigate holdbacks and hidden fees, focusing on the only metric that matters: a fair &apos;Out the Door&apos; (OTD) price."
          />
          <FeatureCard
            icon={<Car className="size-10 text-primary" />}
            title="Save Your Time"
            description="Instead of managing dozens of high-pressure texts and emails while you&apos;re at work, Otto helps to organize the grueling back-and-forth communication for you, making responding easier."
          />
          <FeatureCard
            icon={<Bell className="size-10 text-primary" />}
            title="You&apos;re Always in Control"
            description="You are always the human in the loop. Otto drafts expert responses based on successful negotiation tactics, but you stay in total control—reviewing and approving every move before it is sent."
          />
        </div>
      </section>

      {/* Benefits Section */}
      <section className="bg-gray-100 dark:bg-gray-800 py-20">
        <div className="container mx-auto px-4">
          <h2 className="text-3xl font-bold text-center mb-6 max-w-3xl mx-auto text-foreground">
            The difference between us & other AI tools?
          </h2>
          <p className="text-xl text-muted-foreground text-center mb-12 max-w-3xl mx-auto">
            You are always the human in the loop. Otto drafts expert responses based on successful negotiation tactics, but you stay in total control—reviewing and approving every move before it is sent.
          </p>
          <div className="grid md:grid-cols-2 gap-8 max-w-4xl mx-auto">
            <BenefitCard
              icon={<TrendingUp className="size-8 text-primary" />}
              title="Level the Playing Field"
              description="While dealerships use sophisticated CRM tools, Otto gives you access to decades of professional car-buying expertise and market data to ensure you aren&apos;t overpaying."
            />
            <BenefitCard
              icon={<MessageSquare className="size-8 text-primary" />}
              title="Focus on What Matters"
              description="Otto knows exactly how to navigate holdbacks and hidden fees, focusing on the only metric that matters: a fair &apos;Out the Door&apos; (OTD) price."
            />
            <BenefitCard
              icon={<Car className="size-8 text-primary" />}
              title="Get Your Time Back"
              description="Instead of managing dozens of high-pressure texts and emails while you&apos;re at work, Otto organizes the grueling back-and-forth communication for you."
            />
            <BenefitCard
              icon={<Bell className="size-8 text-primary" />}
              title="Total Control"
              description="Review and approve every message before it&apos;s sent. You make all final decisions—Otto is your advocate, not your replacement."
            />
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="container mx-auto px-4 py-20 text-center">
        <h2 className="text-4xl font-bold mb-6 text-foreground">Ready to level the playing field?</h2>
        <p className="text-xl text-muted-foreground mb-8 max-w-2xl mx-auto">
          Stop overpaying for your next car. Let Otto be your personal negotiation advocate and get the fair price you deserve.
        </p>
        <Link href={registerHref}>
          <Button size="lg">Get Started</Button>
        </Link>
      </section>

      {/* Footer */}
      <footer className="py-4 px-4">
        <div className="max-w-2xl mx-auto text-center text-sm text-muted-foreground">
          <p>&copy; 2025 OttoAI. All rights reserved.</p>
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

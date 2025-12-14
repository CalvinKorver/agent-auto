'use client';

import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { preferencesAPI } from '@/lib/api';

const preferencesSchema = z.object({
  year: z.number().min(2000).max(2030),
  make: z.string().min(1, 'Make is required'),
  model: z.string().min(1, 'Model is required'),
});

type PreferencesFormData = z.infer<typeof preferencesSchema>;

interface PreferencesFormProps {
  onSuccess: () => void;
}

export default function PreferencesForm({ onSuccess }: PreferencesFormProps) {
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<PreferencesFormData>({
    resolver: zodResolver(preferencesSchema),
    defaultValues: {
      year: 2024,
    },
  });

  const onSubmit = async (data: PreferencesFormData) => {
    try {
      setLoading(true);
      setError('');
      await preferencesAPI.create(data.year, data.make, data.model);
      onSuccess();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to save preferences');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 px-4">
      <div className="max-w-md w-full bg-white rounded-lg shadow-md p-8">
        <div className="text-center mb-8">
          <div className="text-4xl mb-4">🚗</div>
          <h1 className="text-2xl font-bold text-gray-900 mb-2">
            Let&apos;s define your target vehicle
          </h1>
          <p className="text-sm text-gray-600">
            These preferences will be locked for this buying session and guide all negotiations
          </p>
        </div>

        {error && (
          <div className="mb-4 p-3 bg-red-50 border border-red-200 text-red-700 rounded-md text-sm">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
          <div>
            <label htmlFor="year" className="block text-sm font-medium text-gray-700 mb-1">
              Year
            </label>
            <input
              id="year"
              type="number"
              {...register('year', { valueAsNumber: true })}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="2024"
            />
            {errors.year && (
              <p className="mt-1 text-sm text-red-600">{errors.year.message}</p>
            )}
          </div>

          <div>
            <label htmlFor="make" className="block text-sm font-medium text-gray-700 mb-1">
              Make
            </label>
            <input
              id="make"
              type="text"
              {...register('make')}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Mazda"
            />
            {errors.make && (
              <p className="mt-1 text-sm text-red-600">{errors.make.message}</p>
            )}
          </div>

          <div>
            <label htmlFor="model" className="block text-sm font-medium text-gray-700 mb-1">
              Model
            </label>
            <input
              id="model"
              type="text"
              {...register('model')}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="CX-90"
            />
            {errors.model && (
              <p className="mt-1 text-sm text-red-600">{errors.model.message}</p>
            )}
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed font-medium"
          >
            {loading ? 'SAVING...' : 'START NEGOTIATING'}
          </button>
        </form>
      </div>
    </div>
  );
}

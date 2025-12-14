import axios from 'axios';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

export const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add auth token to requests if available
api.interceptors.request.use((config) => {
  if (typeof window !== 'undefined') {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
  }
  return config;
});

// Types
export interface User {
  id: string;
  email: string;
  createdAt: string;
  preferences?: UserPreferences;
}

export interface UserPreferences {
  year: number;
  make: string;
  model: string;
}

export interface AuthResponse {
  user: User;
  token: string;
}

export interface ErrorResponse {
  error: string;
}

// Auth API
export const authAPI = {
  register: async (email: string, password: string): Promise<AuthResponse> => {
    const response = await api.post<AuthResponse>('/auth/register', { email, password });
    return response.data;
  },

  login: async (email: string, password: string): Promise<AuthResponse> => {
    const response = await api.post<AuthResponse>('/auth/login', { email, password });
    return response.data;
  },

  me: async (): Promise<User> => {
    const response = await api.get<User>('/auth/me');
    return response.data;
  },
};

// Preferences API
export const preferencesAPI = {
  get: async (): Promise<UserPreferences> => {
    const response = await api.get<UserPreferences>('/preferences');
    return response.data;
  },

  create: async (year: number, make: string, model: string): Promise<UserPreferences> => {
    const response = await api.post<UserPreferences>('/preferences', { year, make, model });
    return response.data;
  },
};

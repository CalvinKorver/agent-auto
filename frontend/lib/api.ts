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
  inboxEmail?: string;
  preferences?: UserPreferences;
  gmailConnected?: boolean;
  gmailEmail?: string;
}

export interface UserPreferences {
  year: number;
  make: string;
  model: string;
}

export interface Thread {
  id: string;
  sellerName: string;
  sellerType: 'private' | 'dealership' | 'other';
  createdAt: string;
  lastMessageAt?: string;
  messageCount: number;
}

export interface Message {
  id: string;
  threadId: string;
  sender: 'user' | 'agent' | 'seller';
  content: string;
  timestamp: string;
  externalMessageId?: string;
  senderEmail?: string;
  subject?: string;
}

export interface InboxMessage {
  id: string;
  sender: 'user' | 'agent' | 'seller';
  senderEmail: string;
  subject: string;
  content: string;
  timestamp: string;
  externalMessageId?: string;
}

export interface TrackedOffer {
  id: string;
  threadId: string;
  messageId?: string;
  offerText: string;
  trackedAt: string;
  sellerName?: string;
  threadType?: string;
}

export interface MessageResponse {
  messages: Message[];
  total: number;
  hasMore: boolean;
}

export interface CreateThreadRequest {
  sellerName: string;
  sellerType: 'private' | 'dealership' | 'other';
}

export interface CreateMessageRequest {
  content: string;
  sender: 'user' | 'seller';
}

export interface CreateUserMessageResponse {
  userMessage: Message;
  agentMessage: Message;
}

export interface CreateSellerMessageResponse {
  sellerMessage: Message;
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

  logout: async (): Promise<void> => {
    await api.post('/auth/logout');
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

  update: async (year: number, make: string, model: string): Promise<UserPreferences> => {
    const response = await api.put<UserPreferences>('/preferences', { year, make, model });
    return response.data;
  },
};

// Thread API
export const threadAPI = {
  getAll: async (): Promise<Thread[]> => {
    const response = await api.get<{ threads: Thread[] }>('/threads');
    return response.data.threads;
  },

  getById: async (threadId: string): Promise<Thread> => {
    const response = await api.get<Thread>(`/threads/${threadId}`);
    return response.data;
  },

  create: async (data: CreateThreadRequest): Promise<Thread> => {
    const response = await api.post<Thread>('/threads', data);
    return response.data;
  },

  archive: async (threadId: string): Promise<void> => {
    await api.delete(`/threads/${threadId}`);
  },
};

// Message API
export const messageAPI = {
  getThreadMessages: async (threadId: string, limit = 50, offset = 0): Promise<MessageResponse> => {
    const response = await api.get<MessageResponse>(`/threads/${threadId}/messages`, {
      params: { limit, offset },
    });
    return response.data;
  },

  createMessage: async (
    threadId: string,
    data: CreateMessageRequest
  ): Promise<CreateUserMessageResponse | CreateSellerMessageResponse> => {
    const response = await api.post<CreateUserMessageResponse | CreateSellerMessageResponse>(
      `/threads/${threadId}/messages`,
      data
    );
    return response.data;
  },

  getInboxMessages: async (limit = 50, offset = 0): Promise<{ messages: InboxMessage[]; total: number; hasMore: boolean }> => {
    const response = await api.get<{ messages: InboxMessage[]; total: number; hasMore: boolean }>('/inbox/messages', {
      params: { limit, offset },
    });
    return response.data;
  },

  assignInboxMessageToThread: async (messageId: string, threadId: string): Promise<void> => {
    await api.put(`/inbox/messages/${messageId}/assign`, { threadId });
  },

  archiveInboxMessage: async (messageId: string): Promise<void> => {
    await api.delete(`/inbox/messages/${messageId}`);
  },
};

// Offer API
export const offerAPI = {
  createOffer: async (threadId: string, offerText: string, messageId?: string): Promise<TrackedOffer> => {
    const response = await api.post<TrackedOffer>(`/threads/${threadId}/offers`, {
      offerText,
      messageId: messageId || null,
    });
    return response.data;
  },

  getAllOffers: async (): Promise<TrackedOffer[]> => {
    const response = await api.get<{ offers: TrackedOffer[] }>('/offers');
    return response.data.offers;
  },
};

// Gmail API
export const gmailAPI = {
  getAuthUrl: async (): Promise<{ authUrl: string; state: string }> => {
    const response = await api.get<{ authUrl: string; state: string }>('/gmail/connect');
    return response.data;
  },

  getStatus: async (): Promise<{ connected: boolean; gmailEmail?: string }> => {
    const response = await api.get<{ connected: boolean; gmailEmail?: string }>('/gmail/status');
    return response.data;
  },

  disconnect: async (): Promise<void> => {
    await api.post('/gmail/disconnect');
  },

  replyViaGmail: async (messageId: string, content: string): Promise<void> => {
    await api.post(`/messages/${messageId}/reply-via-gmail`, { content });
  },
};

import axios from 'axios';
import type { 
  User, 
  URL, 
  Analytics, 
  URLStats, 
  CreateURLRequest, 
  LoginRequest, 
  RegisterRequest, 
  AuthResponse, 
  ApiResponse
} from '@/types';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

const api = axios.create({
  baseURL: `${API_BASE_URL}/api/v1`,
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json',
  },
});

// OAuth uses session-based authentication, no need for token interceptors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Clear any stored auth data, but don't redirect automatically
      localStorage.removeItem('token');
      localStorage.removeItem('user');
    }
    return Promise.reject(error);
  }
);

export const authApi = {
  // Get Google OAuth login URL
  getOAuthLoginUrl: (): Promise<ApiResponse<{ auth_url: string }>> =>
    api.get('/auth/login').then(res => res.data),
  
  // Handle OAuth callback (called automatically by backend)
  handleOAuthCallback: (code: string, state: string): Promise<ApiResponse<User>> =>
    api.get(`/auth/callback?code=${encodeURIComponent(code)}&state=${encodeURIComponent(state)}`).then(res => res.data),
  
  logout: (): Promise<ApiResponse> =>
    api.post('/auth/logout').then(res => res.data),
  
  getProfile: (): Promise<ApiResponse<User>> =>
    api.get('/auth/profile').then(res => res.data),
};

export const urlApi = {
  createURL: (data: CreateURLRequest): Promise<ApiResponse<URL>> =>
    api.post('/urls', data).then(res => res.data),
  
  getUserURLs: (limit = 10, offset = 0): Promise<ApiResponse<{
    urls: URL[];
    total: number;
    limit: number;
    offset: number;
  }>> =>
    api.get(`/urls?limit=${limit}&offset=${offset}`).then(res => res.data),
  
  updateURL: (id: number, data: Partial<CreateURLRequest>): Promise<ApiResponse<URL>> =>
    api.put(`/urls/${id}`, data).then(res => res.data),
  
  deleteURL: (id: number): Promise<ApiResponse> =>
    api.delete(`/urls/${id}`).then(res => res.data),
  
  getURLInfo: (shortCode: string): Promise<ApiResponse<URL>> =>
    api.get(`/urls/${shortCode}/info`).then(res => res.data),
  
  getURLAnalytics: (id: number): Promise<ApiResponse<{
    analytics: Analytics[];
    stats: URLStats;
  }>> =>
    api.get(`/urls/${id}/analytics`).then(res => res.data),
};

export const publicApi = {
  redirect: (shortCode: string): string =>
    `${API_BASE_URL}/${shortCode}`,
};

export default api;
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

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export const authApi = {
  register: (data: RegisterRequest): Promise<AuthResponse> =>
    api.post('/auth/register', data).then(res => res.data),
  
  login: (data: LoginRequest): Promise<AuthResponse> =>
    api.post('/auth/login', data).then(res => res.data),
  
  logout: (): Promise<ApiResponse> =>
    api.post('/auth/logout').then(res => res.data),
  
  getProfile: (): Promise<ApiResponse<User>> =>
    api.get('/auth/profile').then(res => res.data),
  
  refreshToken: (): Promise<ApiResponse<{ token: string }>> =>
    api.post('/auth/refresh').then(res => res.data),
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
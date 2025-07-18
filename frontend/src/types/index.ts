export interface User {
  id: number;
  email: string;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface URL {
  id: number;
  original_url: string;
  short_code: string;
  custom_alias?: string;
  user_id?: number;
  user?: User;
  title?: string;
  description?: string;
  expires_at?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Analytics {
  id: number;
  url_id: number;
  ip_address: string;
  user_agent: string;
  referrer: string;
  country?: string;
  city?: string;
  device?: string;
  os?: string;
  browser?: string;
  clicked_at: string;
  created_at: string;
}

export interface URLStats {
  url_id: number;
  total_clicks: number;
  unique_clicks: number;
  last_clicked?: string;
}

export interface CreateURLRequest {
  original_url: string;
  custom_alias?: string;
  title?: string;
  description?: string;
  expires_at?: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  name: string;
}

export interface AuthResponse {
  success: boolean;
  data: {
    user: User;
    token: string;
  };
  message: string;
}

export interface ApiResponse<T = any> {
  success: boolean;
  data?: T;
  message?: string;
}

export interface ApiError {
  error: string;
  message?: string;
}

export interface PaginatedResponse<T> {
  success: boolean;
  data: {
    items: T[];
    total: number;
    limit: number;
    offset: number;
  };
}
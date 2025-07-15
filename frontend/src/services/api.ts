import axios from 'axios'

const API_BASE_URL = (import.meta as any).env.VITE_API_URL || 'https://s.iafri.com'

// Create axios instance
const api = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
})

// Add auth token to requests
api.interceptors.request.use((config) => {
    const token = localStorage.getItem('auth_token')
    if (token) {
        config.headers.Authorization = `Bearer ${token}`
    }
    return config
})

// Handle auth errors
api.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response?.status === 401) {
            // Clear auth data and redirect to login
            localStorage.removeItem('auth_token')
            localStorage.removeItem('auth_user')

            // Use React Router navigation instead of hard refresh
            // Only redirect if not already on login page
            if (window.location.pathname !== '/login') {
                window.location.href = '/login'
            }
        }
        return Promise.reject(error)
    }
)

// Types
export interface URL {
    id: number
    short_code: string
    original_url: string
    user_id: number
    created_at: string
    updated_at: string
    click_count: number
    is_active: boolean
    expires_at?: string
    user_agent?: string
    ip_address?: string
}

export interface CreateURLRequest {
    url: string
    custom_code?: string
    expires_at?: string
}

export interface UpdateURLRequest {
    original_url?: string
    is_active?: boolean
    expires_at?: string
}

export interface URLAnalytics {
    total_clicks: number
    unique_clicks: number
    clicks_today: number
    clicks_this_week: number
    top_countries: Array<{ country: string; clicks: number }>
    top_referrers: Array<{ referrer: string; clicks: number }>
}

export interface LoginRequest {
    email: string
    password: string
}

export interface RegisterRequest {
    email: string
    password: string
    first_name: string
    last_name: string
}

// Auth API
export const authAPI = {
    login: (data: LoginRequest) => api.post('/api/v1/auth/login', data),
    register: (data: RegisterRequest) => api.post('/api/v1/auth/register', data),
    getProfile: () => api.get('/api/v1/profile'),
    updateProfile: (data: any) => api.put('/api/v1/profile', data),
    changePassword: (data: any) => api.post('/api/v1/profile/change-password', data),
    logout: () => api.post('/api/v1/auth/logout'),
    refreshToken: () => api.post('/api/v1/auth/refresh'),
}

// URLs API
export const urlsAPI = {
    create: (data: CreateURLRequest) => api.post('/api/v1/urls', data),
    getAll: (params?: { limit?: number; offset?: number }) =>
        api.get('/api/v1/urls', { params }),
    getByCode: (shortCode: string) => api.get(`/api/v1/urls/${shortCode}`),
    update: (shortCode: string, data: UpdateURLRequest) =>
        api.put(`/api/v1/urls/${shortCode}`, data),
    delete: (shortCode: string) => api.delete(`/api/v1/urls/${shortCode}`),
    getAnalytics: (shortCode: string, days?: number) =>
        api.get(`/api/v1/urls/${shortCode}/analytics`, { params: { days } }),
    getQRCode: (shortCode: string) => api.get(`/api/v1/urls/${shortCode}/qr`, { responseType: 'blob' }),
}

export default api 
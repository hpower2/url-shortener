import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react'
import { authAPI } from '../services/api'
import toast from 'react-hot-toast'

interface User {
    id: number
    email: string
    first_name: string
    last_name: string
    is_active: boolean
    created_at: string
}

interface AuthContextType {
    user: User | null
    token: string | null
    login: (email: string, password: string) => Promise<void>
    register: (data: RegisterData) => Promise<void>
    logout: () => void
    loading: boolean
    updateUser: (data: Partial<User>) => void
}

interface RegisterData {
    email: string
    password: string
    first_name: string
    last_name: string
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

interface AuthProviderProps {
    children: ReactNode
}

export function AuthProvider({ children }: AuthProviderProps) {
    const [user, setUser] = useState<User | null>(null)
    const [token, setToken] = useState<string | null>(null)
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        // Check for stored auth data on app start
        const storedToken = localStorage.getItem('auth_token')
        const storedUser = localStorage.getItem('auth_user')

        if (storedToken && storedUser) {
            try {
                setToken(storedToken)
                setUser(JSON.parse(storedUser))
            } catch (error) {
                // Clear invalid stored data
                localStorage.removeItem('auth_token')
                localStorage.removeItem('auth_user')
                toast.error('Session expired. Please log in again.')
            }
        }
        setLoading(false)
    }, [])

    const login = async (email: string, password: string) => {
        try {
            const response = await authAPI.login({ email, password })
            const { user: userData, token: authToken } = response.data

            setUser(userData)
            setToken(authToken)

            // Store in localStorage
            localStorage.setItem('auth_token', authToken)
            localStorage.setItem('auth_user', JSON.stringify(userData))

            toast.success(`Welcome back, ${userData.first_name}!`)
        } catch (error: any) {
            const message = error.response?.data?.error || 'Login failed'
            toast.error(message)
            throw error
        }
    }

    const register = async (data: RegisterData) => {
        try {
            const response = await authAPI.register(data)
            const { user: userData, token: authToken } = response.data

            setUser(userData)
            setToken(authToken)

            // Store in localStorage
            localStorage.setItem('auth_token', authToken)
            localStorage.setItem('auth_user', JSON.stringify(userData))

            toast.success(`Welcome, ${userData.first_name}! Your account has been created.`)
        } catch (error: any) {
            const message = error.response?.data?.error || 'Registration failed'
            toast.error(message)
            throw error
        }
    }

    const logout = () => {
        setUser(null)
        setToken(null)
        localStorage.removeItem('auth_token')
        localStorage.removeItem('auth_user')
        toast.success('Logged out successfully')
    }

    const updateUser = (data: Partial<User>) => {
        if (user) {
            const updatedUser = { ...user, ...data }
            setUser(updatedUser)
            localStorage.setItem('auth_user', JSON.stringify(updatedUser))
        }
    }

    const value = {
        user,
        token,
        login,
        register,
        logout,
        loading,
        updateUser,
    }

    return (
        <AuthContext.Provider value={value}>
            {children}
        </AuthContext.Provider>
    )
}

export function useAuth() {
    const context = useContext(AuthContext)
    if (context === undefined) {
        throw new Error('useAuth must be used within an AuthProvider')
    }
    return context
} 
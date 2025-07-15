import React from 'react'
import { Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider, useAuth } from './contexts/AuthContext'
import Layout from './components/Layout'
import Login from './pages/Login'
import Register from './pages/Register'
import Dashboard from './pages/Dashboard'
import Analytics from './pages/Analytics'
import Profile from './pages/Profile'
import Home from './pages/Home'
import NotFound from './pages/NotFound'
import ErrorExpired from './pages/ErrorExpired'
import ErrorInactive from './pages/ErrorInactive'
import ErrorNotFound from './pages/ErrorNotFound'
import ErrorServer from './pages/ErrorServer'
import PublicRoute from './components/PublicRoute'
import PrivateRoute from './components/PrivateRoute'

function AppRoutes() {
    const { user } = useAuth()

    return (
        <Layout>
            <Routes>
                {/* Public routes */}
                <Route path="/" element={<Home />} />

                {/* Auth routes - only accessible when not logged in */}
                <Route path="/login" element={
                    <PublicRoute>
                        <Login />
                    </PublicRoute>
                } />
                <Route path="/register" element={
                    <PublicRoute>
                        <Register />
                    </PublicRoute>
                } />

                {/* Protected routes - only accessible when logged in */}
                <Route path="/dashboard" element={
                    <PrivateRoute>
                        <Dashboard />
                    </PrivateRoute>
                } />
                <Route path="/analytics/:shortCode" element={
                    <PrivateRoute>
                        <Analytics />
                    </PrivateRoute>
                } />
                <Route path="/profile" element={
                    <PrivateRoute>
                        <Profile />
                    </PrivateRoute>
                } />

                {/* Error pages - public routes for backend redirects */}
                <Route path="/error/expired" element={<ErrorExpired />} />
                <Route path="/error/inactive" element={<ErrorInactive />} />
                <Route path="/error/not-found" element={<ErrorNotFound />} />
                <Route path="/error/server-error" element={<ErrorServer />} />

                {/* 404 page for all unmatched routes */}
                <Route path="*" element={<NotFound />} />
            </Routes>
        </Layout>
    )
}

function App() {
    return (
        <AuthProvider>
            <AppRoutes />
        </AuthProvider>
    )
}

export default App 
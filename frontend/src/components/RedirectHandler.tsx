import React, { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import NotFound from '../pages/NotFound'
import axios from 'axios'

interface RedirectHandlerProps {
    baseUrl?: string
}

export default function RedirectHandler({ baseUrl }: RedirectHandlerProps) {
    const { shortCode } = useParams<{ shortCode: string }>()
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState<{
        type: 'not-found' | 'inactive' | 'expired'
        message: string
    } | null>(null)

    useEffect(() => {
        if (shortCode) {
            handleRedirect()
        }
    }, [shortCode])

    const handleRedirect = async () => {
        if (!shortCode) {
            setError({ type: 'not-found', message: 'Invalid URL' })
            setLoading(false)
            return
        }

        try {
            // Make a request to the backend to get the URL
            const apiUrl = baseUrl || (import.meta as any).env.VITE_API_URL || ''
            const response = await axios.get(`${apiUrl}/api/v1/urls/${shortCode}`)
            
            if (response.data && response.data.original_url) {
                // Record the click and redirect
                await axios.post(`${apiUrl}/api/v1/urls/${shortCode}/click`, {
                    user_agent: navigator.userAgent,
                    referer: document.referrer
                })
                
                // Redirect to the original URL
                window.location.href = response.data.original_url
            } else {
                setError({ type: 'not-found', message: 'URL not found' })
            }
        } catch (err: any) {
            console.error('Redirect error:', err)
            
            if (err.response?.status === 404) {
                const errorMessage = err.response?.data?.error || 'URL not found'
                
                if (errorMessage.includes('not active')) {
                    setError({ type: 'inactive', message: errorMessage })
                } else if (errorMessage.includes('expired')) {
                    setError({ type: 'expired', message: errorMessage })
                } else {
                    setError({ type: 'not-found', message: errorMessage })
                }
            } else if (err.response?.status === 410) {
                setError({ type: 'inactive', message: 'This URL has been deactivated' })
            } else {
                setError({ type: 'not-found', message: 'URL not found' })
            }
        } finally {
            setLoading(false)
        }
    }

    if (loading) {
        return (
            <div className="min-h-screen bg-gray-50 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
                <div className="sm:mx-auto sm:w-full sm:max-w-md">
                    <div className="bg-white py-8 px-4 shadow sm:rounded-lg sm:px-10">
                        <div className="text-center">
                            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
                            <p className="mt-4 text-sm text-gray-600">
                                Redirecting...
                            </p>
                        </div>
                    </div>
                </div>
            </div>
        )
    }

    if (error) {
        return <NotFound type={error.type} message={error.message} />
    }

    return null
} 
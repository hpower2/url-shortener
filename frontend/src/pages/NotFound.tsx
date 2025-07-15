import React from 'react'
import { useLocation } from 'react-router-dom'
import ErrorLayout from '../components/ErrorLayout'

interface NotFoundProps {
    type?: 'not-found' | 'inactive' | 'expired'
    message?: string
}

export default function NotFound({ type = 'not-found', message }: NotFoundProps) {
    const location = useLocation()
    
    // For the generic 404 page, we don't have a short code
    // This is different from the error pages which are redirected from backend
    const shortCode = location.pathname !== '/' ? location.pathname.slice(1) : undefined
    
    const getTitle = () => {
        switch (type) {
            case 'inactive':
                return 'Link Deactivated'
            case 'expired':
                return 'Link Expired'
            default:
                return 'Page Not Found'
        }
    }

    const getDescription = () => {
        switch (type) {
            case 'inactive':
                return 'This shortened URL has been deactivated by its owner and is no longer accessible.'
            case 'expired':
                return 'This shortened URL has expired and is no longer accessible.'
            default:
                return message || 'The page you are looking for does not exist or may have been moved.'
        }
    }

    const getStatusCode = () => {
        switch (type) {
            case 'inactive':
            case 'expired':
                return '410'
            default:
                return '404'
        }
    }

    const getErrorType = (): 'expired' | 'inactive' | 'not-found' | 'server-error' => {
        switch (type) {
            case 'expired':
                return 'expired'
            case 'inactive':
                return 'inactive'
            default:
                return 'not-found'
        }
    }

    return (
        <ErrorLayout
            type={getErrorType()}
            shortCode={shortCode}
            title={getTitle()}
            description={getDescription()}
            statusCode={getStatusCode()}
            additionalInfo={type === 'not-found' ? 'This could be a typo in the URL or the page may have been moved or deleted.' : undefined}
        />
    )
} 
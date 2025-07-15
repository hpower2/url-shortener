import React from 'react'
import { useSearchParams } from 'react-router-dom'
import ErrorLayout from '../components/ErrorLayout'

export default function ErrorNotFound() {
    const [searchParams] = useSearchParams()
    const shortCode = searchParams.get('code')
    
    return (
        <ErrorLayout
            type="not-found"
            shortCode={shortCode || undefined}
            title="Link Not Found"
            description="The shortened URL you're looking for doesn't exist or may have been removed."
            statusCode="404"
            additionalInfo="Please check the URL and try again. If you believe this link should exist, contact the person who shared it with you."
        />
    )
} 
import React from 'react'
import { useSearchParams } from 'react-router-dom'
import ErrorLayout from '../components/ErrorLayout'

export default function ErrorServer() {
    const [searchParams] = useSearchParams()
    const shortCode = searchParams.get('code')
    
    return (
        <ErrorLayout
            type="server-error"
            shortCode={shortCode || undefined}
            title="Server Error"
            description="We're experiencing technical difficulties. Please try again later."
            statusCode="500"
            additionalInfo="Our team has been notified of this issue and is working to resolve it. If the problem persists, please contact our support team."
        />
    )
} 
import React from 'react'
import { useSearchParams } from 'react-router-dom'
import ErrorLayout from '../components/ErrorLayout'

export default function ErrorExpired() {
    const [searchParams] = useSearchParams()
    const shortCode = searchParams.get('code')
    
    return (
        <ErrorLayout
            type="expired"
            shortCode={shortCode || undefined}
            title="Link Expired"
            description="This shortened URL has reached its expiration date and is no longer accessible."
            statusCode="410"
            additionalInfo="The owner of this link set an expiration date to ensure the link was only available for a limited time. This helps maintain security and prevents outdated links from being accessed."
        />
    )
} 
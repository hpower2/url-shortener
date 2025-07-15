import React from 'react'
import { useSearchParams } from 'react-router-dom'
import ErrorLayout from '../components/ErrorLayout'

export default function ErrorInactive() {
    const [searchParams] = useSearchParams()
    const shortCode = searchParams.get('code')
    
    return (
        <ErrorLayout
            type="inactive"
            shortCode={shortCode || undefined}
            title="Link Deactivated"
            description="This shortened URL has been deactivated by its owner and is no longer accessible."
            statusCode="410"
            additionalInfo="The owner of this link has manually deactivated it. This could be for security reasons, content changes, or the link is no longer needed."
        />
    )
} 
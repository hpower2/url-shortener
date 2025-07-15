import React from 'react'
import { Link } from 'react-router-dom'
import { 
    HomeIcon, 
    MagnifyingGlassIcon,
    ClockIcon,
    XCircleIcon,
    ExclamationTriangleIcon,
    ServerIcon
} from '@heroicons/react/24/outline'

interface ErrorLayoutProps {
    type: 'expired' | 'inactive' | 'not-found' | 'server-error'
    shortCode?: string
    title: string
    description: string
    statusCode: string
    additionalInfo?: string
}

const iconMap = {
    expired: ClockIcon,
    inactive: XCircleIcon,
    'not-found': ExclamationTriangleIcon,
    'server-error': ServerIcon
}

const colorMap = {
    expired: {
        bg: 'from-orange-50 to-amber-50',
        border: 'border-orange-200',
        icon: 'text-orange-500',
        accent: 'text-orange-600',
        button: 'bg-orange-500 hover:bg-orange-600'
    },
    inactive: {
        bg: 'from-red-50 to-rose-50',
        border: 'border-red-200',
        icon: 'text-red-500',
        accent: 'text-red-600',
        button: 'bg-red-500 hover:bg-red-600'
    },
    'not-found': {
        bg: 'from-purple-50 to-indigo-50',
        border: 'border-purple-200',
        icon: 'text-purple-500',
        accent: 'text-purple-600',
        button: 'bg-purple-500 hover:bg-purple-600'
    },
    'server-error': {
        bg: 'from-gray-50 to-slate-50',
        border: 'border-gray-200',
        icon: 'text-gray-500',
        accent: 'text-gray-600',
        button: 'bg-gray-500 hover:bg-gray-600'
    }
}

export default function ErrorLayout({ 
    type, 
    shortCode, 
    title, 
    description, 
    statusCode, 
    additionalInfo 
}: ErrorLayoutProps) {
    const IconComponent = iconMap[type]
    const colors = colorMap[type]
    
    return (
        <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
            <div className="sm:mx-auto sm:w-full sm:max-w-md">
                <div className="bg-white py-8 px-4 shadow-xl sm:rounded-lg sm:px-10 relative overflow-hidden">
                    {/* Decorative background gradient */}
                    <div className={`absolute inset-0 bg-gradient-to-r ${colors.bg} opacity-50`}></div>
                    
                    {/* Content */}
                    <div className="relative z-10">
                        {/* Icon with animation */}
                        <div className="text-center">
                            <div className={`mx-auto flex items-center justify-center h-20 w-20 rounded-full ${colors.bg} ${colors.border} border-2 mb-6`}>
                                <IconComponent className={`h-10 w-10 ${colors.icon} animate-pulse`} />
                            </div>
                            
                            {/* Status Code */}
                            <h1 className={`text-6xl font-bold ${colors.accent} mb-2`}>
                                {statusCode}
                            </h1>
                            
                            {/* Title */}
                            <h2 className="text-2xl font-bold text-gray-900 mb-4">
                                {title}
                            </h2>
                            
                            {/* Description */}
                            <p className="text-gray-600 mb-6 leading-relaxed">
                                {description}
                            </p>
                            
                            {/* Short Code Info */}
                            {shortCode && (
                                <div className={`${colors.bg} ${colors.border} border rounded-lg p-4 mb-6`}>
                                    <p className="text-sm font-medium text-gray-700 mb-1">
                                        Attempted Short Code:
                                    </p>
                                    <code className={`text-sm font-mono ${colors.accent} font-semibold`}>
                                        {window.location.host}/{shortCode}
                                    </code>
                                </div>
                            )}
                            
                            {/* Additional Info */}
                            {additionalInfo && (
                                <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-6">
                                    <p className="text-sm text-blue-700">
                                        {additionalInfo}
                                    </p>
                                </div>
                            )}
                            
                            {/* Action Buttons */}
                            <div className="space-y-3">
                                <Link
                                    to="/"
                                    className={`w-full flex justify-center items-center px-4 py-3 border border-transparent rounded-md shadow-sm text-sm font-medium text-white ${colors.button} focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-all duration-200 hover:shadow-lg`}
                                >
                                    <HomeIcon className="h-4 w-4 mr-2" />
                                    Create New Short Link
                                </Link>
                                
                                <Link
                                    to="/dashboard"
                                    className="w-full flex justify-center items-center px-4 py-3 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-all duration-200"
                                >
                                    <MagnifyingGlassIcon className="h-4 w-4 mr-2" />
                                    View My Links
                                </Link>
                            </div>
                            
                            {/* Help Text */}
                            <div className="mt-8 text-center">
                                <p className="text-xs text-gray-500">
                                    Need help? Contact our support team for assistance.
                                </p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    )
} 
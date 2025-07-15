import React, { useState, useEffect } from 'react'
import { useParams, Link } from 'react-router-dom'
import { urlsAPI, URLAnalytics } from '../services/api'
import {
    ChartBarIcon,
    ArrowLeftIcon,
    GlobeAltIcon,
    EyeIcon,
    CursorArrowRaysIcon
} from '@heroicons/react/24/outline'
import toast from 'react-hot-toast'

export default function Analytics() {
    const { shortCode } = useParams<{ shortCode: string }>()
    const [analytics, setAnalytics] = useState<URLAnalytics | null>(null)
    const [loading, setLoading] = useState(true)
    const [timeRange, setTimeRange] = useState(30)

    useEffect(() => {
        if (shortCode) {
            fetchAnalytics()
        }
    }, [shortCode, timeRange])

    const fetchAnalytics = async () => {
        if (!shortCode) return

        try {
            const response = await urlsAPI.getAnalytics(shortCode, timeRange)
            setAnalytics(response.data)
        } catch (error: any) {
            toast.error(error.response?.data?.error || 'Failed to fetch analytics')
        } finally {
            setLoading(false)
        }
    }

    if (loading) {
        return (
            <div className="min-h-screen flex items-center justify-center">
                <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-600"></div>
            </div>
        )
    }

    if (!analytics) {
        return (
            <div className="min-h-screen flex items-center justify-center">
                <div className="text-center">
                    <h2 className="text-2xl font-bold text-gray-900 mb-4">Analytics not found</h2>
                    <Link
                        to="/dashboard"
                        className="text-blue-600 hover:text-blue-500"
                    >
                        ← Back to Dashboard
                    </Link>
                </div>
            </div>
        )
    }

    return (
        <div className="min-h-screen bg-gray-50">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
                {/* Header */}
                <div className="mb-8">
                    <div className="flex items-center space-x-4 mb-4">
                        <Link
                            to="/dashboard"
                            className="flex items-center text-gray-500 hover:text-gray-700"
                        >
                            <ArrowLeftIcon className="h-5 w-5 mr-2" />
                            Back to Dashboard
                        </Link>
                    </div>
                    <div className="flex items-center justify-between">
                        <div>
                            <h1 className="text-3xl font-bold text-gray-900">
                                Analytics for /{shortCode}
                            </h1>
                            <p className="text-gray-600 mt-2">
                                Track your link performance and audience insights
                            </p>
                        </div>
                        <div>
                            <select
                                value={timeRange}
                                onChange={(e) => setTimeRange(Number(e.target.value))}
                                className="px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                            >
                                <option value={7}>Last 7 days</option>
                                <option value={30}>Last 30 days</option>
                                <option value={90}>Last 90 days</option>
                                <option value={365}>Last year</option>
                            </select>
                        </div>
                    </div>
                </div>

                {/* Stats Cards */}
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
                    <div className="bg-white rounded-lg shadow p-6">
                        <div className="flex items-center">
                            <div className="flex-shrink-0">
                                <EyeIcon className="h-8 w-8 text-blue-600" />
                            </div>
                            <div className="ml-4">
                                <p className="text-sm font-medium text-gray-500">Total Clicks</p>
                                <p className="text-2xl font-bold text-gray-900">{analytics.total_clicks}</p>
                            </div>
                        </div>
                    </div>

                    <div className="bg-white rounded-lg shadow p-6">
                        <div className="flex items-center">
                            <div className="flex-shrink-0">
                                <CursorArrowRaysIcon className="h-8 w-8 text-green-600" />
                            </div>
                            <div className="ml-4">
                                <p className="text-sm font-medium text-gray-500">Unique Clicks</p>
                                <p className="text-2xl font-bold text-gray-900">{analytics.unique_clicks}</p>
                            </div>
                        </div>
                    </div>

                    <div className="bg-white rounded-lg shadow p-6">
                        <div className="flex items-center">
                            <div className="flex-shrink-0">
                                <ChartBarIcon className="h-8 w-8 text-purple-600" />
                            </div>
                            <div className="ml-4">
                                <p className="text-sm font-medium text-gray-500">Today</p>
                                <p className="text-2xl font-bold text-gray-900">{analytics.clicks_today}</p>
                            </div>
                        </div>
                    </div>

                    <div className="bg-white rounded-lg shadow p-6">
                        <div className="flex items-center">
                            <div className="flex-shrink-0">
                                <GlobeAltIcon className="h-8 w-8 text-orange-600" />
                            </div>
                            <div className="ml-4">
                                <p className="text-sm font-medium text-gray-500">This Week</p>
                                <p className="text-2xl font-bold text-gray-900">{analytics.clicks_this_week}</p>
                            </div>
                        </div>
                    </div>
                </div>

                {/* Analytics Tables */}
                <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
                    {/* Top Countries */}
                    <div className="bg-white rounded-lg shadow">
                        <div className="px-6 py-4 border-b border-gray-200">
                            <h2 className="text-lg font-medium text-gray-900">Top Countries</h2>
                        </div>
                        <div className="p-6">
                            {analytics.top_countries && analytics.top_countries.length > 0 ? (
                                <div className="space-y-4">
                                    {analytics.top_countries.map((country, index) => (
                                        <div key={index} className="flex items-center justify-between">
                                            <div className="flex items-center space-x-3">
                                                <div className="flex-shrink-0">
                                                    <div className="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center">
                                                        <span className="text-sm font-medium text-blue-600">
                                                            {index + 1}
                                                        </span>
                                                    </div>
                                                </div>
                                                <div>
                                                    <p className="text-sm font-medium text-gray-900">
                                                        {country.country || 'Unknown'}
                                                    </p>
                                                </div>
                                            </div>
                                            <div className="text-sm text-gray-500">
                                                {country.clicks} clicks
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            ) : (
                                <p className="text-gray-500 text-center py-8">No geographic data available</p>
                            )}
                        </div>
                    </div>

                    {/* Top Referrers */}
                    <div className="bg-white rounded-lg shadow">
                        <div className="px-6 py-4 border-b border-gray-200">
                            <h2 className="text-lg font-medium text-gray-900">Top Referrers</h2>
                        </div>
                        <div className="p-6">
                            {analytics.top_referrers && analytics.top_referrers.length > 0 ? (
                                <div className="space-y-4">
                                    {analytics.top_referrers.map((referrer, index) => (
                                        <div key={index} className="flex items-center justify-between">
                                            <div className="flex items-center space-x-3">
                                                <div className="flex-shrink-0">
                                                    <div className="w-8 h-8 bg-green-100 rounded-full flex items-center justify-center">
                                                        <span className="text-sm font-medium text-green-600">
                                                            {index + 1}
                                                        </span>
                                                    </div>
                                                </div>
                                                <div className="min-w-0 flex-1">
                                                    <p className="text-sm font-medium text-gray-900 truncate">
                                                        {referrer.referrer || 'Direct'}
                                                    </p>
                                                </div>
                                            </div>
                                            <div className="text-sm text-gray-500">
                                                {referrer.clicks} clicks
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            ) : (
                                <p className="text-gray-500 text-center py-8">No referrer data available</p>
                            )}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    )
} 
import { Link } from 'react-router-dom'
import { LinkIcon, ChartBarIcon, QrCodeIcon, ShieldCheckIcon } from '@heroicons/react/24/outline'
import React from 'react'

export default function Home() {
    return (
        <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
            {/* Hero Section */}
            <div className="relative overflow-hidden">
                <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
                    <div className="relative pb-16 pt-6 sm:pb-24">
                        <div className="mx-auto max-w-2xl text-center">
                            <h1 className="text-4xl font-bold tracking-tight text-gray-900 sm:text-6xl">
                                Shorten URLs with
                                <span className="text-blue-600"> Power</span>
                            </h1>
                            <p className="mt-6 text-lg leading-8 text-gray-600">
                                Create short, memorable links with advanced analytics, QR codes, and custom codes.
                                Perfect for businesses, marketers, and individuals who want more from their links.
                            </p>
                            <div className="mt-10 flex items-center justify-center gap-x-6">
                                <Link
                                    to="/register"
                                    className="rounded-md bg-blue-600 px-3.5 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600"
                                >
                                    Get Started Free
                                </Link>
                                <Link
                                    to="/login"
                                    className="text-sm font-semibold leading-6 text-gray-900"
                                >
                                    Sign In <span aria-hidden="true">→</span>
                                </Link>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            {/* Features Section */}
            <div className="py-24 bg-white">
                <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
                    <div className="mx-auto max-w-2xl text-center">
                        <h2 className="text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">
                            Everything you need to manage your links
                        </h2>
                        <p className="mt-6 text-lg leading-8 text-gray-600">
                            Our URL shortener comes packed with professional features to help you track, manage, and optimize your links.
                        </p>
                    </div>
                    <div className="mx-auto mt-16 max-w-2xl sm:mt-20 lg:mt-24 lg:max-w-none">
                        <dl className="grid max-w-xl grid-cols-1 gap-x-8 gap-y-16 lg:max-w-none lg:grid-cols-3">
                            <div className="flex flex-col">
                                <dt className="flex items-center gap-x-3 text-base font-semibold leading-7 text-gray-900">
                                    <LinkIcon className="h-5 w-5 flex-none text-blue-600" />
                                    Custom Short Codes
                                </dt>
                                <dd className="mt-4 flex flex-auto flex-col text-base leading-7 text-gray-600">
                                    <p className="flex-auto">
                                        Create memorable, branded short links with custom codes that reflect your brand or campaign.
                                    </p>
                                </dd>
                            </div>
                            <div className="flex flex-col">
                                <dt className="flex items-center gap-x-3 text-base font-semibold leading-7 text-gray-900">
                                    <ChartBarIcon className="h-5 w-5 flex-none text-blue-600" />
                                    Advanced Analytics
                                </dt>
                                <dd className="mt-4 flex flex-auto flex-col text-base leading-7 text-gray-600">
                                    <p className="flex-auto">
                                        Track clicks, geographic data, referrers, and user behavior with detailed analytics and reports.
                                    </p>
                                </dd>
                            </div>
                            <div className="flex flex-col">
                                <dt className="flex items-center gap-x-3 text-base font-semibold leading-7 text-gray-900">
                                    <QrCodeIcon className="h-5 w-5 flex-none text-blue-600" />
                                    QR Code Generation
                                </dt>
                                <dd className="mt-4 flex flex-auto flex-col text-base leading-7 text-gray-600">
                                    <p className="flex-auto">
                                        Generate QR codes for your short links instantly, perfect for print materials and offline campaigns.
                                    </p>
                                </dd>
                            </div>
                            <div className="flex flex-col">
                                <dt className="flex items-center gap-x-3 text-base font-semibold leading-7 text-gray-900">
                                    <ShieldCheckIcon className="h-5 w-5 flex-none text-blue-600" />
                                    Secure & Private
                                </dt>
                                <dd className="mt-4 flex flex-auto flex-col text-base leading-7 text-gray-600">
                                    <p className="flex-auto">
                                        Your data is protected with JWT authentication and user isolation. Only you can see your links.
                                    </p>
                                </dd>
                            </div>
                            <div className="flex flex-col">
                                <dt className="flex items-center gap-x-3 text-base font-semibold leading-7 text-gray-900">
                                    <LinkIcon className="h-5 w-5 flex-none text-blue-600" />
                                    Link Management
                                </dt>
                                <dd className="mt-4 flex flex-auto flex-col text-base leading-7 text-gray-600">
                                    <p className="flex-auto">
                                        Edit, disable, or delete your links anytime. Set expiration dates and manage your link portfolio.
                                    </p>
                                </dd>
                            </div>
                            <div className="flex flex-col">
                                <dt className="flex items-center gap-x-3 text-base font-semibold leading-7 text-gray-900">
                                    <ChartBarIcon className="h-5 w-5 flex-none text-blue-600" />
                                    Real-time Tracking
                                </dt>
                                <dd className="mt-4 flex flex-auto flex-col text-base leading-7 text-gray-600">
                                    <p className="flex-auto">
                                        Monitor your link performance in real-time with click tracking and geographic insights.
                                    </p>
                                </dd>
                            </div>
                        </dl>
                    </div>
                </div>
            </div>

            {/* CTA Section */}
            <div className="bg-blue-600">
                <div className="px-6 py-24 sm:px-6 sm:py-32 lg:px-8">
                    <div className="mx-auto max-w-2xl text-center">
                        <h2 className="text-3xl font-bold tracking-tight text-white sm:text-4xl">
                            Ready to supercharge your links?
                        </h2>
                        <p className="mx-auto mt-6 max-w-xl text-lg leading-8 text-blue-100">
                            Join thousands of users who trust our platform for their link management needs.
                            Start creating powerful short links today.
                        </p>
                        <div className="mt-10 flex items-center justify-center gap-x-6">
                            <Link
                                to="/register"
                                className="rounded-md bg-white px-3.5 py-2.5 text-sm font-semibold text-blue-600 shadow-sm hover:bg-blue-50 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white"
                            >
                                Create Free Account
                            </Link>
                            <Link
                                to="/login"
                                className="text-sm font-semibold leading-6 text-white"
                            >
                                Already have an account? Sign in <span aria-hidden="true">→</span>
                            </Link>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    )
} 
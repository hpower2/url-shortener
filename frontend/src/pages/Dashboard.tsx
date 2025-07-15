import React, { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import { urlsAPI, URL, CreateURLRequest, UpdateURLRequest } from '../services/api'
import {
    PlusIcon,
    LinkIcon,
    ChartBarIcon,
    QrCodeIcon,
    TrashIcon,
    ClipboardIcon,
    PencilIcon,
    XMarkIcon,
    ClockIcon
} from '@heroicons/react/24/outline'
import toast from 'react-hot-toast'

export default function Dashboard() {
    const { user } = useAuth()
    const [urls, setUrls] = useState<URL[]>([])
    const [loading, setLoading] = useState(true)
    const [creating, setCreating] = useState(false)
    const [showCreateForm, setShowCreateForm] = useState(false)
    const [formData, setFormData] = useState<CreateURLRequest>({
        url: '',
        custom_code: '',
        expires_at: ''
    })

    // Save form data to localStorage whenever it changes
    useEffect(() => {
        if (formData.url || formData.custom_code || formData.expires_at) {
            localStorage.setItem('url-shortener-form-data', JSON.stringify(formData))
        }
    }, [formData])

    // Edit URL state
    const [editingUrl, setEditingUrl] = useState<URL | null>(null)
    const [updating, setUpdating] = useState(false)
    const [editFormData, setEditFormData] = useState<UpdateURLRequest>({
        original_url: '',
        is_active: true,
        expires_at: ''
    })

    useEffect(() => {
        fetchUrls()

        // Restore form data from localStorage
        const savedFormData = localStorage.getItem('url-shortener-form-data')
        if (savedFormData) {
            try {
                const parsedData = JSON.parse(savedFormData)
                setFormData(parsedData)
            } catch (error) {
                // Clear invalid data
                localStorage.removeItem('url-shortener-form-data')
            }
        }
    }, [])

    const fetchUrls = async () => {
        try {
            const response = await urlsAPI.getAll({ limit: 50, offset: 0 })
            setUrls(response.data.urls || [])
        } catch (error: any) {
            // Don't show error for 401 - let the interceptor handle it
            if (error.response?.status !== 401) {
                toast.error('Failed to fetch URLs')
            }
        } finally {
            setLoading(false)
        }
    }

    const handleCreate = async (e: React.FormEvent) => {
        e.preventDefault()
        setCreating(true)

        try {
            // Prepare the data, excluding empty expires_at
            const requestData: CreateURLRequest = {
                url: formData.url,
                custom_code: formData.custom_code || undefined,
                expires_at: formData.expires_at ? new Date(formData.expires_at).toISOString() : undefined
            }

            const response = await urlsAPI.create(requestData)
            setUrls([response.data, ...urls])
            setFormData({ url: '', custom_code: '', expires_at: '' })
            setShowCreateForm(false)

            // Clear saved form data after successful creation
            localStorage.removeItem('url-shortener-form-data')

            toast.success('URL created successfully!')
        } catch (error: any) {
            toast.error(error.response?.data?.error || 'Failed to create URL')
        } finally {
            setCreating(false)
        }
    }

    const handleDelete = async (shortCode: string) => {
        if (!confirm('Are you sure you want to delete this URL?')) return

        try {
            await urlsAPI.delete(shortCode)
            setUrls(urls.filter(url => url.short_code !== shortCode))
            toast.success('URL deleted successfully!')
        } catch (error) {
            toast.error('Failed to delete URL')
        }
    }

    const handleEditClick = (url: URL) => {
        setEditingUrl(url)
        setEditFormData({
            original_url: url.original_url,
            is_active: url.is_active,
            expires_at: url.expires_at ? new Date(url.expires_at).toISOString().slice(0, 16) : ''
        })
    }

    const handleUpdate = async (e: React.FormEvent) => {
        e.preventDefault()
        if (!editingUrl) return

        setUpdating(true)
        try {
            const requestData: UpdateURLRequest = {
                original_url: editFormData.original_url,
                is_active: editFormData.is_active,
                expires_at: editFormData.expires_at ? new Date(editFormData.expires_at).toISOString() : undefined
            }

            const response = await urlsAPI.update(editingUrl.short_code, requestData)
            setUrls(urls.map(url =>
                url.short_code === editingUrl.short_code ? response.data : url
            ))
            setEditingUrl(null)
            setEditFormData({ original_url: '', is_active: true, expires_at: '' })
            toast.success('URL updated successfully!')
        } catch (error: any) {
            toast.error(error.response?.data?.error || 'Failed to update URL')
        } finally {
            setUpdating(false)
        }
    }

    const handleCancelEdit = () => {
        setEditingUrl(null)
        setEditFormData({ original_url: '', is_active: true, expires_at: '' })
    }

    const copyToClipboard = (text: string) => {
        navigator.clipboard.writeText(text)
        toast.success('Copied to clipboard!')
    }

    const downloadQR = async (shortCode: string) => {
        try {
            const response = await urlsAPI.getQRCode(shortCode)
            const blob = new Blob([response.data], { type: 'image/png' })
            const url = window.URL.createObjectURL(blob)
            const a = document.createElement('a')
            a.href = url
            a.download = `qr-${shortCode}.png`
            document.body.appendChild(a)
            a.click()
            window.URL.revokeObjectURL(url)
            document.body.removeChild(a)
            toast.success('QR code downloaded!')
        } catch (error) {
            toast.error('Failed to download QR code')
        }
    }

    const getShortUrl = (shortCode: string) => {
        return `${(import.meta as any).env.VITE_API_URL}/${shortCode}`
    }

    // Helper function to check if URL is expired
    const isUrlExpired = (url: URL) => {
        if (!url.expires_at) return false
        return new Date() > new Date(url.expires_at)
    }

    // Helper function to check if URL is expiring soon (within 7 days)
    const isUrlExpiringSoon = (url: URL) => {
        if (!url.expires_at || isUrlExpired(url)) return false
        const daysUntilExpiry = Math.floor((new Date(url.expires_at).getTime() - new Date().getTime()) / (1000 * 60 * 60 * 24))
        return daysUntilExpiry <= 7 && daysUntilExpiry > 0
    }

    // Helper function to get URL status for styling
    const getUrlStatus = (url: URL) => {
        if (isUrlExpired(url)) return 'expired'
        if (!url.is_active) return 'inactive'
        if (isUrlExpiringSoon(url)) return 'expiring-soon'
        return 'active'
    }

    if (loading) {
        return (
            <div className="min-h-screen flex items-center justify-center">
                <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-600"></div>
            </div>
        )
    }

    return (
        <div className="min-h-screen bg-gray-50">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
                {/* Header */}
                <div className="mb-8">
                    <h1 className="text-3xl font-bold text-gray-900">
                        Welcome back, {user?.first_name}!
                    </h1>
                    <p className="text-gray-600 mt-2">
                        Manage your shortened URLs and track their performance
                    </p>
                </div>

                {/* Create URL Section */}
                <div className="bg-white rounded-lg shadow mb-8">
                    <div className="p-6">
                        {!showCreateForm ? (
                            <button
                                onClick={() => setShowCreateForm(true)}
                                className="w-full flex items-center justify-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700"
                            >
                                <PlusIcon className="h-5 w-5 mr-2" />
                                Create New Short URL
                            </button>
                        ) : (
                            <form onSubmit={handleCreate} className="space-y-4">
                                <div>
                                    <label htmlFor="url" className="block text-sm font-medium text-gray-700">
                                        Original URL *
                                    </label>
                                    <input
                                        type="url"
                                        id="url"
                                        required
                                        className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                                        placeholder="https://example.com/very/long/url"
                                        value={formData.url}
                                        onChange={(e) => setFormData({ ...formData, url: e.target.value })}
                                    />
                                </div>

                                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                    <div>
                                        <label htmlFor="custom_code" className="block text-sm font-medium text-gray-700">
                                            Custom Code (optional)
                                        </label>
                                        <input
                                            type="text"
                                            id="custom_code"
                                            className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                                            placeholder="my-custom-link"
                                            value={formData.custom_code}
                                            onChange={(e) => setFormData({ ...formData, custom_code: e.target.value })}
                                        />
                                    </div>

                                    <div>
                                        <label htmlFor="expires_at" className="block text-sm font-medium text-gray-700">
                                            Expires At (optional)
                                        </label>
                                        <input
                                            type="datetime-local"
                                            id="expires_at"
                                            className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                                            value={formData.expires_at}
                                            onChange={(e) => setFormData({ ...formData, expires_at: e.target.value })}
                                        />
                                    </div>
                                </div>

                                <div className="flex space-x-3">
                                    <button
                                        type="submit"
                                        disabled={creating}
                                        className="flex-1 bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 disabled:opacity-50"
                                    >
                                        {creating ? 'Creating...' : 'Create URL'}
                                    </button>
                                    <button
                                        type="button"
                                        onClick={() => setShowCreateForm(false)}
                                        className="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50"
                                    >
                                        Cancel
                                    </button>
                                </div>
                            </form>
                        )}
                    </div>
                </div>

                {/* URLs List */}
                <div className="bg-white rounded-lg shadow">
                    <div className="px-6 py-4 border-b border-gray-200">
                        <h2 className="text-lg font-medium text-gray-900">Your URLs</h2>
                    </div>

                    {urls.length === 0 ? (
                        <div className="p-6 text-center text-gray-500">
                            <LinkIcon className="h-12 w-12 mx-auto mb-4 text-gray-400" />
                            <p>No URLs created yet. Create your first short URL above!</p>
                        </div>
                    ) : (
                        <div className="divide-y divide-gray-200">
                            {urls.map((url) => (
                                <div key={url.id} className="p-6">
                                    <div className="flex items-center justify-between">
                                        <div className="flex-1 min-w-0">
                                            <div className="flex items-center space-x-3">
                                                <div className="flex-shrink-0">
                                                    <LinkIcon className={`h-5 w-5 ${
                                                        getUrlStatus(url) === 'expired' ? 'text-orange-400' : 
                                                        getUrlStatus(url) === 'inactive' ? 'text-red-400' : 
                                                        getUrlStatus(url) === 'expiring-soon' ? 'text-yellow-400' : 
                                                        'text-gray-400'
                                                    }`} />
                                                </div>
                                                <div className="flex-1 min-w-0">
                                                    <div className="flex items-center space-x-2">
                                                        <p className={`text-sm font-medium truncate ${
                                                            getUrlStatus(url) === 'active' ? 'text-gray-900' : 'text-gray-400'
                                                        }`}>
                                                            {getShortUrl(url.short_code)}
                                                        </p>
                                                        {getUrlStatus(url) === 'expired' && (
                                                            <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-orange-100 text-orange-800">
                                                                <ClockIcon className="h-3 w-3 mr-1" />
                                                                Expired
                                                            </span>
                                                        )}
                                                        {getUrlStatus(url) === 'inactive' && (
                                                            <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-red-100 text-red-800">
                                                                Inactive
                                                            </span>
                                                        )}
                                                        {getUrlStatus(url) === 'expiring-soon' && (
                                                            <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-yellow-100 text-yellow-800">
                                                                <ClockIcon className="h-3 w-3 mr-1" />
                                                                Expiring Soon
                                                            </span>
                                                        )}
                                                    </div>
                                                    <p className="text-sm text-gray-500 truncate">
                                                        {url.original_url}
                                                    </p>
                                                </div>
                                            </div>
                                            <div className="mt-2 flex items-center text-sm text-gray-500 space-x-4">
                                                <span>{url.click_count} clicks</span>
                                                <span>Created {new Date(url.created_at).toLocaleDateString()}</span>
                                                {url.expires_at && (
                                                    <span className={
                                                        isUrlExpired(url) ? 'text-orange-600 font-medium' : 
                                                        isUrlExpiringSoon(url) ? 'text-yellow-600 font-medium' : ''
                                                    }>
                                                        {isUrlExpired(url) ? 'Expired' : 'Expires'} {new Date(url.expires_at).toLocaleDateString()}
                                                        {isUrlExpired(url) && (
                                                            <span className="ml-1 text-xs">
                                                                ({Math.floor((new Date().getTime() - new Date(url.expires_at).getTime()) / (1000 * 60 * 60 * 24))} days ago)
                                                            </span>
                                                        )}
                                                        {isUrlExpiringSoon(url) && (
                                                            <span className="ml-1 text-xs">
                                                                ({Math.floor((new Date(url.expires_at).getTime() - new Date().getTime()) / (1000 * 60 * 60 * 24))} days left)
                                                            </span>
                                                        )}
                                                    </span>
                                                )}
                                            </div>
                                        </div>

                                        <div className="flex items-center space-x-2">
                                            <button
                                                onClick={() => copyToClipboard(getShortUrl(url.short_code))}
                                                className="p-2 text-gray-400 hover:text-gray-600"
                                                title="Copy URL"
                                            >
                                                <ClipboardIcon className="h-5 w-5" />
                                            </button>
                                            <button
                                                onClick={() => downloadQR(url.short_code)}
                                                className="p-2 text-gray-400 hover:text-gray-600"
                                                title="Download QR Code"
                                            >
                                                <QrCodeIcon className="h-5 w-5" />
                                            </button>
                                            <Link
                                                to={`/analytics/${url.short_code}`}
                                                className="p-2 text-gray-400 hover:text-gray-600"
                                                title="View Analytics"
                                            >
                                                <ChartBarIcon className="h-5 w-5" />
                                            </Link>
                                            <button
                                                onClick={() => handleEditClick(url)}
                                                className="p-2 text-blue-400 hover:text-blue-600"
                                                title="Edit URL"
                                            >
                                                <PencilIcon className="h-5 w-5" />
                                            </button>
                                            <button
                                                onClick={() => handleDelete(url.short_code)}
                                                className="p-2 text-red-400 hover:text-red-600"
                                                title="Delete URL"
                                            >
                                                <TrashIcon className="h-5 w-5" />
                                            </button>
                                        </div>
                                    </div>
                                </div>
                            ))}
                        </div>
                    )}
                </div>

                {/* Edit URL Modal */}
                {editingUrl && (
                    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
                        <div className="bg-white rounded-lg p-6 w-full max-w-md">
                            <div className="flex justify-between items-center mb-4">
                                <h3 className="text-lg font-medium text-gray-900">Edit URL</h3>
                                <button
                                    onClick={handleCancelEdit}
                                    className="text-gray-400 hover:text-gray-600"
                                >
                                    <XMarkIcon className="h-6 w-6" />
                                </button>
                            </div>

                            <form onSubmit={handleUpdate} className="space-y-4">
                                <div>
                                    <label htmlFor="edit-url" className="block text-sm font-medium text-gray-700">
                                        Original URL *
                                    </label>
                                    <input
                                        type="url"
                                        id="edit-url"
                                        required
                                        className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                                        value={editFormData.original_url}
                                        onChange={(e) => setEditFormData({ ...editFormData, original_url: e.target.value })}
                                    />
                                </div>

                                <div>
                                    <label htmlFor="edit-expires" className="block text-sm font-medium text-gray-700">
                                        Expires At (optional)
                                    </label>
                                    <input
                                        type="datetime-local"
                                        id="edit-expires"
                                        className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                                        value={editFormData.expires_at}
                                        onChange={(e) => setEditFormData({ ...editFormData, expires_at: e.target.value })}
                                    />
                                </div>

                                <div className="flex items-center">
                                    <input
                                        type="checkbox"
                                        id="edit-active"
                                        checked={editFormData.is_active}
                                        onChange={(e) => setEditFormData({ ...editFormData, is_active: e.target.checked })}
                                        className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                                    />
                                    <label htmlFor="edit-active" className="ml-2 block text-sm text-gray-900">
                                        Active
                                    </label>
                                </div>

                                <div className="flex justify-end space-x-3 pt-4">
                                    <button
                                        type="button"
                                        onClick={handleCancelEdit}
                                        className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
                                    >
                                        Cancel
                                    </button>
                                    <button
                                        type="submit"
                                        disabled={updating}
                                        className="px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md hover:bg-blue-700 disabled:opacity-50"
                                    >
                                        {updating ? 'Updating...' : 'Update URL'}
                                    </button>
                                </div>
                            </form>
                        </div>
                    </div>
                )}
            </div>
        </div>
    )
} 
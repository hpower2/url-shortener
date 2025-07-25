import { Fragment } from 'react'
import { Disclosure, Menu, Transition } from '@headlessui/react'
import { Bars3Icon, XMarkIcon, LinkIcon } from '@heroicons/react/24/outline'
import { Link, useLocation } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'

interface LayoutProps {
    children: React.ReactNode
}

function classNames(...classes: string[]) {
    return classes.filter(Boolean).join(' ')
}

export default function Layout({ children }: LayoutProps) {
    const { user, logout } = useAuth()
    const location = useLocation()

    const navigation = user
        ? [
            { name: 'Dashboard', href: '/dashboard', current: location.pathname === '/dashboard' },
            { name: 'Profile', href: '/profile', current: location.pathname === '/profile' },
        ]
        : [
            { name: 'Home', href: '/', current: location.pathname === '/' },
            { name: 'Login', href: '/login', current: location.pathname === '/login' },
            { name: 'Register', href: '/register', current: location.pathname === '/register' },
        ]

    return (
        <div className="min-h-screen bg-gray-50">
            <Disclosure as="nav" className="bg-white shadow">
                {({ open }) => (
                    <>
                        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
                            <div className="flex h-16 justify-between">
                                <div className="flex">
                                    <div className="flex flex-shrink-0 items-center">
                                        <Link to={user ? '/dashboard' : '/'} className="flex items-center">
                                            <LinkIcon className="h-8 w-8 text-primary-600" />
                                            <span className="ml-2 text-xl font-bold text-gray-900">URL Shortener</span>
                                        </Link>
                                    </div>
                                    <div className="hidden sm:ml-6 sm:flex sm:space-x-8">
                                        {navigation.map((item) => (
                                            <Link
                                                key={item.name}
                                                to={item.href}
                                                className={classNames(
                                                    item.current
                                                        ? 'border-primary-500 text-gray-900'
                                                        : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700',
                                                    'inline-flex items-center border-b-2 px-1 pt-1 text-sm font-medium'
                                                )}
                                            >
                                                {item.name}
                                            </Link>
                                        ))}
                                    </div>
                                </div>

                                {user && (
                                    <div className="hidden sm:ml-6 sm:flex sm:items-center">
                                        <Menu as="div" className="relative ml-3">
                                            <div>
                                                <Menu.Button className="flex rounded-full bg-white text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2">
                                                    <span className="sr-only">Open user menu</span>
                                                    <div className="h-8 w-8 rounded-full bg-primary-600 flex items-center justify-center">
                                                        <span className="text-sm font-medium text-white">
                                                            {user.first_name[0]}{user.last_name[0]}
                                                        </span>
                                                    </div>
                                                </Menu.Button>
                                            </div>
                                            <Transition
                                                as={Fragment}
                                                enter="transition ease-out duration-200"
                                                enterFrom="transform opacity-0 scale-95"
                                                enterTo="transform opacity-100 scale-100"
                                                leave="transition ease-in duration-75"
                                                leaveFrom="transform opacity-100 scale-100"
                                                leaveTo="transform opacity-0 scale-95"
                                            >
                                                <Menu.Items className="absolute right-0 z-10 mt-2 w-48 origin-top-right rounded-md bg-white py-1 shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
                                                    <Menu.Item>
                                                        {({ active }) => (
                                                            <Link
                                                                to="/profile"
                                                                className={classNames(
                                                                    active ? 'bg-gray-100' : '',
                                                                    'block px-4 py-2 text-sm text-gray-700'
                                                                )}
                                                            >
                                                                Your Profile
                                                            </Link>
                                                        )}
                                                    </Menu.Item>
                                                    <Menu.Item>
                                                        {({ active }) => (
                                                            <button
                                                                onClick={logout}
                                                                className={classNames(
                                                                    active ? 'bg-gray-100' : '',
                                                                    'block w-full text-left px-4 py-2 text-sm text-gray-700'
                                                                )}
                                                            >
                                                                Sign out
                                                            </button>
                                                        )}
                                                    </Menu.Item>
                                                </Menu.Items>
                                            </Transition>
                                        </Menu>
                                    </div>
                                )}

                                <div className="-mr-2 flex items-center sm:hidden">
                                    <Disclosure.Button className="inline-flex items-center justify-center rounded-md bg-white p-2 text-gray-400 hover:bg-gray-100 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-inset">
                                        <span className="sr-only">Open main menu</span>
                                        {open ? (
                                            <XMarkIcon className="block h-6 w-6" aria-hidden="true" />
                                        ) : (
                                            <Bars3Icon className="block h-6 w-6" aria-hidden="true" />
                                        )}
                                    </Disclosure.Button>
                                </div>
                            </div>
                        </div>

                        <Disclosure.Panel className="sm:hidden">
                            <div className="space-y-1 pb-3 pt-2">
                                {navigation.map((item) => (
                                    <Disclosure.Button
                                        key={item.name}
                                        as={Link}
                                        to={item.href}
                                        className={classNames(
                                            item.current
                                                ? 'border-primary-500 bg-primary-50 text-primary-700'
                                                : 'border-transparent text-gray-600 hover:border-gray-300 hover:bg-gray-50 hover:text-gray-800',
                                            'block border-l-4 py-2 pl-3 pr-4 text-base font-medium'
                                        )}
                                    >
                                        {item.name}
                                    </Disclosure.Button>
                                ))}
                                {user && (
                                    <>
                                        <div className="border-t border-gray-200 pt-4">
                                            <div className="flex items-center px-4">
                                                <div className="flex-shrink-0">
                                                    <div className="h-10 w-10 rounded-full bg-primary-600 flex items-center justify-center">
                                                        <span className="text-sm font-medium text-white">
                                                            {user.first_name[0]}{user.last_name[0]}
                                                        </span>
                                                    </div>
                                                </div>
                                                <div className="ml-3">
                                                    <div className="text-base font-medium text-gray-800">
                                                        {user.first_name} {user.last_name}
                                                    </div>
                                                    <div className="text-sm font-medium text-gray-500">{user.email}</div>
                                                </div>
                                            </div>
                                            <div className="mt-3 space-y-1">
                                                <Disclosure.Button
                                                    as={Link}
                                                    to="/profile"
                                                    className="block px-4 py-2 text-base font-medium text-gray-500 hover:bg-gray-100 hover:text-gray-800"
                                                >
                                                    Your Profile
                                                </Disclosure.Button>
                                                <Disclosure.Button
                                                    as="button"
                                                    onClick={logout}
                                                    className="block w-full text-left px-4 py-2 text-base font-medium text-gray-500 hover:bg-gray-100 hover:text-gray-800"
                                                >
                                                    Sign out
                                                </Disclosure.Button>
                                            </div>
                                        </div>
                                    </>
                                )}
                            </div>
                        </Disclosure.Panel>
                    </>
                )}
            </Disclosure>

            <main className="mx-auto max-w-7xl py-6 sm:px-6 lg:px-8">
                {children}
            </main>
        </div>
    )
} 
import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';
import { useTheme } from '../../contexts/ThemeContext';
import {
    UserIcon,
    Cog6ToothIcon,
    ArrowRightOnRectangleIcon,
    Bars3Icon,
    XMarkIcon,
    SunIcon,
    MoonIcon
} from '@heroicons/react/24/outline';
import { Menu, Transition } from '@headlessui/react';

const Navbar: React.FC = () => {
    const { user, logout } = useAuth();
    const { isDarkMode, toggleTheme } = useTheme();
    const navigate = useNavigate();
    const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

    const handleLogout = () => {
        logout();
        navigate('/login');
    };

    const getUserInitials = (firstName: string, lastName: string) => {
        return `${firstName.charAt(0)}${lastName.charAt(0)}`.toUpperCase();
    };

    return (
        <nav className="bg-white border-b border-gray-200 fixed w-full top-0 z-50">
            <div className="w-full px-4 lg:px-6">
                <div className="flex justify-between items-center h-16">
                    {/* Logo and brand */}
                    <div className="flex items-center">
                        <Link to="/" className="flex items-center group">
                            <div className="flex-shrink-0">
                                <div className="flex items-center space-x-2">
                                    <div className="w-8 h-8 bg-gradient-to-br from-purple-500 via-pink-500 to-orange-400 rounded-lg flex items-center justify-center shadow-lg">
                                        <span className="text-white font-bold text-sm">W</span>
                                    </div>
                                    <h1 className="text-2xl font-extrabold text-gray-900 group-hover:text-gray-700 transition-colors duration-200" style={{ fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif', letterSpacing: '-0.02em' }}>
                                        WanderSphere
                                    </h1>
                                </div>
                            </div>
                        </Link>
                    </div>

                    {/* Desktop menu */}
                    <div className="hidden md:flex md:items-center md:space-x-6">
                        {/* Search bar - Instagram style */}
                        <div className="relative hidden lg:block">
                            <input
                                type="text"
                                placeholder="Search"
                                className="w-64 px-4 py-2 bg-gray-50 border border-gray-200 rounded-lg text-sm placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                            />
                        </div>

                        {/* Theme toggle */}
                        <button
                            onClick={toggleTheme}
                            className="p-2 rounded-full hover:bg-gray-100 text-gray-600 hover:text-gray-800 transition-all duration-200"
                            aria-label="Toggle theme"
                        >
                            {isDarkMode ? (
                                <SunIcon className="h-6 w-6" />
                            ) : (
                                <MoonIcon className="h-6 w-6" />
                            )}
                        </button>

                        {/* User menu */}
                        <Menu as="div" className="relative">
                            <div>
                                <Menu.Button className="flex items-center text-sm rounded-full hover:bg-gray-50 p-1 transition-all duration-200">
                                    <span className="sr-only">Open user menu</span>
                                    <div className="story-border w-8 h-8">
                                        {user?.profile_picture ? (
                                            <img
                                                className="w-full h-full rounded-full object-cover"
                                                src={user.profile_picture}
                                                alt={`${user.first_name} ${user.last_name}`}
                                            />
                                        ) : (
                                            <div className="w-full h-full rounded-full bg-gray-300 flex items-center justify-center text-gray-600 font-medium text-xs">
                                                {user ? getUserInitials(user.first_name, user.last_name) : 'U'}
                                            </div>
                                        )}
                                    </div>
                                    <span className="ml-3 text-gray-900 font-medium hidden xl:block">
                                        {user ? user.user_name : 'User'}
                                    </span>
                                </Menu.Button>
                            </div>

                            <Transition
                                as={React.Fragment}
                                enter="transition ease-out duration-100"
                                enterFrom="transform opacity-0 scale-95"
                                enterTo="transform opacity-100 scale-100"
                                leave="transition ease-in duration-75"
                                leaveFrom="transform opacity-100 scale-100"
                                leaveTo="transform opacity-0 scale-95"
                            >
                                <Menu.Items className="origin-top-right absolute right-0 mt-2 w-56 rounded-lg shadow-lg bg-white border border-gray-200 ring-1 ring-black ring-opacity-5 focus:outline-none">
                                    <div className="py-2">
                                        <Menu.Item>
                                            {({ active }) => (
                                                <Link
                                                    to={`/profile/${user?.user_id}`}
                                                    className={`${active ? 'bg-gray-50 text-gray-900' : 'text-gray-700'
                                                        } flex items-center px-4 py-3 text-sm transition-colors duration-200 mx-2 rounded-lg`}
                                                >
                                                    <UserIcon className="mr-3 h-5 w-5" />
                                                    My Profile
                                                </Link>
                                            )}
                                        </Menu.Item>
                                        <Menu.Item>
                                            {({ active }) => (
                                                <Link
                                                    to="/profile/edit"
                                                    className={`${active ? 'bg-gray-50 text-gray-900' : 'text-gray-700'
                                                        } flex items-center px-4 py-3 text-sm transition-colors duration-200 mx-2 rounded-lg`}
                                                >
                                                    <Cog6ToothIcon className="mr-3 h-5 w-5" />
                                                    Settings
                                                </Link>
                                            )}
                                        </Menu.Item>
                                        <hr className="my-2 border-gray-200" />
                                        <Menu.Item>
                                            {({ active }) => (
                                                <button
                                                    onClick={handleLogout}
                                                    className={`${active ? 'bg-red-50 text-red-600' : 'text-gray-700'
                                                        } flex w-full items-center px-4 py-3 text-sm transition-colors duration-200 mx-2 rounded-lg`}
                                                >
                                                    <ArrowRightOnRectangleIcon className="mr-3 h-5 w-5" />
                                                    Sign out
                                                </button>
                                            )}
                                        </Menu.Item>
                                    </div>
                                </Menu.Items>
                            </Transition>
                        </Menu>
                    </div>

                    {/* Mobile menu button */}
                    <div className="md:hidden flex items-center space-x-2">
                        {/* Mobile theme toggle */}
                        <button
                            onClick={toggleTheme}
                            className="p-2 rounded-full hover:bg-gray-100 text-gray-600 hover:text-gray-800 transition-all duration-200"
                            aria-label="Toggle theme"
                        >
                            {isDarkMode ? (
                                <SunIcon className="h-6 w-6" />
                            ) : (
                                <MoonIcon className="h-6 w-6" />
                            )}
                        </button>

                        <button
                            onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
                            className="p-2 rounded-full hover:bg-gray-100 text-gray-600 hover:text-gray-800 transition-all duration-200"
                        >
                            {isMobileMenuOpen ? (
                                <XMarkIcon className="h-6 w-6" />
                            ) : (
                                <Bars3Icon className="h-6 w-6" />
                            )}
                        </button>
                    </div>
                </div>
            </div>

            {/* Mobile menu */}
            {isMobileMenuOpen && (
                <div className="md:hidden animate-slideDown">
                    <div className="px-4 pt-2 pb-3 space-y-2 bg-white/95 backdrop-blur-sm border-t border-gray-200">
                        <Link
                            to={`/profile/${user?.user_id}`}
                            className="flex items-center px-3 py-2 text-base font-medium text-gray-700 hover:text-gray-900 hover:bg-gray-50 rounded-lg transition-colors duration-200"
                            onClick={() => setIsMobileMenuOpen(false)}
                        >
                            <UserIcon className="mr-3 h-5 w-5" />
                            My Profile
                        </Link>
                        <Link
                            to="/profile/edit"
                            className="flex items-center px-3 py-2 text-base font-medium text-gray-700 hover:text-gray-900 hover:bg-gray-50 rounded-lg transition-colors duration-200"
                            onClick={() => setIsMobileMenuOpen(false)}
                        >
                            <Cog6ToothIcon className="mr-3 h-5 w-5" />
                            Settings
                        </Link>
                        <hr className="my-2 border-gray-200" />
                        <button
                            onClick={() => {
                                handleLogout();
                                setIsMobileMenuOpen(false);
                            }}
                            className="flex w-full items-center px-3 py-2 text-base font-medium text-red-600 hover:bg-red-50 rounded-lg transition-colors duration-200"
                        >
                            <ArrowRightOnRectangleIcon className="mr-3 h-5 w-5" />
                            Sign out
                        </button>
                    </div>
                </div>
            )}
        </nav>
    );
};

export default Navbar; 
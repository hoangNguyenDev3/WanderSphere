import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import {
    HomeIcon,
    PlusCircleIcon,
    MagnifyingGlassIcon,
    UserGroupIcon,
} from '@heroicons/react/24/outline';
import {
    HomeIcon as HomeIconSolid,
    PlusCircleIcon as PlusCircleIconSolid,
    MagnifyingGlassIcon as MagnifyingGlassIconSolid,
    UserGroupIcon as UserGroupIconSolid,
} from '@heroicons/react/24/solid';
import { useAuth } from '../../contexts/AuthContext';

interface SidebarLinkProps {
    to: string;
    icon: React.ComponentType<{ className?: string }>;
    activeIcon: React.ComponentType<{ className?: string }>;
    children: React.ReactNode;
    exact?: boolean;
}

const SidebarLink: React.FC<SidebarLinkProps> = ({
    to,
    icon: Icon,
    activeIcon: ActiveIcon,
    children,
    exact = false
}) => {
    const location = useLocation();
    const isActive = exact ? location.pathname === to : location.pathname === to;

    return (
        <Link
            to={to}
            className={`
                group flex items-center px-4 py-3 text-sm font-medium rounded-xl transition-all duration-200 
                ${isActive
                    ? 'bg-white text-gray-900 shadow-md border border-gray-200'
                    : 'text-gray-600 hover:text-gray-900 hover:bg-gray-50'
                }
                relative overflow-hidden min-h-[48px]
            `}
        >
            <div className="flex items-center w-full">
                {isActive ? (
                    <ActiveIcon className="mr-3 h-6 w-6 flex-shrink-0" />
                ) : (
                    <Icon className="mr-3 h-6 w-6 flex-shrink-0 group-hover:scale-110 transition-transform duration-200" />
                )}
                <span className="truncate">{children}</span>
            </div>
            {isActive && (
                <div className="absolute inset-0 bg-gradient-to-r from-white/10 to-transparent pointer-events-none" />
            )}
        </Link>
    );
};

interface MobileLinkProps {
    to: string;
    icon: React.ComponentType<{ className?: string }>;
    activeIcon: React.ComponentType<{ className?: string }>;
    children: React.ReactNode;
    exact?: boolean;
    currentPath: string;
}

const MobileLink: React.FC<MobileLinkProps> = ({
    to,
    icon: Icon,
    activeIcon: ActiveIcon,
    children,
    exact = false,
    currentPath
}) => {
    const isActive = exact ? currentPath === to : currentPath === to;

    return (
        <Link
            to={to}
            className={`
                flex flex-col items-center justify-center py-2 px-1 rounded-lg transition-all duration-200 min-h-[60px]
                ${isActive
                    ? 'text-primary bg-primary/10'
                    : 'text-muted-foreground hover:text-foreground hover:bg-muted'
                }
            `}
        >
            {isActive ? (
                <ActiveIcon className="h-6 w-6 mb-1" />
            ) : (
                <Icon className="h-6 w-6 mb-1" />
            )}
            <span className="text-xs font-medium text-center leading-tight">
                {children}
            </span>
        </Link>
    );
};

const Sidebar: React.FC = () => {
    const { user } = useAuth();
    const location = useLocation();

    const navigation = [
        {
            name: 'Home',
            to: '/',
            icon: HomeIcon,
            activeIcon: HomeIconSolid,
            exact: true,
        },
        {
            name: 'Search',
            to: '/search',
            icon: MagnifyingGlassIcon,
            activeIcon: MagnifyingGlassIconSolid,
            exact: true,
        },
        {
            name: 'Create',
            to: '/create',
            icon: PlusCircleIcon,
            activeIcon: PlusCircleIconSolid,
            exact: true,
        },
        {
            name: 'Connections',
            to: `/profile/${user?.user_id}/following`,
            icon: UserGroupIcon,
            activeIcon: UserGroupIconSolid,
            exact: true,
        },
    ];

    return (
        <>
            {/* Desktop Sidebar */}
            <div className="hidden lg:block lg:fixed lg:left-0 lg:top-16 lg:w-64 lg:z-30" style={{ height: 'calc(100vh - 4rem)' }}>
                <div className="flex flex-col h-full bg-white border-r border-gray-200">
                    {/* Navigation */}
                    <nav className="flex-1 px-4 space-y-2 py-6 overflow-y-auto">
                        {navigation.map((item) => (
                            <SidebarLink
                                key={item.name}
                                to={item.to}
                                icon={item.icon}
                                activeIcon={item.activeIcon}
                                exact={item.exact}
                            >
                                {item.name}
                            </SidebarLink>
                        ))}
                    </nav>

                    {/* User info section - Fixed at bottom */}
                    {user && (
                        <div className="flex-shrink-0 p-4 border-t border-gray-100">
                            <Link to={`/profile/${user.user_id}`}>
                                <div className="bg-gray-50 rounded-xl p-4 hover:bg-gray-100 transition-all duration-200 cursor-pointer">
                                    <div className="flex items-center">
                                        <div className="story-border flex-shrink-0">
                                            {user.profile_picture ? (
                                                <img
                                                    className="h-10 w-10 rounded-full object-cover"
                                                    src={user.profile_picture}
                                                    alt={`${user.first_name} ${user.last_name}`}
                                                />
                                            ) : (
                                                <div className="h-10 w-10 rounded-full bg-gray-300 flex items-center justify-center text-gray-600 font-medium">
                                                    {user.first_name.charAt(0) + user.last_name.charAt(0)}
                                                </div>
                                            )}
                                        </div>
                                        <div className="ml-3 min-w-0 flex-1">
                                            <p className="text-username truncate">
                                                @{user.user_name}
                                            </p>
                                            <p className="text-muted-instagram truncate">
                                                {user.first_name} {user.last_name}
                                            </p>
                                        </div>
                                    </div>
                                </div>
                            </Link>
                        </div>
                    )}
                </div>
            </div>

            {/* Mobile Bottom Navigation */}
            <div className="lg:hidden fixed bottom-0 left-0 right-0 z-40">
                <div className="bg-white/95 backdrop-blur-md border-t border-gray-200">
                    <div className="grid grid-cols-4 px-2 py-1">
                        {navigation.map((item) => (
                            <MobileLink
                                key={item.name}
                                to={item.to}
                                icon={item.icon}
                                activeIcon={item.activeIcon}
                                exact={item.exact}
                                currentPath={location.pathname}
                            >
                                {item.name}
                            </MobileLink>
                        ))}
                    </div>
                </div>
            </div >
        </>
    );
};

export default Sidebar; 
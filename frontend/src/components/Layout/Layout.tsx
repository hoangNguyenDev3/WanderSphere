import React from 'react';
import { useLocation } from 'react-router-dom';
import Navbar from './Navbar';
import Sidebar from './Sidebar';
import RightSidebar from './RightSidebar';

interface LayoutProps {
    children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
    const location = useLocation();

    // Show right sidebar only on home page and feed-related pages
    const showRightSidebar = ['/', '/home'].includes(location.pathname);

    return (
        <div className="min-h-screen bg-gray-50 text-gray-900">
            <Navbar />

            <Sidebar />

            <div className="pt-16">
                {/* Main content area with conditional layout */}
                <main className={`lg:ml-64 ${showRightSidebar ? 'lg:mr-80' : ''}`}>
                    <div className={`mx-auto px-4 sm:px-6 lg:px-8 py-8 ${showRightSidebar
                        ? 'max-w-2xl' // Narrower for home page with sidebar
                        : 'max-w-4xl' // Wider for other pages
                        }`}>
                        {children}
                    </div>
                </main>

                {/* Right sidebar - only show on specific pages */}
                {showRightSidebar && (
                    <div className="hidden lg:block fixed right-0 top-16 h-screen overflow-y-auto z-30">
                        <RightSidebar />
                    </div>
                )}
            </div>
        </div>
    );
};

export default Layout; 
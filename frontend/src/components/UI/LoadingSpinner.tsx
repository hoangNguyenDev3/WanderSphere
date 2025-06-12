import React from 'react';

interface LoadingSpinnerProps {
    size?: 'sm' | 'md' | 'lg' | 'xl';
    className?: string;
    color?: 'primary' | 'secondary' | 'white' | 'current';
}

const LoadingSpinner: React.FC<LoadingSpinnerProps> = ({
    size = 'md',
    className = '',
    color = 'primary'
}) => {
    const sizes = {
        sm: 'h-4 w-4',
        md: 'h-6 w-6',
        lg: 'h-8 w-8',
        xl: 'h-12 w-12'
    };

    const colors = {
        primary: 'text-primary',
        secondary: 'text-secondary',
        white: 'text-white',
        current: 'text-current'
    };

    const strokeWidth = size === 'sm' ? '3' : size === 'md' ? '2.5' : '2';

    return (
        <div
            className={`inline-block animate-spin ${sizes[size]} ${colors[color]} ${className}`}
            role="status"
            aria-label="Loading"
        >
            <svg
                className="w-full h-full"
                viewBox="0 0 24 24"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
            >
                <circle
                    cx="12"
                    cy="12"
                    r="10"
                    stroke="currentColor"
                    strokeWidth={strokeWidth}
                    strokeLinecap="round"
                    strokeDasharray="31.416"
                    strokeDashoffset="31.416"
                    fill="none"
                    opacity="0.2"
                />
                <circle
                    cx="12"
                    cy="12"
                    r="10"
                    stroke="currentColor"
                    strokeWidth={strokeWidth}
                    strokeLinecap="round"
                    strokeDasharray="31.416"
                    strokeDashoffset="23.562"
                    fill="none"
                    className="animate-spin"
                    style={{
                        transformOrigin: '50% 50%',
                        animation: 'spin 1s linear infinite'
                    }}
                />
            </svg>
            <span className="sr-only">Loading...</span>
        </div>
    );
};

export default LoadingSpinner; 
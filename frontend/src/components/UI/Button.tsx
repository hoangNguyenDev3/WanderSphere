import React from 'react';
import LoadingSpinner from './LoadingSpinner';

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
    variant?: 'primary' | 'secondary' | 'outline' | 'ghost' | 'danger' | 'success';
    size?: 'sm' | 'md' | 'lg';
    isLoading?: boolean;
    leftIcon?: React.ReactNode;
    rightIcon?: React.ReactNode;
    fullWidth?: boolean;
}

const Button: React.FC<ButtonProps> = ({
    children,
    variant = 'primary',
    size = 'md',
    isLoading = false,
    leftIcon,
    rightIcon,
    fullWidth = false,
    className = '',
    disabled,
    ...props
}) => {
    const baseClasses = 'inline-flex items-center justify-center font-medium rounded-xl transition-all duration-200 btn-focus btn-modern relative overflow-hidden';

    const variants = {
        primary: 'bg-primary hover:bg-primary/90 text-primary-foreground shadow-lg hover:shadow-xl hover:shadow-primary/25',
        secondary: 'bg-secondary hover:bg-secondary/90 text-secondary-foreground',
        outline: 'border-2 border-border text-foreground hover:bg-muted hover:border-muted-foreground',
        ghost: 'text-muted-foreground hover:text-foreground hover:bg-muted',
        danger: 'bg-destructive hover:bg-destructive/90 text-destructive-foreground shadow-lg hover:shadow-xl hover:shadow-destructive/25',
        success: 'bg-success-600 hover:bg-success-700 dark:bg-success-500 dark:hover:bg-success-600 text-white shadow-lg hover:shadow-xl hover:shadow-success-500/25'
    };

    const sizes = {
        sm: 'px-3 py-1.5 text-sm h-8',
        md: 'px-4 py-2 text-sm h-10',
        lg: 'px-6 py-3 text-base h-12'
    };

    const classes = [
        baseClasses,
        variants[variant],
        sizes[size],
        fullWidth ? 'w-full' : '',
        (disabled || isLoading) ? 'opacity-50 cursor-not-allowed hover:shadow-none hover:transform-none' : 'hover:scale-105 active:scale-95',
        className
    ].filter(Boolean).join(' ');

    return (
        <button
            className={classes}
            disabled={disabled || isLoading}
            {...props}
        >
            {isLoading ? (
                <LoadingSpinner size="sm" className="mr-2" />
            ) : leftIcon ? (
                <span className="mr-2 flex-shrink-0">{leftIcon}</span>
            ) : null}

            <span className="truncate">{children}</span>

            {rightIcon && !isLoading && (
                <span className="ml-2 flex-shrink-0">{rightIcon}</span>
            )}
        </button>
    );
};

export default Button; 
import React, { forwardRef } from 'react';

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
    label?: string;
    error?: string;
    fullWidth?: boolean;
    leftIcon?: React.ReactNode;
    rightIcon?: React.ReactNode;
}

const Input = forwardRef<HTMLInputElement, InputProps>(({
    label,
    error,
    fullWidth = false,
    leftIcon,
    rightIcon,
    className = '',
    ...props
}, ref) => {
    const baseClasses = 'block w-full px-4 py-3 text-sm border rounded-xl transition-all duration-200 placeholder:text-muted-foreground input-modern';
    const errorClasses = error
        ? 'border-destructive focus:border-destructive focus:ring-destructive/20'
        : 'border-border focus:border-primary focus:ring-primary/20';

    const classes = [
        baseClasses,
        errorClasses,
        leftIcon ? 'pl-12' : '',
        rightIcon ? 'pr-12' : '',
        fullWidth ? 'w-full' : '',
        className
    ].filter(Boolean).join(' ');

    return (
        <div className={fullWidth ? 'w-full' : ''}>
            {label && (
                <label className="block text-sm font-medium text-foreground mb-2">
                    {label}
                </label>
            )}
            <div className="relative">
                {leftIcon && (
                    <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                        <span className="text-muted-foreground">{leftIcon}</span>
                    </div>
                )}
                <input
                    ref={ref}
                    className={classes}
                    {...props}
                />
                {rightIcon && (
                    <div className="absolute inset-y-0 right-0 pr-3 flex items-center">
                        <span className="text-muted-foreground">{rightIcon}</span>
                    </div>
                )}
            </div>
            {error && (
                <p className="mt-2 text-sm text-destructive animate-fadeIn">
                    {error}
                </p>
            )}
        </div>
    );
});

Input.displayName = 'Input';

export default Input; 
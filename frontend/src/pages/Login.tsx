import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { useAuth } from '../contexts/AuthContext';
import { useTheme } from '../contexts/ThemeContext';
import { LoginRequest } from '../types/api';
import Button from '../components/UI/Button';
import Input from '../components/UI/Input';
import { EyeIcon, EyeSlashIcon, SunIcon, MoonIcon } from '@heroicons/react/24/outline';

const Login: React.FC = () => {
    const { login } = useAuth();
    const { isDarkMode, toggleTheme } = useTheme();
    const navigate = useNavigate();
    const [showPassword, setShowPassword] = useState(false);
    const [isLoading, setIsLoading] = useState(false);

    const {
        register,
        handleSubmit,
        formState: { errors },
        setError
    } = useForm<LoginRequest>();

    const onSubmit = async (data: LoginRequest) => {
        setIsLoading(true);

        try {
            const success = await login(data);
            if (success) {
                navigate('/');
            } else {
                setError('root', {
                    message: 'Invalid username or password'
                });
            }
        } catch (error) {
            setError('root', {
                message: 'An error occurred. Please try again.'
            });
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-primary-50 via-background to-secondary-50 dark:from-gray-900 dark:via-background dark:to-gray-800 py-12 px-4 sm:px-6 lg:px-8 transition-colors duration-300">
            {/* Theme toggle */}
            <button
                onClick={toggleTheme}
                className="fixed top-4 right-4 p-3 rounded-full bg-card/80 backdrop-blur-sm border border-border shadow-lg hover:shadow-xl transition-all duration-200 btn-focus z-50"
                aria-label="Toggle theme"
            >
                {isDarkMode ? (
                    <SunIcon className="h-5 w-5 text-foreground" />
                ) : (
                    <MoonIcon className="h-5 w-5 text-foreground" />
                )}
            </button>

            <div className="max-w-md w-full space-y-8">
                <div className="text-center">
                    {/* Logo */}
                    <div className="flex justify-center mb-6">
                        <div className="flex items-center space-x-3">
                            <div className="w-12 h-12 bg-gradient-to-r from-primary-500 to-primary-600 rounded-2xl flex items-center justify-center shadow-lg">
                                <span className="text-white font-bold text-xl">W</span>
                            </div>
                            <h1 className="text-4xl font-bold bg-gradient-to-r from-primary-600 to-primary-500 bg-clip-text text-transparent">
                                WanderSphere
                            </h1>
                        </div>
                    </div>

                    <div className="card p-8 shadow-xl border-0 bg-card/50 backdrop-blur-sm">
                        <h2 className="text-3xl font-bold text-foreground mb-2">
                            Welcome back
                        </h2>
                        <p className="text-muted-foreground mb-8">
                            Sign in to continue your journey
                        </p>

                        <form className="space-y-6" onSubmit={handleSubmit(onSubmit)}>
                            <div className="space-y-4">
                                <Input
                                    label="Username"
                                    type="text"
                                    autoComplete="username"
                                    fullWidth
                                    error={errors.user_name?.message}
                                    className="input-modern"
                                    {...register('user_name', {
                                        required: 'Username is required',
                                        minLength: {
                                            value: 4,
                                            message: 'Username must be at least 4 characters'
                                        }
                                    })}
                                />

                                <Input
                                    label="Password"
                                    type={showPassword ? 'text' : 'password'}
                                    autoComplete="current-password"
                                    fullWidth
                                    error={errors.password?.message}
                                    className="input-modern"
                                    rightIcon={
                                        <button
                                            type="button"
                                            onClick={() => setShowPassword(!showPassword)}
                                            className="focus:outline-none text-muted-foreground hover:text-foreground transition-colors"
                                        >
                                            {showPassword ? (
                                                <EyeSlashIcon className="h-5 w-5" />
                                            ) : (
                                                <EyeIcon className="h-5 w-5" />
                                            )}
                                        </button>
                                    }
                                    {...register('password', {
                                        required: 'Password is required',
                                        minLength: {
                                            value: 4,
                                            message: 'Password must be at least 4 characters'
                                        }
                                    })}
                                />
                            </div>

                            {errors.root && (
                                <div className="bg-destructive/10 border border-destructive/20 rounded-lg p-3 text-destructive text-sm text-center animate-fadeIn">
                                    {errors.root.message}
                                </div>
                            )}

                            <Button
                                type="submit"
                                fullWidth
                                isLoading={isLoading}
                                disabled={isLoading}
                                className="btn-modern bg-gradient-to-r from-primary-600 to-primary-500 hover:from-primary-500 hover:to-primary-400 text-primary-foreground font-semibold py-3 shadow-lg"
                            >
                                {isLoading ? 'Signing in...' : 'Sign in'}
                            </Button>
                        </form>

                        <div className="mt-6 text-center">
                            <p className="text-sm text-muted-foreground">
                                Don't have an account?{' '}
                                <Link
                                    to="/signup"
                                    className="font-semibold text-primary hover:text-primary/80 transition-colors"
                                >
                                    Create one here
                                </Link>
                            </p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Login; 
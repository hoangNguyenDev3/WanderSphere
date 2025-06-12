import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { useAuth } from '../contexts/AuthContext';
import { CreateUserRequest } from '../types/api';
import Button from '../components/UI/Button';
import Input from '../components/UI/Input';
import { EyeIcon, EyeSlashIcon } from '@heroicons/react/24/outline';

const Signup: React.FC = () => {
    const { signup } = useAuth();
    const navigate = useNavigate();
    const [showPassword, setShowPassword] = useState(false);
    const [isLoading, setIsLoading] = useState(false);

    const {
        register,
        handleSubmit,
        formState: { errors },
        setError,
        watch
    } = useForm<CreateUserRequest>();

    const password = watch('password');

    const onSubmit = async (data: CreateUserRequest) => {
        setIsLoading(true);

        try {
            const success = await signup(data);
            if (success) {
                navigate('/login');
            } else {
                setError('root', {
                    message: 'Failed to create account. Please try again.'
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
        <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
            <div className="max-w-md w-full space-y-8">
                <div>
                    <div className="flex justify-center">
                        <h1 className="text-4xl font-bold text-primary-600">
                            WanderSphere
                        </h1>
                    </div>
                    <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
                        Create your account
                    </h2>
                    <p className="mt-2 text-center text-sm text-gray-600">
                        Or{' '}
                        <Link
                            to="/login"
                            className="font-medium text-primary-600 hover:text-primary-500"
                        >
                            sign in to your existing account
                        </Link>
                    </p>
                </div>

                <form className="mt-8 space-y-6" onSubmit={handleSubmit(onSubmit)}>
                    <div className="space-y-4">
                        <Input
                            label="Username"
                            type="text"
                            autoComplete="username"
                            fullWidth
                            error={errors.user_name?.message}
                            {...register('user_name', {
                                required: 'Username is required',
                                minLength: {
                                    value: 4,
                                    message: 'Username must be at least 4 characters'
                                },
                                pattern: {
                                    value: /^[a-zA-Z0-9_-]+$/,
                                    message: 'Username can only contain letters, numbers, underscores, and hyphens'
                                }
                            })}
                        />

                        <Input
                            label="Email"
                            type="email"
                            autoComplete="email"
                            fullWidth
                            error={errors.email?.message}
                            {...register('email', {
                                required: 'Email is required',
                                pattern: {
                                    value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i,
                                    message: 'Invalid email address'
                                }
                            })}
                        />

                        <div className="grid grid-cols-2 gap-4">
                            <Input
                                label="First Name"
                                type="text"
                                autoComplete="given-name"
                                fullWidth
                                error={errors.first_name?.message}
                                {...register('first_name', {
                                    required: 'First name is required'
                                })}
                            />

                            <Input
                                label="Last Name"
                                type="text"
                                autoComplete="family-name"
                                fullWidth
                                error={errors.last_name?.message}
                                {...register('last_name', {
                                    required: 'Last name is required'
                                })}
                            />
                        </div>

                        <Input
                            label="Date of Birth"
                            type="date"
                            fullWidth
                            error={errors.date_of_birth?.message}
                            {...register('date_of_birth', {
                                required: 'Date of birth is required',
                                validate: (value) => {
                                    const birthDate = new Date(value);
                                    const today = new Date();
                                    const age = today.getFullYear() - birthDate.getFullYear();

                                    if (age < 13) {
                                        return 'You must be at least 13 years old';
                                    }

                                    if (birthDate > today) {
                                        return 'Date of birth cannot be in the future';
                                    }

                                    return true;
                                }
                            })}
                        />

                        <Input
                            label="Password"
                            type={showPassword ? 'text' : 'password'}
                            autoComplete="new-password"
                            fullWidth
                            error={errors.password?.message}
                            rightIcon={
                                <button
                                    type="button"
                                    onClick={() => setShowPassword(!showPassword)}
                                    className="focus:outline-none"
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
                        <div className="text-red-600 text-sm text-center">
                            {errors.root.message}
                        </div>
                    )}

                    <div>
                        <Button
                            type="submit"
                            fullWidth
                            isLoading={isLoading}
                            disabled={isLoading}
                        >
                            Create Account
                        </Button>
                    </div>
                </form>

                <div className="text-center">
                    <p className="text-sm text-gray-600">
                        Already have an account?{' '}
                        <Link
                            to="/login"
                            className="font-medium text-primary-600 hover:text-primary-500"
                        >
                            Sign in here
                        </Link>
                    </p>
                </div>
            </div>
        </div>
    );
};

export default Signup; 
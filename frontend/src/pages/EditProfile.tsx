import React, { useState, useRef } from 'react';
import { useForm } from 'react-hook-form';
import { useMutation, useQueryClient } from 'react-query';
import { useNavigate } from 'react-router-dom';
import { toast } from 'react-hot-toast';
import { userAPI, postsAPI } from '../services/api';
import { useAuth } from '../contexts/AuthContext';
import { EditUserRequest } from '../types/api';
import Button from '../components/UI/Button';
import Input from '../components/UI/Input';
import {
    CameraIcon,
    ArrowLeftIcon,
} from '@heroicons/react/24/outline';

interface EditProfileForm {
    first_name: string;
    last_name: string;
    date_of_birth: string;
    password?: string;
    confirmPassword?: string;
}

const EditProfile: React.FC = () => {
    const navigate = useNavigate();
    const { user } = useAuth();
    const queryClient = useQueryClient();

    const [profileImage, setProfileImage] = useState<File | null>(null);
    const [coverImage, setCoverImage] = useState<File | null>(null);
    const [profilePreview, setProfilePreview] = useState<string | null>(null);
    const [coverPreview, setCoverPreview] = useState<string | null>(null);

    const profileInputRef = useRef<HTMLInputElement>(null);
    const coverInputRef = useRef<HTMLInputElement>(null);

    const { register, handleSubmit, formState: { errors }, watch } = useForm<EditProfileForm>({
        defaultValues: {
            first_name: user?.first_name || '',
            last_name: user?.last_name || '',
            date_of_birth: user?.date_of_birth ? user.date_of_birth.split('T')[0] : '',
        }
    });

    // Update profile mutation
    const updateProfileMutation = useMutation(
        (data: EditUserRequest) => userAPI.updateProfile(data),
        {
            onSuccess: () => {
                toast.success('Profile updated successfully!');
                queryClient.invalidateQueries(['profile', user?.user_id]);
                navigate(`/profile/${user?.user_id}`);
            },
            onError: (error: any) => {
                toast.error(error.response?.data?.error || 'Failed to update profile');
            },
        }
    );

    // Upload file to backend
    const uploadFileToBackend = async (file: File): Promise<string> => {
        try {
            // Create FormData for multipart upload
            const formData = new FormData();
            formData.append('file', file);

            // Upload to backend
            const uploadResponse = await fetch('/api/v1/binaries/upload', {
                method: 'POST',
                body: formData,
                credentials: 'include', // Include cookies for authentication
            });

            if (!uploadResponse.ok) {
                throw new Error('Failed to upload file');
            }

            const result = await uploadResponse.json();

            if (!result.success || !result.data) {
                throw new Error('Upload failed');
            }

            // Return the file URL
            return result.data.url;
        } catch (error) {
            console.error('Upload error:', error);
            throw error;
        }
    };

    // Handle image selection
    const handleImageSelect = (file: File, type: 'profile' | 'cover') => {
        if (!file.type.startsWith('image/')) {
            toast.error('Please select a valid image file');
            return;
        }

        if (file.size > 10 * 1024 * 1024) {
            toast.error('Image size must be less than 10MB');
            return;
        }

        const preview = URL.createObjectURL(file);

        if (type === 'profile') {
            setProfileImage(file);
            setProfilePreview(preview);
        } else {
            setCoverImage(file);
            setCoverPreview(preview);
        }
    };

    const onSubmit = async (data: EditProfileForm) => {
        try {
            let updates: EditUserRequest = {
                first_name: data.first_name,
                last_name: data.last_name,
                date_of_birth: data.date_of_birth,
            };

            // Add password if provided
            if (data.password) {
                if (data.password !== data.confirmPassword) {
                    toast.error('Passwords do not match');
                    return;
                }
                updates.password = data.password;
            }

            // Upload profile image if selected
            if (profileImage) {
                toast.loading('Uploading profile picture...');
                const profileUrl = await uploadFileToBackend(profileImage);
                updates.profile_picture = profileUrl;
                toast.dismiss();
            }

            // Upload cover image if selected
            if (coverImage) {
                toast.loading('Uploading cover photo...');
                const coverUrl = await uploadFileToBackend(coverImage);
                updates.cover_picture = coverUrl;
                toast.dismiss();
            }

            updateProfileMutation.mutate(updates);
        } catch (error) {
            toast.dismiss();
            toast.error('Failed to upload images. Please try again.');
        }
    };

    const isLoading = updateProfileMutation.isLoading;

    if (!user) {
        return null;
    }

    return (
        <div className="max-w-2xl mx-auto">
            <div className="card p-6">
                {/* Header */}
                <div className="flex items-center justify-between mb-6">
                    <div className="flex items-center space-x-3">
                        <button
                            onClick={() => navigate(-1)}
                            className="text-gray-400 hover:text-gray-600"
                        >
                            <ArrowLeftIcon className="w-6 h-6" />
                        </button>
                        <h1 className="text-2xl font-bold text-gray-900">Edit Profile</h1>
                    </div>
                </div>

                <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
                    {/* Cover Photo */}
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            Cover Photo
                        </label>
                        <div className="relative">
                            <div className="h-32 bg-gradient-to-r from-primary-400 to-primary-600 rounded-lg overflow-hidden">
                                {coverPreview ? (
                                    <img
                                        src={coverPreview}
                                        alt="Cover preview"
                                        className="w-full h-full object-cover"
                                    />
                                ) : user.cover_picture ? (
                                    <img
                                        src={user.cover_picture}
                                        alt="Current cover"
                                        className="w-full h-full object-cover"
                                    />
                                ) : null}
                            </div>
                            <button
                                type="button"
                                onClick={() => coverInputRef.current?.click()}
                                className="absolute bottom-2 right-2 bg-white bg-opacity-90 hover:bg-opacity-100 rounded-full p-2 shadow-lg transition-all"
                            >
                                <CameraIcon className="w-5 h-5 text-gray-600" />
                            </button>
                            <input
                                ref={coverInputRef}
                                type="file"
                                accept="image/*"
                                onChange={(e) => e.target.files?.[0] && handleImageSelect(e.target.files[0], 'cover')}
                                className="hidden"
                            />
                        </div>
                    </div>

                    {/* Profile Picture */}
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            Profile Picture
                        </label>
                        <div className="flex items-center space-x-4">
                            <div className="relative">
                                <div className="w-20 h-20 rounded-full overflow-hidden bg-gray-100">
                                    {profilePreview ? (
                                        <img
                                            src={profilePreview}
                                            alt="Profile preview"
                                            className="w-full h-full object-cover"
                                        />
                                    ) : user.profile_picture ? (
                                        <img
                                            src={user.profile_picture}
                                            alt="Current profile"
                                            className="w-full h-full object-cover"
                                        />
                                    ) : (
                                        <div className="w-full h-full bg-gray-300 flex items-center justify-center">
                                            <span className="text-gray-600 text-lg font-bold">
                                                {user.first_name[0]}{user.last_name[0]}
                                            </span>
                                        </div>
                                    )}
                                </div>
                                <button
                                    type="button"
                                    onClick={() => profileInputRef.current?.click()}
                                    className="absolute bottom-0 right-0 bg-primary-600 hover:bg-primary-700 rounded-full p-1.5 text-white shadow-lg"
                                >
                                    <CameraIcon className="w-3 h-3" />
                                </button>
                                <input
                                    ref={profileInputRef}
                                    type="file"
                                    accept="image/*"
                                    onChange={(e) => e.target.files?.[0] && handleImageSelect(e.target.files[0], 'profile')}
                                    className="hidden"
                                />
                            </div>
                            <div>
                                <p className="text-sm text-gray-600">Upload a new profile picture</p>
                                <p className="text-xs text-gray-400">JPG, PNG or GIF, max 10MB</p>
                            </div>
                        </div>
                    </div>

                    {/* Basic Information */}
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <Input
                            label="First Name"
                            type="text"
                            {...register('first_name', {
                                required: 'First name is required',
                                minLength: { value: 2, message: 'First name must be at least 2 characters' }
                            })}
                            error={errors.first_name?.message}
                        />

                        <Input
                            label="Last Name"
                            type="text"
                            {...register('last_name', {
                                required: 'Last name is required',
                                minLength: { value: 2, message: 'Last name must be at least 2 characters' }
                            })}
                            error={errors.last_name?.message}
                        />
                    </div>

                    <Input
                        label="Date of Birth"
                        type="date"
                        {...register('date_of_birth', {
                            required: 'Date of birth is required'
                        })}
                        error={errors.date_of_birth?.message}
                    />

                    {/* Password Section */}
                    <div className="border-t border-gray-200 pt-6">
                        <h3 className="text-lg font-medium text-gray-900 mb-4">Change Password</h3>
                        <p className="text-sm text-gray-600 mb-4">
                            Leave blank to keep your current password
                        </p>

                        <div className="space-y-4">
                            <Input
                                label="New Password"
                                type="password"
                                {...register('password', {
                                    minLength: { value: 6, message: 'Password must be at least 6 characters' }
                                })}
                                error={errors.password?.message}
                                placeholder="Enter new password"
                            />

                            {watch('password') && (
                                <Input
                                    label="Confirm New Password"
                                    type="password"
                                    {...register('confirmPassword', {
                                        validate: value =>
                                            !watch('password') || value === watch('password') || 'Passwords do not match'
                                    })}
                                    error={errors.confirmPassword?.message}
                                    placeholder="Confirm new password"
                                />
                            )}
                        </div>
                    </div>

                    {/* Actions */}
                    <div className="flex justify-end space-x-3 pt-6 border-t border-gray-200">
                        <Button
                            type="button"
                            variant="secondary"
                            onClick={() => navigate(-1)}
                            disabled={isLoading}
                        >
                            Cancel
                        </Button>
                        <Button
                            type="submit"
                            isLoading={isLoading}
                        >
                            Save Changes
                        </Button>
                    </div>
                </form>
            </div>
        </div>
    );
};

export default EditProfile; 
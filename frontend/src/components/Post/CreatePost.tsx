import React, { useState, useRef, useCallback } from 'react';
import { useForm } from 'react-hook-form';
import { useMutation } from 'react-query';
import { toast } from 'react-hot-toast';
import { postsAPI } from '../../services/api';
import { CreatePostRequest } from '../../types/api';
import Button from '../UI/Button';
import {
    PhotoIcon,
    XMarkIcon,
    FaceSmileIcon,
    MapPinIcon,
    UserGroupIcon,
} from '@heroicons/react/24/outline';

interface CreatePostProps {
    onPostCreated: () => void;
    onCancel: () => void;
}

interface CreatePostForm {
    content_text: string;
}

interface UploadedImage {
    file: File;
    preview: string;
    url?: string;
    uploading?: boolean;
}

const CreatePost: React.FC<CreatePostProps> = ({ onPostCreated, onCancel }) => {
    const [images, setImages] = useState<UploadedImage[]>([]);
    const [dragActive, setDragActive] = useState(false);
    const fileInputRef = useRef<HTMLInputElement>(null);

    const { register, handleSubmit, formState: { errors }, watch, reset } = useForm<CreatePostForm>();
    const contentText = watch('content_text', '');

    // Create post mutation
    const createPostMutation = useMutation(
        (data: CreatePostRequest) => postsAPI.createPost(data),
        {
            onSuccess: () => {
                toast.success('Post created successfully!');
                reset();
                setImages([]);
                onPostCreated();
            },
            onError: (error: any) => {
                toast.error(error.response?.data?.error || 'Failed to create post');
            },
        }
    );

    // File upload to backend
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

    // Handle file selection
    const handleFileSelect = useCallback(async (files: FileList) => {
        const validFiles: File[] = [];

        for (let i = 0; i < files.length; i++) {
            const file = files[i];

            // Validate file type
            if (!file.type.startsWith('image/')) {
                toast.error(`${file.name} is not an image file`);
                continue;
            }

            // Validate file size (10MB limit)
            if (file.size > 10 * 1024 * 1024) {
                toast.error(`${file.name} is too large. Maximum size is 10MB`);
                continue;
            }

            validFiles.push(file);
        }

        // Check total images limit (4 images max)
        if (images.length + validFiles.length > 4) {
            toast.error('You can only upload up to 4 images per post');
            const allowedCount = 4 - images.length;
            validFiles.splice(allowedCount);
        }

        // Create preview for each valid file
        const newImages: UploadedImage[] = validFiles.map(file => ({
            file,
            preview: URL.createObjectURL(file),
            uploading: false,
        }));

        setImages(prev => [...prev, ...newImages]);
    }, [images.length]);

    // Drag and drop handlers
    const handleDrag = (e: React.DragEvent) => {
        e.preventDefault();
        e.stopPropagation();
        if (e.type === "dragenter" || e.type === "dragover") {
            setDragActive(true);
        } else if (e.type === "dragleave") {
            setDragActive(false);
        }
    };

    const handleDrop = (e: React.DragEvent) => {
        e.preventDefault();
        e.stopPropagation();
        setDragActive(false);

        if (e.dataTransfer.files && e.dataTransfer.files.length > 0) {
            handleFileSelect(e.dataTransfer.files);
        }
    };

    // Remove image
    const removeImage = (index: number) => {
        setImages(prev => {
            const newImages = [...prev];
            URL.revokeObjectURL(newImages[index].preview);
            newImages.splice(index, 1);
            return newImages;
        });
    };

    // Handle form submission
    const onSubmit = async (data: CreatePostForm) => {
        if (!data.content_text.trim() && images.length === 0) {
            toast.error('Please add some content or images to your post');
            return;
        }

        try {
            // Upload images to S3
            const uploadPromises = images.map(async (image, index) => {
                setImages(prev => prev.map((img, i) =>
                    i === index ? { ...img, uploading: true } : img
                ));

                try {
                    const url = await uploadFileToBackend(image.file);
                    setImages(prev => prev.map((img, i) =>
                        i === index ? { ...img, url, uploading: false } : img
                    ));
                    return url;
                } catch (error) {
                    setImages(prev => prev.map((img, i) =>
                        i === index ? { ...img, uploading: false } : img
                    ));
                    throw error;
                }
            });

            const imageUrls = await Promise.all(uploadPromises);

            // Create post
            const postData: CreatePostRequest = {
                content_text: data.content_text.trim(),
                content_image_path: imageUrls.length > 0 ? imageUrls : undefined,
                visible: true,
            };

            createPostMutation.mutate(postData);
        } catch (error) {
            toast.error('Failed to upload images. Please try again.');
        }
    };

    const isLoading = createPostMutation.isLoading || images.some(img => img.uploading);

    return (
        <div className="post-instagram">
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
                {/* Header */}
                <div className="flex items-center justify-between p-4 border-b border-gray-200">
                    <h3 className="text-lg font-semibold text-gray-900">Create Post</h3>
                    <button
                        type="button"
                        onClick={onCancel}
                        className="text-gray-500 hover:text-gray-700 transition-colors"
                    >
                        <XMarkIcon className="w-6 h-6" />
                    </button>
                </div>

                {/* Text Content */}
                <div className="px-4">
                    <textarea
                        {...register('content_text', {
                            maxLength: { value: 2000, message: 'Post content cannot exceed 2000 characters' }
                        })}
                        placeholder="What's on your mind?"
                        className="w-full min-h-[120px] border-0 resize-none focus:ring-0 text-base placeholder-gray-400"
                        style={{ outline: 'none' }}
                    />
                    {errors.content_text && (
                        <p className="text-red-500 text-sm mt-1">{errors.content_text.message}</p>
                    )}
                    <div className="text-right text-xs text-gray-400 mt-1">
                        {contentText.length}/2000
                    </div>
                </div>

                {/* Image Upload Area */}
                <div className="px-4">
                    <div
                        className={`relative border-2 border-dashed rounded-lg p-6 transition-colors ${dragActive
                            ? 'border-blue-400 bg-blue-50'
                            : 'border-gray-300 hover:border-gray-400'
                            }`}
                        onDragEnter={handleDrag}
                        onDragLeave={handleDrag}
                        onDragOver={handleDrag}
                        onDrop={handleDrop}
                    >
                        <input
                            ref={fileInputRef}
                            type="file"
                            multiple
                            accept="image/*"
                            onChange={(e) => e.target.files && handleFileSelect(e.target.files)}
                            className="hidden"
                        />

                        {images.length === 0 ? (
                            <div className="text-center">
                                <PhotoIcon className="w-12 h-12 text-gray-400 mx-auto mb-3" />
                                <p className="text-gray-500 mb-2">Drag and drop images here, or click to select</p>
                                <button
                                    type="button"
                                    onClick={() => fileInputRef.current?.click()}
                                    className="text-primary-600 hover:text-primary-700 font-medium"
                                >
                                    Choose Files
                                </button>
                                <p className="text-xs text-gray-400 mt-2">
                                    Up to 4 images, max 10MB each
                                </p>
                            </div>
                        ) : (
                            <div className="grid grid-cols-2 gap-3">
                                {images.map((image, index) => (
                                    <div key={index} className="relative group">
                                        <img
                                            src={image.preview}
                                            alt={`Upload ${index + 1}`}
                                            className="w-full h-32 object-cover rounded-lg"
                                        />
                                        {image.uploading && (
                                            <div className="absolute inset-0 bg-black bg-opacity-50 flex items-center justify-center rounded-lg">
                                                <div className="animate-spin rounded-full h-6 w-6 border-2 border-white border-t-transparent"></div>
                                            </div>
                                        )}
                                        <button
                                            type="button"
                                            onClick={() => removeImage(index)}
                                            className="absolute top-2 right-2 w-6 h-6 bg-red-500 text-white rounded-full flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity"
                                        >
                                            <XMarkIcon className="w-4 h-4" />
                                        </button>
                                    </div>
                                ))}
                                {images.length < 4 && (
                                    <button
                                        type="button"
                                        onClick={() => fileInputRef.current?.click()}
                                        className="h-32 border-2 border-dashed border-gray-300 rounded-lg flex items-center justify-center text-gray-400 hover:border-gray-400 hover:text-gray-500 transition-colors"
                                    >
                                        <PhotoIcon className="w-8 h-8" />
                                    </button>
                                )}
                            </div>
                        )}
                    </div>
                </div>

                {/* Post Options */}
                <div className="flex items-center justify-between px-4 py-3 border-t border-gray-200">
                    <div className="flex items-center space-x-6">
                        <button
                            type="button"
                            onClick={() => fileInputRef.current?.click()}
                            className="text-gray-600 hover:text-gray-800 transition-colors"
                        >
                            <PhotoIcon className="w-6 h-6" />
                        </button>
                        <button
                            type="button"
                            className="text-gray-600 hover:text-gray-800 transition-colors"
                        >
                            <FaceSmileIcon className="w-6 h-6" />
                        </button>
                        <button
                            type="button"
                            className="text-gray-600 hover:text-gray-800 transition-colors"
                        >
                            <MapPinIcon className="w-6 h-6" />
                        </button>
                        <button
                            type="button"
                            className="text-gray-600 hover:text-gray-800 transition-colors"
                        >
                            <UserGroupIcon className="w-6 h-6" />
                        </button>
                    </div>

                    <div className="flex items-center space-x-3">
                        <button
                            type="button"
                            className="btn-secondary-instagram"
                            onClick={onCancel}
                            disabled={isLoading}
                        >
                            Cancel
                        </button>
                        <button
                            type="submit"
                            className="btn-instagram"
                            disabled={isLoading || (!contentText.trim() && images.length === 0)}
                        >
                            {isLoading ? 'Posting...' : 'Post'}
                        </button>
                    </div>
                </div>
            </form>
        </div>
    );
};

export default CreatePost; 
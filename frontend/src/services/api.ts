import axios, { AxiosResponse } from 'axios';
import {
    LoginRequest,
    LoginResponse,
    CreateUserRequest,
    EditUserRequest,
    User,
    Post,
    CreatePostRequest,
    CreatePostResponse,
    EditPostRequest,
    CreatePostCommentRequest,
    UserFollowerResponse,
    UserFollowingResponse,
    UserPostsResponse,
    NewsfeedResponse,
    GetS3PresignedUrlRequest,
    GetS3PresignedUrlResponse,
    MessageResponse,
    ApiResponse,
} from '../types/api';

// Create axios instance with base configuration
const api = axios.create({
    baseURL: '/api/v1',
    withCredentials: true,
    headers: {
        'Content-Type': 'application/json',
    },
});

// Response interceptor to handle errors
api.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response?.status === 401) {
            // Handle unauthorized access
            localStorage.removeItem('user');
            window.location.href = '/login';
        }
        return Promise.reject(error);
    }
);

// Helper function to handle API responses
const handleResponse = <T>(response: AxiosResponse<T>): ApiResponse<T> => ({
    data: response.data,
    status: response.status,
});

const handleError = (error: any): ApiResponse => ({
    error: error.response?.data || { error: 'Network Error', message: error.message },
    status: error.response?.status || 500,
});

// Authentication API
export const authAPI = {
    login: async (data: LoginRequest): Promise<ApiResponse<LoginResponse>> => {
        try {
            const response = await api.post<LoginResponse>('/users/login', data);
            return handleResponse(response);
        } catch (error) {
            return handleError(error);
        }
    },

    signup: async (data: CreateUserRequest): Promise<ApiResponse<MessageResponse>> => {
        try {
            const response = await api.post<MessageResponse>('/users/signup', data);
            return handleResponse(response);
        } catch (error) {
            return handleError(error);
        }
    },

    logout: async (): Promise<ApiResponse<MessageResponse>> => {
        try {
            const response = await api.post<MessageResponse>('/auth/logout');
            return handleResponse(response);
        } catch (error) {
            return handleError(error);
        }
    },
};

// User API
export const userAPI = {
    getProfile: async (userId: number): Promise<ApiResponse<User>> => {
        try {
            const response = await api.get<User>(`/users/${userId}`);
            return handleResponse(response);
        } catch (error) {
            return handleError(error);
        }
    },

    updateProfile: async (data: EditUserRequest): Promise<ApiResponse<MessageResponse>> => {
        try {
            const response = await api.put<MessageResponse>('/users/edit', data);
            return handleResponse(response);
        } catch (error) {
            return handleError(error);
        }
    },
};

// Posts API
export const postsAPI = {
    createPost: async (data: CreatePostRequest): Promise<ApiResponse<CreatePostResponse>> => {
        try {
            const response = await api.post<CreatePostResponse>('/posts', data);
            return handleResponse(response);
        } catch (error) {
            return handleError(error);
        }
    },

    getPost: async (postId: number): Promise<ApiResponse<Post>> => {
        try {
            const response = await api.get<Post>(`/posts/${postId}`);
            return handleResponse(response);
        } catch (error) {
            return handleError(error);
        }
    },

    updatePost: async (postId: number, data: EditPostRequest): Promise<ApiResponse<MessageResponse>> => {
        try {
            const response = await api.put<MessageResponse>(`/posts/${postId}`, data);
            return handleResponse(response);
        } catch (error) {
            return handleError(error);
        }
    },

    deletePost: async (postId: number): Promise<ApiResponse<MessageResponse>> => {
        try {
            const response = await api.delete<MessageResponse>(`/posts/${postId}`);
            return handleResponse(response);
        } catch (error) {
            return handleError(error);
        }
    },

    likePost: async (postId: number): Promise<ApiResponse<MessageResponse>> => {
        try {
            const response = await api.post<MessageResponse>(`/posts/${postId}/likes`);
            return handleResponse(response);
        } catch (error) {
            return handleError(error);
        }
    },

    commentOnPost: async (postId: number, data: CreatePostCommentRequest): Promise<ApiResponse<MessageResponse>> => {
        try {
            const response = await api.post<MessageResponse>(`/posts/${postId}`, data);
            return handleResponse(response);
        } catch (error) {
            return handleError(error);
        }
    },

    getPresignedUrl: async (data: GetS3PresignedUrlRequest): Promise<ApiResponse<GetS3PresignedUrlResponse>> => {
        try {
            const response = await api.get<GetS3PresignedUrlResponse>('/posts/url', {
                params: data,
            });
            return handleResponse(response);
        } catch (error) {
            return handleError(error);
        }
    },
};

// Social API (Friends/Following)
export const socialAPI = {
    getFollowers: async (userId: number): Promise<ApiResponse<UserFollowerResponse>> => {
        try {
            const response = await api.get<UserFollowerResponse>(`/friends/${userId}/followers`);
            return handleResponse(response);
        } catch (error) {
            return handleError(error);
        }
    },

    getFollowing: async (userId: number): Promise<ApiResponse<UserFollowingResponse>> => {
        try {
            const response = await api.get<UserFollowingResponse>(`/friends/${userId}/followings`);
            return handleResponse(response);
        } catch (error) {
            return handleError(error);
        }
    },

    getUserPosts: async (userId: number): Promise<ApiResponse<UserPostsResponse>> => {
        try {
            const response = await api.get<UserPostsResponse>(`/friends/${userId}/posts`);
            return handleResponse(response);
        } catch (error) {
            return handleError(error);
        }
    },

    followUser: async (userId: number): Promise<ApiResponse<MessageResponse>> => {
        try {
            const response = await api.post<MessageResponse>(`/friends/${userId}`);
            return handleResponse(response);
        } catch (error) {
            return handleError(error);
        }
    },

    unfollowUser: async (userId: number): Promise<ApiResponse<MessageResponse>> => {
        try {
            const response = await api.delete<MessageResponse>(`/friends/${userId}`);
            return handleResponse(response);
        } catch (error) {
            return handleError(error);
        }
    },
};

// Newsfeed API
export const newsfeedAPI = {
    getNewsfeed: async (): Promise<ApiResponse<NewsfeedResponse>> => {
        try {
            const response = await api.get<NewsfeedResponse>('/newsfeed');
            return handleResponse(response);
        } catch (error) {
            return handleError(error);
        }
    },
};

// File Upload API
export const fileAPI = {
    uploadFile: async (file: File, presignedUrl: string): Promise<boolean> => {
        try {
            await axios.put(presignedUrl, file, {
                headers: {
                    'Content-Type': file.type,
                },
            });
            return true;
        } catch (error) {
            console.error('File upload error:', error);
            return false;
        }
    },
};

export default api; 
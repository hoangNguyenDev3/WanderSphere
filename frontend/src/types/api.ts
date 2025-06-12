// User types
export interface User {
    user_id: number;
    user_name: string;
    first_name: string;
    last_name: string;
    date_of_birth: string;
    email: string;
    profile_picture?: string;
    cover_picture?: string;
}

// Auth types
export interface LoginRequest {
    user_name: string;
    password: string;
}

export interface LoginResponse {
    message: string;
    user: User;
}

export interface CreateUserRequest {
    user_name: string;
    password: string;
    first_name: string;
    last_name: string;
    date_of_birth: string;
    email: string;
}

export interface EditUserRequest {
    password?: string;
    first_name?: string;
    last_name?: string;
    date_of_birth?: string;
    profile_picture?: string;
    cover_picture?: string;
}

// Post types
export interface Comment {
    comment_id: number;
    user_id: number;
    post_id: number;
    content_text: string;
}

export interface Post {
    post_id: number;
    user_id: number;
    content_text: string;
    content_image_path: string[];
    created_at: string;
    comments: Comment[];
    users_liked: number[];
}

export interface CreatePostRequest {
    content_text: string;
    content_image_path?: string[];
    visible?: boolean;
}

export interface CreatePostResponse {
    message: string;
    post_id: number;
}

export interface EditPostRequest {
    content_text?: string;
    content_image_path?: string[];
    visible?: boolean;
}

export interface CreatePostCommentRequest {
    content_text: string;
}

// Social types
export interface UserFollowerResponse {
    followers_ids: number[];
}

export interface UserFollowingResponse {
    followings_ids: number[];
}

export interface UserPostsResponse {
    posts_ids: number[];
}

export interface NewsfeedResponse {
    posts_ids: number[];
}

// File upload types
export interface GetS3PresignedUrlRequest {
    file_name: string;
    file_type: string;
}

export interface GetS3PresignedUrlResponse {
    url: string;
    expiration_time: string;
}

// Generic response types
export interface MessageResponse {
    message: string;
    status?: string;
}

export interface ErrorResponse {
    error: string;
    message?: string;
    code?: number;
}

// API response wrapper
export interface ApiResponse<T = any> {
    data?: T;
    error?: ErrorResponse;
    status: number;
}

// Component types
export interface UserProfile extends User {
    followers_count?: number;
    following_count?: number;
    posts_count?: number;
}

export interface PostWithUser extends Post {
    user?: User;
    is_liked?: boolean;
}

export interface FeedPost extends PostWithUser {
    time_ago?: string;
} 
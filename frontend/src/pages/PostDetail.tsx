import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from 'react-query';
import { useForm } from 'react-hook-form';
import { toast } from 'react-hot-toast';
import { postsAPI, userAPI } from '../services/api';
import { useAuth } from '../contexts/AuthContext';
import { Post, PostWithUser, Comment, CreatePostCommentRequest, User } from '../types/api';
import PostCard from '../components/Post/PostCard';
import LoadingSpinner from '../components/UI/LoadingSpinner';
import Button from '../components/UI/Button';
import {
    ArrowLeftIcon,
    PaperAirplaneIcon,
} from '@heroicons/react/24/outline';

interface CommentForm {
    content_text: string;
}

interface CommentWithUser extends Comment {
    user?: User;
}

const PostDetail: React.FC = () => {
    const { postId } = useParams<{ postId: string }>();
    const navigate = useNavigate();
    const { user: currentUser } = useAuth();
    const queryClient = useQueryClient();

    const [post, setPost] = useState<PostWithUser | null>(null);
    const [comments, setComments] = useState<CommentWithUser[]>([]);

    const { register, handleSubmit, reset, formState: { errors } } = useForm<CommentForm>();

    // Fetch post
    const {
        data: postData,
        isLoading: isLoadingPost,
        error: postError,
    } = useQuery(
        ['post', postId],
        () => postId ? postsAPI.getPost(parseInt(postId)) : null,
        {
            enabled: !!postId,
        }
    );

    // Fetch post author when we have post data
    useEffect(() => {
        const fetchPostWithUser = async () => {
            if (!postData?.data) {
                setPost(null);
                return;
            }

            try {
                // Fetch the user who created the post
                const userResponse = await userAPI.getProfile(postData.data.user_id);

                setPost({
                    ...postData.data,
                    user: userResponse.data || undefined,
                });

                // Fetch comment users
                if (postData.data.comments) {
                    const commentPromises = postData.data.comments.map(async (comment) => {
                        try {
                            const commentUserResponse = await userAPI.getProfile(comment.user_id);
                            return {
                                ...comment,
                                user: commentUserResponse.data || undefined,
                            };
                        } catch (error) {
                            return {
                                ...comment,
                                user: undefined,
                            };
                        }
                    });

                    const commentsWithUsers = await Promise.all(commentPromises);
                    setComments(commentsWithUsers);
                }
            } catch (error) {
                console.error('Error fetching post user:', error);
                setPost({
                    ...postData.data,
                    user: undefined,
                });
            }
        };

        fetchPostWithUser();
    }, [postData]);

    // Add comment mutation
    const addCommentMutation = useMutation(
        (data: CreatePostCommentRequest) => {
            if (!postId) throw new Error('No post ID');
            return postsAPI.commentOnPost(parseInt(postId), data);
        },
        {
            onSuccess: () => {
                toast.success('Comment added successfully!');
                reset();
                queryClient.invalidateQueries(['post', postId]);
            },
            onError: (error: any) => {
                toast.error(error.response?.data?.error || 'Failed to add comment');
            },
        }
    );

    const onSubmitComment = (data: CommentForm) => {
        if (!data.content_text.trim()) return;

        addCommentMutation.mutate({
            content_text: data.content_text.trim(),
        });
    };

    const handlePostLike = (postId: number) => {
        // Optimistic update handled in PostCard
    };

    const handlePostComment = (postId: number) => {
        // Already on post detail page
    };

    const handlePostEdit = (postId: number) => {
        // Open edit modal or navigate to edit page
        console.log('Edit post:', postId);
    };

    const handlePostDelete = (postId: number) => {
        // Navigate back after deletion
        navigate(-1);
    };

    if (isLoadingPost) {
        return (
            <div className="flex justify-center items-center min-h-64">
                <LoadingSpinner size="lg" />
            </div>
        );
    }

    if (postError || !post) {
        return (
            <div className="text-center py-12">
                <p className="text-gray-500 mb-4">Post not found</p>
                <Button onClick={() => navigate(-1)}>Go Back</Button>
            </div>
        );
    }

    return (
        <div className="max-w-2xl mx-auto">
            {/* Header */}
            <div className="flex items-center space-x-3 mb-6">
                <button
                    onClick={() => navigate(-1)}
                    className="text-gray-400 hover:text-gray-600"
                >
                    <ArrowLeftIcon className="w-6 h-6" />
                </button>
                <h1 className="text-xl font-semibold text-gray-900">Post</h1>
            </div>

            {/* Post */}
            <div className="mb-6">
                <PostCard
                    post={post}
                    onLike={handlePostLike}
                    onComment={handlePostComment}
                    onEdit={handlePostEdit}
                    onDelete={handlePostDelete}
                />
            </div>

            {/* Comments Section */}
            <div className="card p-6">
                <h2 className="text-lg font-semibold text-gray-900 mb-4">
                    Comments ({comments.length})
                </h2>

                {/* Add Comment Form */}
                {currentUser && (
                    <form onSubmit={handleSubmit(onSubmitComment)} className="mb-6">
                        <div className="flex space-x-3">
                            <div className="flex-shrink-0">
                                {currentUser.profile_picture ? (
                                    <img
                                        src={currentUser.profile_picture}
                                        alt={currentUser.user_name}
                                        className="w-8 h-8 rounded-full object-cover"
                                    />
                                ) : (
                                    <div className="w-8 h-8 rounded-full bg-primary-100 flex items-center justify-center">
                                        <span className="text-primary-600 text-sm font-medium">
                                            {currentUser.first_name[0]}{currentUser.last_name[0]}
                                        </span>
                                    </div>
                                )}
                            </div>
                            <div className="flex-1">
                                <textarea
                                    {...register('content_text', {
                                        required: 'Comment cannot be empty',
                                        maxLength: { value: 500, message: 'Comment cannot exceed 500 characters' }
                                    })}
                                    placeholder="Write a comment..."
                                    rows={3}
                                    className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent resize-none"
                                />
                                {errors.content_text && (
                                    <p className="text-red-500 text-sm mt-1">{errors.content_text.message}</p>
                                )}
                                <div className="flex justify-between items-center mt-2">
                                    <span className="text-xs text-gray-400">
                                        Max 500 characters
                                    </span>
                                    <Button
                                        type="submit"
                                        size="sm"
                                        isLoading={addCommentMutation.isLoading}
                                        disabled={addCommentMutation.isLoading}
                                    >
                                        <PaperAirplaneIcon className="w-4 h-4 mr-1" />
                                        Comment
                                    </Button>
                                </div>
                            </div>
                        </div>
                    </form>
                )}

                {/* Comments List */}
                <div className="space-y-4">
                    {comments.length === 0 ? (
                        <div className="text-center py-8">
                            <p className="text-gray-500">No comments yet.</p>
                            <p className="text-gray-400 text-sm mt-1">
                                Be the first to comment on this post!
                            </p>
                        </div>
                    ) : (
                        comments.map((comment) => (
                            <div key={comment.comment_id} className="flex space-x-3">
                                <div className="flex-shrink-0">
                                    {comment.user?.profile_picture ? (
                                        <img
                                            src={comment.user.profile_picture}
                                            alt={comment.user.user_name}
                                            className="w-8 h-8 rounded-full object-cover"
                                        />
                                    ) : (
                                        <div className="w-8 h-8 rounded-full bg-gray-100 flex items-center justify-center">
                                            <span className="text-gray-600 text-sm font-medium">
                                                {comment.user ?
                                                    `${comment.user.first_name[0]}${comment.user.last_name[0]}` :
                                                    '?'
                                                }
                                            </span>
                                        </div>
                                    )}
                                </div>
                                <div className="flex-1">
                                    <div className="bg-gray-50 rounded-lg p-3">
                                        <div className="flex items-center space-x-2 mb-1">
                                            <h4 className="font-medium text-gray-900">
                                                {comment.user ?
                                                    `${comment.user.first_name} ${comment.user.last_name}` :
                                                    'Unknown User'
                                                }
                                            </h4>
                                            {comment.user && (
                                                <span className="text-gray-500 text-sm">
                                                    @{comment.user.user_name}
                                                </span>
                                            )}
                                        </div>
                                        <p className="text-gray-700">{comment.content_text}</p>
                                    </div>
                                    <div className="flex items-center space-x-4 mt-2 text-sm text-gray-500">
                                        <button className="hover:text-primary-600 transition-colors">
                                            Like
                                        </button>
                                        <button className="hover:text-primary-600 transition-colors">
                                            Reply
                                        </button>
                                        {comment.user_id === currentUser?.user_id && (
                                            <button className="hover:text-red-600 transition-colors">
                                                Delete
                                            </button>
                                        )}
                                    </div>
                                </div>
                            </div>
                        ))
                    )}
                </div>
            </div>
        </div>
    );
};

export default PostDetail; 
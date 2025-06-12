import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { format } from 'date-fns';
import { PostWithUser, User } from '../../types/api';
import { useAuth } from '../../contexts/AuthContext';
import { postsAPI } from '../../services/api';
import {
    HeartIcon,
    ChatBubbleOvalLeftIcon,
    ShareIcon,
    EllipsisHorizontalIcon,
} from '@heroicons/react/24/outline';
import {
    HeartIcon as HeartIconSolid,
} from '@heroicons/react/24/solid';
import { Menu, Transition } from '@headlessui/react';
import toast from 'react-hot-toast';

interface PostCardProps {
    post: PostWithUser;
    onLike?: (postId: number) => void;
    onComment?: (postId: number) => void;
    onEdit?: (postId: number) => void;
    onDelete?: (postId: number) => void;
}

const PostCard: React.FC<PostCardProps> = ({
    post,
    onLike,
    onComment,
    onEdit,
    onDelete,
}) => {
    const { user: currentUser } = useAuth();
    const [isLiking, setIsLiking] = useState(false);
    const [localPost, setLocalPost] = useState(post);

    const isOwnPost = currentUser?.user_id === post.user_id;
    const isLiked = localPost.users_liked.includes(currentUser?.user_id || 0);

    const getUserInitials = (user: User) => {
        return `${user.first_name.charAt(0)}${user.last_name.charAt(0)}`.toUpperCase();
    };

    const formatDate = (dateString: string) => {
        try {
            return format(new Date(dateString), 'MMM d, yyyy • h:mm a');
        } catch {
            return 'Unknown date';
        }
    };

    const handleLike = async () => {
        if (isLiking) return;

        setIsLiking(true);
        try {
            const response = await postsAPI.likePost(post.post_id);

            if (response.data) {
                // Toggle like status optimistically
                const updatedPost = {
                    ...localPost,
                    users_liked: isLiked
                        ? localPost.users_liked.filter(id => id !== currentUser?.user_id)
                        : [...localPost.users_liked, currentUser?.user_id || 0]
                };
                setLocalPost(updatedPost);

                if (onLike) {
                    onLike(post.post_id);
                }
            }
        } catch (error) {
            toast.error('Failed to like post');
        } finally {
            setIsLiking(false);
        }
    };

    const handleDelete = async () => {
        // Use toast with confirmation action instead of window.confirm
        toast((t) => (
            <div className="flex flex-col space-y-3">
                <span className="font-medium">Delete this post?</span>
                <span className="text-sm text-gray-600">This action cannot be undone.</span>
                <div className="flex space-x-2">
                    <button
                        onClick={async () => {
                            toast.dismiss(t.id);
                            try {
                                const response = await postsAPI.deletePost(post.post_id);
                                if (response.data) {
                                    toast.success('Post deleted successfully');
                                    if (onDelete) {
                                        onDelete(post.post_id);
                                    }
                                }
                            } catch (error) {
                                toast.error('Failed to delete post');
                            }
                        }}
                        className="px-3 py-1 bg-red-600 text-white text-sm rounded hover:bg-red-700 transition-colors"
                    >
                        Delete
                    </button>
                    <button
                        onClick={() => toast.dismiss(t.id)}
                        className="px-3 py-1 bg-gray-200 text-gray-800 text-sm rounded hover:bg-gray-300 transition-colors"
                    >
                        Cancel
                    </button>
                </div>
            </div>
        ), {
            duration: 8000,
            style: {
                maxWidth: '400px',
            }
        });
    };

    return (
        <div className="post-instagram animate-fadeIn">
            {/* Header */}
            <div className="flex items-center justify-between p-4">
                <div className="flex items-center space-x-3">
                    <Link to={`/profile/${post.user_id}`}>
                        <div className="story-border">
                            {post.user?.profile_picture ? (
                                <img
                                    src={post.user.profile_picture}
                                    alt={`${post.user.first_name} ${post.user.last_name}`}
                                    className="w-8 h-8 rounded-full object-cover"
                                />
                            ) : (
                                <div className="w-8 h-8 rounded-full bg-gray-300 flex items-center justify-center text-gray-600 text-xs font-medium">
                                    {post.user ? getUserInitials(post.user) : 'U'}
                                </div>
                            )}
                        </div>
                    </Link>

                    <div>
                        <Link
                            to={`/profile/${post.user_id}`}
                            className="text-username hover:text-gray-600 transition-colors"
                        >
                            {post.user ? post.user.user_name : 'unknown_user'}
                        </Link>
                        <p className="text-muted-instagram">
                            {formatDate(post.created_at)}
                        </p>
                    </div>
                </div>

                {isOwnPost && (
                    <Menu as="div" className="relative">
                        <Menu.Button className="p-1 text-gray-600 hover:text-gray-800 rounded-full hover:bg-gray-50 transition-colors">
                            <EllipsisHorizontalIcon className="w-5 h-5" />
                        </Menu.Button>

                        <Transition
                            as={React.Fragment}
                            enter="transition ease-out duration-100"
                            enterFrom="transform opacity-0 scale-95"
                            enterTo="transform opacity-100 scale-100"
                            leave="transition ease-in duration-75"
                            leaveFrom="transform opacity-100 scale-100"
                            leaveTo="transform opacity-0 scale-95"
                        >
                            <Menu.Items className="absolute right-0 mt-2 w-48 origin-top-right bg-white rounded-lg shadow-lg ring-1 ring-gray-200 focus:outline-none">
                                <div className="py-2">
                                    <Menu.Item>
                                        {({ active }) => (
                                            <button
                                                onClick={() => onEdit && onEdit(post.post_id)}
                                                className={`${active ? 'bg-gray-50' : ''
                                                    } block w-full text-left px-4 py-3 text-sm text-gray-700 hover:bg-gray-50 transition-colors`}
                                            >
                                                Edit Post
                                            </button>
                                        )}
                                    </Menu.Item>
                                    <Menu.Item>
                                        {({ active }) => (
                                            <button
                                                onClick={handleDelete}
                                                className={`${active ? 'bg-red-50' : ''
                                                    } block w-full text-left px-4 py-3 text-sm text-red-600 hover:bg-red-50 transition-colors`}
                                            >
                                                Delete Post
                                            </button>
                                        )}
                                    </Menu.Item>
                                </div>
                            </Menu.Items>
                        </Transition>
                    </Menu>
                )}
            </div>

            {/* Images */}
            {post.content_image_path && post.content_image_path.length > 0 && (
                <div className="mb-3">
                    <div className={`grid gap-1 ${post.content_image_path.length === 1 ? 'grid-cols-1' : 'grid-cols-2'}`}>
                        {post.content_image_path.map((imagePath, index) => (
                            <img
                                key={index}
                                src={imagePath}
                                alt={`Content ${index + 1}`}
                                className={`w-full object-cover ${post.content_image_path.length === 1
                                    ? 'max-h-72 aspect-[4/3]'
                                    : 'aspect-square max-h-36'
                                    }`}
                                loading="lazy"
                            />
                        ))}
                    </div>
                </div>
            )}

            {/* Content */}
            {post.content_text && (
                <div className="px-4 pb-3">
                    <p className="text-caption">
                        <span className="text-username mr-2">{post.user ? post.user.user_name : 'unknown_user'}</span>
                        {post.content_text}
                    </p>
                </div>
            )}

            {/* Actions */}
            <div className="px-4 pb-3">
                <div className="flex items-center justify-between mb-2">
                    <div className="flex items-center space-x-4">
                        <button
                            onClick={handleLike}
                            disabled={isLiking}
                            className={`transition-transform duration-150 ${isLiked ? 'scale-110' : 'hover:scale-105'}`}
                        >
                            {isLiked ? (
                                <HeartIconSolid className="w-6 h-6 text-red-500" />
                            ) : (
                                <HeartIcon className="w-6 h-6 text-gray-700 hover:text-gray-500" />
                            )}
                        </button>

                        <button
                            onClick={() => onComment && onComment(post.post_id)}
                            className="hover:scale-105 transition-transform duration-150"
                        >
                            <ChatBubbleOvalLeftIcon className="w-6 h-6 text-gray-700 hover:text-gray-500" />
                        </button>

                        <button className="hover:scale-105 transition-transform duration-150">
                            <ShareIcon className="w-6 h-6 text-gray-700 hover:text-gray-500" />
                        </button>
                    </div>
                </div>

                {/* Like count */}
                {localPost.users_liked.length > 0 && (
                    <p className="text-username mb-1">
                        {localPost.users_liked.length} {localPost.users_liked.length === 1 ? 'like' : 'likes'}
                    </p>
                )}

                {/* Comments preview */}
                {post.comments.length > 0 && (
                    <p className="text-muted-instagram">
                        View all {post.comments.length} comments
                    </p>
                )}
            </div>
        </div>
    );
};

export default PostCard; 
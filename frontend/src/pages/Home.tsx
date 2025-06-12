import React, { useState, useEffect } from 'react';
import { useQuery } from 'react-query';
import { newsfeedAPI, postsAPI, userAPI } from '../services/api';
import { PostWithUser, Post } from '../types/api';
import PostCard from '../components/Post/PostCard';
import CreatePost from '../components/Post/CreatePost';
import LoadingSpinner from '../components/UI/LoadingSpinner';
import { PlusIcon } from '@heroicons/react/24/outline';
import { useAuth } from '../contexts/AuthContext';

const Home: React.FC = () => {
    const { user } = useAuth();
    const [posts, setPosts] = useState<PostWithUser[]>([]);
    const [showCreatePost, setShowCreatePost] = useState(false);

    // Fetch newsfeed
    const {
        data: newsfeedData,
        isLoading: isLoadingNewsfeed,
        error: newsfeedError,
        refetch: refetchNewsfeed
    } = useQuery(
        'newsfeed',
        () => newsfeedAPI.getNewsfeed(),
        {
            refetchOnMount: true,
            staleTime: 5 * 60 * 1000, // 5 minutes
        }
    );

    // Fetch individual posts when we have post IDs
    useEffect(() => {
        const fetchPosts = async () => {
            const postIds = newsfeedData?.data?.posts_ids || [];

            if (postIds.length === 0) {
                setPosts([]);
                return;
            }

            // Set loading state while fetching posts
            setPosts([]);

            try {
                const postPromises = postIds.map(async (postId) => {
                    const response = await postsAPI.getPost(postId);
                    return response.data;
                });

                const fetchedPosts = await Promise.all(postPromises);
                const validPosts = fetchedPosts.filter(Boolean) as Post[];

                // Fetch user information for each post
                const postsWithUserPromises = validPosts.map(async (post) => {
                    try {
                        const userResponse = await userAPI.getProfile(post.user_id);

                        if (userResponse.data && !userResponse.error) {
                            return {
                                ...post,
                                user: userResponse.data,
                            } as PostWithUser;
                        } else {
                            console.error(`Failed to fetch user ${post.user_id}:`, userResponse);
                            // Create a fallback user with basic info
                            const fallbackUser = {
                                user_id: post.user_id,
                                user_name: `user${post.user_id}`,
                                first_name: 'User',
                                last_name: `${post.user_id}`,
                                date_of_birth: '1990-01-01',
                                email: `user${post.user_id}@example.com`,
                            };
                            return {
                                ...post,
                                user: fallbackUser,
                            } as PostWithUser;
                        }
                    } catch (error) {
                        console.error(`Error fetching user ${post.user_id}:`, error);
                        // Create a fallback user with basic info
                        const fallbackUser = {
                            user_id: post.user_id,
                            user_name: `user${post.user_id}`,
                            first_name: 'User',
                            last_name: `${post.user_id}`,
                            date_of_birth: '1990-01-01',
                            email: `user${post.user_id}@example.com`,
                        };
                        return {
                            ...post,
                            user: fallbackUser,
                        } as PostWithUser;
                    }
                });

                const postsWithUser = await Promise.all(postsWithUserPromises);
                setPosts(postsWithUser);
            } catch (error) {
                console.error('Error fetching posts:', error);
            }
        };

        fetchPosts();
    }, [newsfeedData?.data?.posts_ids]);

    const handlePostCreated = () => {
        setShowCreatePost(false);
        refetchNewsfeed();
    };

    const handlePostLike = (postId: number) => {
        // Optimistic update handled in PostCard
        // Could trigger a refetch here if needed
    };

    const handlePostComment = (postId: number) => {
        // Navigate to post detail or open comment modal
        console.log('Comment on post:', postId);
    };

    const handlePostEdit = (postId: number) => {
        // Open edit modal or navigate to edit page
        console.log('Edit post:', postId);
    };

    const handlePostDelete = (postId: number) => {
        // Remove from local state
        setPosts(prevPosts => prevPosts.filter(post => post.post_id !== postId));
        refetchNewsfeed();
    };

    if (isLoadingNewsfeed) {
        return (
            <div className="flex justify-center items-center min-h-64">
                <LoadingSpinner size="lg" />
            </div>
        );
    }

    if (newsfeedError) {
        return (
            <div className="text-center py-12">
                <p className="text-gray-500 mb-4">Failed to load your newsfeed</p>
                <button
                    onClick={() => refetchNewsfeed()}
                    className="text-primary-600 hover:text-primary-700"
                >
                    Try again
                </button>
            </div>
        );
    }

    return (
        <div className="w-full">
            {/* Create Post Section */}
            <div className="mb-4">
                {showCreatePost ? (
                    <CreatePost
                        onPostCreated={handlePostCreated}
                        onCancel={() => setShowCreatePost(false)}
                    />
                ) : (
                    <div className="post-instagram">
                        <div className="p-4">
                            <button
                                onClick={() => setShowCreatePost(true)}
                                className="w-full flex items-center space-x-3 text-left"
                            >
                                <div className="story-border w-10 h-10 flex-shrink-0">
                                    {user?.profile_picture ? (
                                        <img
                                            src={user.profile_picture}
                                            alt={user.user_name}
                                            className="w-full h-full rounded-full object-cover"
                                        />
                                    ) : (
                                        <div className="w-full h-full bg-gray-300 rounded-full flex items-center justify-center">
                                            <span className="text-gray-600 text-xs font-medium">
                                                {user ? user.first_name[0] + user.last_name[0] : 'U'}
                                            </span>
                                        </div>
                                    )}
                                </div>
                                <div className="flex-1 py-3 px-4 bg-gray-50 rounded-full text-gray-500 hover:bg-gray-100 transition-colors">
                                    What's on your mind?
                                </div>
                            </button>
                        </div>
                    </div>
                )}
            </div>

            {/* Posts Feed */}
            <div className="space-y-4">
                {posts.length === 0 && !isLoadingNewsfeed ? (
                    <div className="post-instagram">
                        <div className="text-center py-12">
                            <div className="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
                                <PlusIcon className="w-8 h-8 text-gray-400" />
                            </div>
                            <h3 className="text-lg font-medium text-gray-900 mb-2">
                                No posts yet
                            </h3>
                            <p className="text-muted-instagram mb-6">
                                Start following people or create your first post to see content here.
                            </p>
                            <button
                                onClick={() => setShowCreatePost(true)}
                                className="btn-instagram"
                            >
                                <PlusIcon className="w-5 h-5 mr-2" />
                                Create Your First Post
                            </button>
                        </div>
                    </div>
                ) : (
                    posts.map((post) => (
                        <PostCard
                            key={post.post_id}
                            post={post}
                            onLike={handlePostLike}
                            onComment={handlePostComment}
                            onEdit={handlePostEdit}
                            onDelete={handlePostDelete}
                        />
                    ))
                )}
            </div>

            {/* Load more posts */}
            {posts.length > 0 && (
                <div className="text-center py-8">
                    <button className="text-link-instagram">
                        Load more posts
                    </button>
                </div>
            )}
        </div>
    );
};

export default Home; 
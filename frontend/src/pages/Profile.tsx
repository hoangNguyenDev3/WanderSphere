import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from 'react-query';
import { toast } from 'react-hot-toast';
import { userAPI, socialAPI, postsAPI } from '../services/api';
import { useAuth } from '../contexts/AuthContext';
import { User, Post, PostWithUser } from '../types/api';
import PostCard from '../components/Post/PostCard';
import LoadingSpinner from '../components/UI/LoadingSpinner';
import Button from '../components/UI/Button';
import {
    PencilIcon,
    UserPlusIcon,
    UserMinusIcon,
    CalendarDaysIcon,
    EnvelopeIcon,
} from '@heroicons/react/24/outline';

const Profile: React.FC = () => {
    const { userId } = useParams<{ userId: string }>();
    const navigate = useNavigate();
    const { user: currentUser } = useAuth();
    const queryClient = useQueryClient();

    const [posts, setPosts] = useState<PostWithUser[]>([]);
    const [isFollowing, setIsFollowing] = useState(false);

    const profileUserId = userId ? parseInt(userId) : currentUser?.user_id;
    const isOwnProfile = profileUserId === currentUser?.user_id;

    // Fetch user profile
    const {
        data: profileData,
        isLoading: isLoadingProfile,
        error: profileError,
    } = useQuery(
        ['profile', profileUserId],
        () => profileUserId ? userAPI.getProfile(profileUserId) : null,
        {
            enabled: !!profileUserId,
        }
    );

    // Fetch user posts
    const {
        data: userPostsData,
        isLoading: isLoadingPosts,
    } = useQuery(
        ['userPosts', profileUserId],
        () => profileUserId ? socialAPI.getUserPosts(profileUserId) : null,
        {
            enabled: !!profileUserId,
        }
    );

    // Fetch followers
    const {
        data: followersData,
    } = useQuery(
        ['followers', profileUserId],
        () => profileUserId ? socialAPI.getFollowers(profileUserId) : null,
        {
            enabled: !!profileUserId,
        }
    );

    // Fetch following
    const {
        data: followingData,
    } = useQuery(
        ['following', profileUserId],
        () => profileUserId ? socialAPI.getFollowing(profileUserId) : null,
        {
            enabled: !!profileUserId,
        }
    );

    // Check if current user is following this profile
    useEffect(() => {
        if (currentUser && followersData?.data?.followers_ids) {
            setIsFollowing(followersData.data.followers_ids.includes(currentUser.user_id));
        }
    }, [currentUser, followersData]);

    // Fetch individual posts
    useEffect(() => {
        const fetchPosts = async () => {
            if (!userPostsData?.data?.posts_ids) {
                setPosts([]);
                return;
            }

            try {
                const postPromises = userPostsData.data.posts_ids.map(async (postId) => {
                    const response = await postsAPI.getPost(postId);
                    return response.data;
                });

                const fetchedPosts = await Promise.all(postPromises);
                const validPosts = fetchedPosts.filter(Boolean) as Post[];

                const postsWithUser = validPosts.map(post => ({
                    ...post,
                    user: profileData?.data,
                })) as PostWithUser[];

                setPosts(postsWithUser);
            } catch (error) {
                console.error('Error fetching posts:', error);
            }
        };

        fetchPosts();
    }, [userPostsData, profileData]);

    // Follow/Unfollow mutation
    const followMutation = useMutation(
        (action: 'follow' | 'unfollow') => {
            return action === 'follow'
                ? socialAPI.followUser(profileUserId!)
                : socialAPI.unfollowUser(profileUserId!);
        },
        {
            onSuccess: (_, action) => {
                setIsFollowing(action === 'follow');
                queryClient.invalidateQueries(['followers', profileUserId]);
                toast.success(action === 'follow' ? 'User followed!' : 'User unfollowed!');
            },
            onError: () => {
                toast.error('Failed to update follow status');
            },
        }
    );

    const handleFollowToggle = () => {
        if (!profileUserId) return;
        followMutation.mutate(isFollowing ? 'unfollow' : 'follow');
    };

    const handlePostEdit = (postId: number) => {
        // Navigate to post edit or open modal
        console.log('Edit post:', postId);
    };

    const handlePostDelete = (postId: number) => {
        setPosts(prevPosts => prevPosts.filter(post => post.post_id !== postId));
    };

    const handlePostLike = (postId: number) => {
        // Optimistic update handled in PostCard
    };

    const handlePostComment = (postId: number) => {
        navigate(`/post/${postId}`);
    };

    if (isLoadingProfile) {
        return (
            <div className="flex justify-center items-center min-h-64">
                <LoadingSpinner size="lg" />
            </div>
        );
    }

    if (profileError || !profileData?.data) {
        return (
            <div className="text-center py-12">
                <p className="text-gray-500 mb-4">Profile not found</p>
                <button className="btn-instagram" onClick={() => navigate(-1)}>Go Back</button>
            </div>
        );
    }

    const user = profileData.data;
    const followersCount = followersData?.data?.followers_ids?.length || 0;
    const followingCount = followingData?.data?.followings_ids?.length || 0;
    const postsCount = posts.length;

    return (
        <div className="max-w-4xl mx-auto">
            {/* Profile Header */}
            <div className="bg-white rounded-lg border border-gray-200 shadow-sm mb-6 overflow-hidden">
                {/* Cover Photo */}
                <div className="h-32 sm:h-48 bg-gradient-to-br from-purple-400 via-pink-400 to-orange-400 relative">
                    {user.cover_picture && (
                        <img
                            src={user.cover_picture}
                            alt="Cover"
                            className="w-full h-full object-cover"
                        />
                    )}
                </div>

                {/* Profile Info */}
                <div className="px-6 pb-6 relative">
                    <div className="flex flex-col sm:flex-row sm:items-start sm:space-x-6 -mt-16">
                        {/* Profile Picture */}
                        <div className="flex-shrink-0 relative z-10 mb-4 sm:mb-0">
                            <div className="story-border w-24 h-24 sm:w-32 sm:h-32">
                                {user.profile_picture ? (
                                    <img
                                        src={user.profile_picture}
                                        alt={user.user_name}
                                        className="w-full h-full rounded-full object-cover"
                                    />
                                ) : (
                                    <div className="w-full h-full rounded-full bg-gray-300 flex items-center justify-center">
                                        <span className="text-gray-600 text-2xl font-bold">
                                            {user.first_name[0]}{user.last_name[0]}
                                        </span>
                                    </div>
                                )}
                            </div>
                        </div>

                        {/* User Info */}
                        <div className="flex-1 min-w-0 pt-16">
                            <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between">
                                <div className="flex-1">
                                    <h1 className="text-xl sm:text-2xl font-semibold text-gray-900">
                                        @{user.user_name}
                                    </h1>
                                    <p className="text-gray-600 mt-1">
                                        {user.first_name} {user.last_name}
                                    </p>
                                </div>

                                {/* Action Buttons */}
                                <div className="mt-4 sm:mt-0 flex space-x-3">
                                    {isOwnProfile ? (
                                        <button
                                            className="btn-secondary-instagram"
                                            onClick={() => navigate('/profile/edit')}
                                        >
                                            <PencilIcon className="w-4 h-4 mr-2" />
                                            Edit Profile
                                        </button>
                                    ) : (
                                        <button
                                            onClick={handleFollowToggle}
                                            disabled={followMutation.isLoading}
                                            className={isFollowing ? "btn-secondary-instagram" : "btn-instagram"}
                                        >
                                            {isFollowing ? (
                                                <>
                                                    <UserMinusIcon className="w-4 h-4 mr-2" />
                                                    Unfollow
                                                </>
                                            ) : (
                                                <>
                                                    <UserPlusIcon className="w-4 h-4 mr-2" />
                                                    Follow
                                                </>
                                            )}
                                        </button>
                                    )}
                                </div>
                            </div>

                            {/* Stats */}
                            <div className="mt-4 flex space-x-8">
                                <div className="text-center">
                                    <div className="text-lg font-semibold text-gray-900">{postsCount}</div>
                                    <div className="text-muted-instagram">posts</div>
                                </div>
                                <div
                                    className="text-center cursor-pointer hover:opacity-70 transition-opacity"
                                    onClick={() => navigate(`/profile/${user.user_id}/followers`)}
                                >
                                    <div className="text-lg font-semibold text-gray-900">{followersCount}</div>
                                    <div className="text-muted-instagram">followers</div>
                                </div>
                                <div
                                    className="text-center cursor-pointer hover:opacity-70 transition-opacity"
                                    onClick={() => navigate(`/profile/${user.user_id}/following`)}
                                >
                                    <div className="text-lg font-semibold text-gray-900">{followingCount}</div>
                                    <div className="text-muted-instagram">following</div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            {/* Posts */}
            <div className="space-y-4">
                <h2 className="text-lg font-semibold text-gray-900 px-2">Posts</h2>

                {isLoadingPosts ? (
                    <div className="flex justify-center py-12">
                        <LoadingSpinner size="lg" />
                    </div>
                ) : posts.length === 0 ? (
                    <div className="text-center py-12">
                        <p className="text-gray-500">
                            {isOwnProfile ? "You haven't posted anything yet." : "No posts to show."}
                        </p>
                        {isOwnProfile && (
                            <button
                                className="btn-instagram mt-4"
                                onClick={() => navigate('/')}
                            >
                                Create Your First Post
                            </button>
                        )}
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
        </div>
    );
};

export default Profile; 
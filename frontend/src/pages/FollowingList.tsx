import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from 'react-query';
import { toast } from 'react-hot-toast';
import { socialAPI, userAPI } from '../services/api';
import { useAuth } from '../contexts/AuthContext';
import { User } from '../types/api';
import LoadingSpinner from '../components/UI/LoadingSpinner';
import Button from '../components/UI/Button';
import {
    ArrowLeftIcon,
    UserPlusIcon,
    UserMinusIcon,
} from '@heroicons/react/24/outline';

interface UserWithFollowStatus extends User {
    isFollowing?: boolean;
}

const FollowingList: React.FC = () => {
    const { userId } = useParams<{ userId: string }>();
    const navigate = useNavigate();
    const { user: currentUser } = useAuth();
    const queryClient = useQueryClient();

    const [following, setFollowing] = useState<UserWithFollowStatus[]>([]);
    const [currentUserFollowing, setCurrentUserFollowing] = useState<number[]>([]);

    const profileUserId = userId ? parseInt(userId) : null;
    const isOwnProfile = profileUserId === currentUser?.user_id;

    // Fetch following
    const {
        data: followingData,
        isLoading: isLoadingFollowing,
        error: followingError,
    } = useQuery(
        ['following', profileUserId],
        () => profileUserId ? socialAPI.getFollowing(profileUserId) : null,
        {
            enabled: !!profileUserId,
        }
    );

    // Fetch current user's following list (if viewing someone else's list)
    const {
        data: currentUserFollowingData,
    } = useQuery(
        ['following', currentUser?.user_id],
        () => currentUser && !isOwnProfile ? socialAPI.getFollowing(currentUser.user_id) : null,
        {
            enabled: !!currentUser && !isOwnProfile,
        }
    );

    // Fetch profile owner info
    const {
        data: profileOwnerData,
    } = useQuery(
        ['profile', profileUserId],
        () => profileUserId ? userAPI.getProfile(profileUserId) : null,
        {
            enabled: !!profileUserId,
        }
    );

    // Update current user following list
    useEffect(() => {
        if (isOwnProfile && followingData?.data?.followings_ids) {
            setCurrentUserFollowing(followingData.data.followings_ids);
        } else if (currentUserFollowingData?.data?.followings_ids) {
            setCurrentUserFollowing(currentUserFollowingData.data.followings_ids);
        }
    }, [followingData, currentUserFollowingData, isOwnProfile]);

    // Fetch following details when we have following IDs
    useEffect(() => {
        const fetchFollowingDetails = async () => {
            if (!followingData?.data?.followings_ids || followingData.data.followings_ids.length === 0) {
                setFollowing([]);
                return;
            }

            try {
                const followingPromises = followingData.data.followings_ids.map(async (followingId) => {
                    const response = await userAPI.getProfile(followingId);
                    return {
                        ...response.data!,
                        isFollowing: isOwnProfile ? true : currentUserFollowing.includes(followingId),
                    };
                });

                const followingWithDetails = await Promise.all(followingPromises);
                setFollowing(followingWithDetails.filter(Boolean));
            } catch (error) {
                console.error('Error fetching following details:', error);
                setFollowing([]);
            }
        };

        fetchFollowingDetails();
    }, [followingData, currentUserFollowing, isOwnProfile]);

    // Follow/Unfollow mutation
    const followMutation = useMutation(
        ({ userId, action }: { userId: number; action: 'follow' | 'unfollow' }) => {
            return action === 'follow'
                ? socialAPI.followUser(userId)
                : socialAPI.unfollowUser(userId);
        },
        {
            onSuccess: (_, { userId, action }) => {
                if (isOwnProfile) {
                    // If viewing own following list, remove unfollowed users
                    if (action === 'unfollow') {
                        setFollowing(prev => prev.filter(user => user.user_id !== userId));
                    }
                } else {
                    // Update follow status
                    setFollowing(prev => prev.map(user =>
                        user.user_id === userId
                            ? { ...user, isFollowing: action === 'follow' }
                            : user
                    ));
                }

                // Update current user following list
                setCurrentUserFollowing(prev =>
                    action === 'follow'
                        ? [...prev, userId]
                        : prev.filter(id => id !== userId)
                );

                queryClient.invalidateQueries(['following', currentUser?.user_id]);
                if (isOwnProfile) {
                    queryClient.invalidateQueries(['following', profileUserId]);
                }
                toast.success(action === 'follow' ? 'User followed!' : 'User unfollowed!');
            },
            onError: () => {
                toast.error('Failed to update follow status');
            },
        }
    );

    const handleFollowToggle = (userId: number, isCurrentlyFollowing: boolean) => {
        followMutation.mutate({
            userId,
            action: isCurrentlyFollowing ? 'unfollow' : 'follow',
        });
    };

    const handleUserClick = (userId: number) => {
        navigate(`/profile/${userId}`);
    };

    if (isLoadingFollowing) {
        return (
            <div className="flex justify-center items-center min-h-64">
                <LoadingSpinner size="lg" />
            </div>
        );
    }

    if (followingError || !profileUserId) {
        return (
            <div className="text-center py-12">
                <p className="text-gray-500 mb-4">Unable to load following list</p>
                <button className="btn-instagram" onClick={() => navigate(-1)}>Go Back</button>
            </div>
        );
    }

    return (
        <div className="max-w-2xl mx-auto">
            {/* Header */}
            <div className="flex items-center space-x-3 mb-6">
                <button
                    onClick={() => navigate(-1)}
                    className="text-gray-600 hover:text-gray-800 transition-colors"
                >
                    <ArrowLeftIcon className="w-6 h-6" />
                </button>
                <div>
                    <h1 className="text-xl font-semibold text-gray-900">Following</h1>
                    {profileOwnerData?.data && (
                        <p className="text-muted-instagram">
                            {profileOwnerData.data.user_name}
                        </p>
                    )}
                </div>
            </div>

            {/* Following List */}
            <div className="bg-white rounded-lg border border-gray-200 shadow-sm">
                {following.length === 0 ? (
                    <div className="text-center py-12">
                        <p className="text-gray-500">
                            {isOwnProfile ? "You're not following anyone yet" : "Not following anyone yet"}
                        </p>
                        {isOwnProfile && (
                            <p className="text-gray-400 text-sm mt-1">
                                Discover people to follow on the home feed!
                            </p>
                        )}
                    </div>
                ) : (
                    <div className="divide-y divide-gray-200">
                        {following.map((user) => (
                            <div key={user.user_id} className="p-4 flex items-center justify-between">
                                <div
                                    className="flex items-center space-x-3 cursor-pointer flex-1"
                                    onClick={() => handleUserClick(user.user_id)}
                                >
                                    <div className="story-border w-12 h-12">
                                        {user.profile_picture ? (
                                            <img
                                                src={user.profile_picture}
                                                alt={user.user_name}
                                                className="w-full h-full rounded-full object-cover"
                                            />
                                        ) : (
                                            <div className="w-full h-full rounded-full bg-gray-300 flex items-center justify-center">
                                                <span className="text-gray-600 font-medium text-sm">
                                                    {user.first_name[0]}{user.last_name[0]}
                                                </span>
                                            </div>
                                        )}
                                    </div>
                                    <div className="flex-1 min-w-0">
                                        <h3 className="text-username truncate">
                                            {user.user_name}
                                        </h3>
                                        <p className="text-muted-instagram truncate">
                                            {user.first_name} {user.last_name}
                                        </p>
                                    </div>
                                </div>

                                {/* Follow/Unfollow Button */}
                                {currentUser && user.user_id !== currentUser.user_id && (
                                    <button
                                        className={user.isFollowing ? "btn-secondary-instagram" : "btn-instagram"}
                                        onClick={() => handleFollowToggle(user.user_id, user.isFollowing || false)}
                                        disabled={followMutation.isLoading}
                                    >
                                        {user.isFollowing ? 'Following' : 'Follow'}
                                    </button>
                                )}
                            </div>
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
};

export default FollowingList; 
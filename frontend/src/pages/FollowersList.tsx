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

const FollowersList: React.FC = () => {
    const { userId } = useParams<{ userId: string }>();
    const navigate = useNavigate();
    const { user: currentUser } = useAuth();
    const queryClient = useQueryClient();

    const [followers, setFollowers] = useState<UserWithFollowStatus[]>([]);
    const [currentUserFollowing, setCurrentUserFollowing] = useState<number[]>([]);

    const profileUserId = userId ? parseInt(userId) : null;

    // Fetch followers
    const {
        data: followersData,
        isLoading: isLoadingFollowers,
        error: followersError,
    } = useQuery(
        ['followers', profileUserId],
        () => profileUserId ? socialAPI.getFollowers(profileUserId) : null,
        {
            enabled: !!profileUserId,
        }
    );

    // Fetch current user's following list
    const {
        data: currentUserFollowingData,
    } = useQuery(
        ['following', currentUser?.user_id],
        () => currentUser ? socialAPI.getFollowing(currentUser.user_id) : null,
        {
            enabled: !!currentUser,
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
        if (currentUserFollowingData?.data?.followings_ids) {
            setCurrentUserFollowing(currentUserFollowingData.data.followings_ids);
        }
    }, [currentUserFollowingData]);

    // Fetch followers details when we have follower IDs
    useEffect(() => {
        const fetchFollowersDetails = async () => {
            if (!followersData?.data?.followers_ids || followersData.data.followers_ids.length === 0) {
                setFollowers([]);
                return;
            }

            try {
                const followerPromises = followersData.data.followers_ids.map(async (followerId) => {
                    const response = await userAPI.getProfile(followerId);
                    return {
                        ...response.data!,
                        isFollowing: currentUserFollowing.includes(followerId),
                    };
                });

                const followersWithDetails = await Promise.all(followerPromises);
                setFollowers(followersWithDetails.filter(Boolean));
            } catch (error) {
                console.error('Error fetching followers details:', error);
                setFollowers([]);
            }
        };

        fetchFollowersDetails();
    }, [followersData, currentUserFollowing]);

    // Follow/Unfollow mutation
    const followMutation = useMutation(
        ({ userId, action }: { userId: number; action: 'follow' | 'unfollow' }) => {
            return action === 'follow'
                ? socialAPI.followUser(userId)
                : socialAPI.unfollowUser(userId);
        },
        {
            onSuccess: (_, { userId, action }) => {
                // Update local state
                setFollowers(prev => prev.map(follower =>
                    follower.user_id === userId
                        ? { ...follower, isFollowing: action === 'follow' }
                        : follower
                ));

                // Update current user following list
                setCurrentUserFollowing(prev =>
                    action === 'follow'
                        ? [...prev, userId]
                        : prev.filter(id => id !== userId)
                );

                queryClient.invalidateQueries(['following', currentUser?.user_id]);
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

    if (isLoadingFollowers) {
        return (
            <div className="flex justify-center items-center min-h-64">
                <LoadingSpinner size="lg" />
            </div>
        );
    }

    if (followersError || !profileUserId) {
        return (
            <div className="text-center py-12">
                <p className="text-gray-500 mb-4">Unable to load followers</p>
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
                <div>
                    <h1 className="text-xl font-semibold text-gray-900">Followers</h1>
                    {profileOwnerData?.data && (
                        <p className="text-gray-500">
                            {profileOwnerData.data.first_name} {profileOwnerData.data.last_name}
                        </p>
                    )}
                </div>
            </div>

            {/* Followers List */}
            <div className="card">
                {followers.length === 0 ? (
                    <div className="text-center py-12">
                        <p className="text-gray-500">No followers yet</p>
                    </div>
                ) : (
                    <div className="divide-y divide-gray-200">
                        {followers.map((follower) => (
                            <div key={follower.user_id} className="p-4 flex items-center justify-between">
                                <div
                                    className="flex items-center space-x-3 cursor-pointer flex-1"
                                    onClick={() => handleUserClick(follower.user_id)}
                                >
                                    {follower.profile_picture ? (
                                        <img
                                            src={follower.profile_picture}
                                            alt={follower.user_name}
                                            className="w-12 h-12 rounded-full object-cover"
                                        />
                                    ) : (
                                        <div className="w-12 h-12 rounded-full bg-gray-100 flex items-center justify-center">
                                            <span className="text-gray-600 font-medium">
                                                {follower.first_name[0]}{follower.last_name[0]}
                                            </span>
                                        </div>
                                    )}
                                    <div className="flex-1 min-w-0">
                                        <h3 className="font-medium text-gray-900 truncate">
                                            {follower.first_name} {follower.last_name}
                                        </h3>
                                        <p className="text-gray-500 text-sm truncate">
                                            @{follower.user_name}
                                        </p>
                                    </div>
                                </div>

                                {/* Follow/Unfollow Button */}
                                {currentUser && follower.user_id !== currentUser.user_id && (
                                    <Button
                                        size="sm"
                                        variant={follower.isFollowing ? "secondary" : "primary"}
                                        onClick={() => handleFollowToggle(follower.user_id, follower.isFollowing || false)}
                                        disabled={followMutation.isLoading}
                                    >
                                        {follower.isFollowing ? (
                                            <>
                                                <UserMinusIcon className="w-4 h-4 mr-1" />
                                                Unfollow
                                            </>
                                        ) : (
                                            <>
                                                <UserPlusIcon className="w-4 h-4 mr-1" />
                                                Follow
                                            </>
                                        )}
                                    </Button>
                                )}
                            </div>
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
};

export default FollowersList; 
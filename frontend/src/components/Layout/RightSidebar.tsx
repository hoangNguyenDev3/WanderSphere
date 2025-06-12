import React, { useState, useEffect, useCallback } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';
import { socialAPI, userAPI } from '../../services/api';
import { User } from '../../types/api';

const RightSidebar: React.FC = () => {
    const { user } = useAuth();
    const [suggestedUsers, setSuggestedUsers] = useState<User[]>([]);
    const [loading, setLoading] = useState(true);

    const fetchSuggestedUsers = useCallback(async () => {
        if (!user?.user_id) return;

        try {
            setLoading(true);

            // First, get the current user's following list
            const followingResponse = await socialAPI.getFollowing(user.user_id);
            const followingIds = followingResponse.data?.followings_ids || [];

            // Try a larger range of user IDs to find users
            const maxUserId = 20; // Adjust this based on your user base
            const potentialUserIds = Array.from({ length: maxUserId }, (_, i) => i + 1);

            // Randomize the order to get different suggestions each time
            const shuffledIds = potentialUserIds.sort(() => Math.random() - 0.5);

            const userPromises = shuffledIds.map(async (id) => {
                try {
                    const response = await userAPI.getProfile(id);
                    return response.data;
                } catch (error) {
                    return null;
                }
            });

            const users = await Promise.all(userPromises);
            const validUsers = users.filter((suggestedUser): suggestedUser is User =>
                suggestedUser !== null &&
                suggestedUser !== undefined &&
                suggestedUser.user_id !== user?.user_id && // Not the current user
                !followingIds.includes(suggestedUser.user_id) // Not already following
            );

            // Take only the first 5 suggestions
            setSuggestedUsers(validUsers.slice(0, 5));
        } catch (error) {
            console.error('Error fetching suggested users:', error);
        } finally {
            setLoading(false);
        }
    }, [user?.user_id]);

    useEffect(() => {
        fetchSuggestedUsers();
    }, [fetchSuggestedUsers]);

    const handleFollow = async (userId: number) => {
        try {
            await socialAPI.followUser(userId);
            // Remove from suggested users after following
            setSuggestedUsers(prev => prev.filter(user => user.user_id !== userId));

            // If we have less than 3 suggestions left, fetch more
            if (suggestedUsers.length <= 3) {
                // Small delay to ensure backend is updated
                setTimeout(() => {
                    fetchSuggestedUsers();
                }, 500);
            }
        } catch (error) {
            console.error('Error following user:', error);
        }
    };

    if (!user) return null;

    return (
        <div className="w-80 p-6 space-y-6">
            {/* Suggested for You */}
            <div className="bg-white rounded-lg p-4 shadow-sm border border-gray-200">
                <div className="flex items-center justify-between mb-4">
                    <h3 className="text-muted-instagram font-medium">Suggested for you</h3>
                    <Link
                        to="/search"
                        className="text-link-instagram text-sm hover:text-gray-600"
                    >
                        See All
                    </Link>
                </div>

                {loading ? (
                    <div className="space-y-3">
                        {[1, 2, 3].map((i) => (
                            <div key={i} className="flex items-center space-x-3">
                                <div className="w-10 h-10 bg-gray-200 rounded-full animate-pulse"></div>
                                <div className="flex-1">
                                    <div className="h-4 bg-gray-200 rounded animate-pulse mb-1"></div>
                                    <div className="h-3 bg-gray-200 rounded animate-pulse w-2/3"></div>
                                </div>
                                <div className="w-16 h-8 bg-gray-200 rounded animate-pulse"></div>
                            </div>
                        ))}
                    </div>
                ) : (
                    <div className="space-y-3">
                        {loading ? (
                            <div className="space-y-3">
                                {[1, 2, 3].map((i) => (
                                    <div key={i} className="flex items-center space-x-3 animate-pulse">
                                        <div className="w-10 h-10 bg-gray-200 rounded-full"></div>
                                        <div className="flex-1">
                                            <div className="h-3 bg-gray-200 rounded w-20 mb-1"></div>
                                            <div className="h-2 bg-gray-200 rounded w-16"></div>
                                        </div>
                                        <div className="h-6 bg-gray-200 rounded w-12"></div>
                                    </div>
                                ))}
                            </div>
                        ) : suggestedUsers.length === 0 ? (
                            <p className="text-muted-instagram text-sm">No new users to suggest</p>
                        ) : (
                            suggestedUsers.map((suggestedUser) => (
                                <div key={suggestedUser.user_id} className="flex items-center space-x-3">
                                    <Link to={`/profile/${suggestedUser.user_id}`}>
                                        {suggestedUser.profile_picture ? (
                                            <img
                                                src={suggestedUser.profile_picture}
                                                alt={suggestedUser.user_name}
                                                className="w-10 h-10 rounded-full object-cover"
                                            />
                                        ) : (
                                            <div className="w-10 h-10 bg-gray-300 rounded-full flex items-center justify-center">
                                                <span className="text-gray-600 text-xs font-medium">
                                                    {suggestedUser.first_name[0] + suggestedUser.last_name[0]}
                                                </span>
                                            </div>
                                        )}
                                    </Link>
                                    <div className="flex-1 min-w-0">
                                        <Link
                                            to={`/profile/${suggestedUser.user_id}`}
                                            className="text-username hover:text-gray-600 transition-colors text-sm block truncate"
                                        >
                                            {suggestedUser.user_name}
                                        </Link>
                                        <p className="text-muted-instagram text-xs truncate">
                                            {suggestedUser.first_name} {suggestedUser.last_name}
                                        </p>
                                    </div>
                                    <button
                                        onClick={() => handleFollow(suggestedUser.user_id)}
                                        className="text-link-instagram text-sm font-medium hover:text-gray-600 transition-colors"
                                    >
                                        Follow
                                    </button>
                                </div>
                            ))
                        )}
                    </div>
                )}
            </div>
            {/* Footer */}
            <div className="text-muted-instagram text-xs leading-relaxed">
                <div className="space-x-2">
                    <button className="hover:underline">About</button>
                    <span>·</span>
                    <button className="hover:underline">Help</button>
                    <span>·</span>
                    <button className="hover:underline">Press</button>
                    <span>·</span>
                    <button className="hover:underline">API</button>
                    <span>·</span>
                    <button className="hover:underline">Jobs</button>
                    <span>·</span>
                    <button className="hover:underline">Privacy</button>
                    <span>·</span>
                    <button className="hover:underline">Terms</button>
                </div>
                <div className="mt-3">
                    © 2025 WanderSphere
                </div>
            </div>
        </div>
    );
};

export default RightSidebar; 
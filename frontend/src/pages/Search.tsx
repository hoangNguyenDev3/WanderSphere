import React, { useState, useEffect } from 'react';
import { useQuery, useMutation, useQueryClient } from 'react-query';
import { useNavigate } from 'react-router-dom';
import { toast } from 'react-hot-toast';
import { userAPI, socialAPI } from '../services/api';
import { useAuth } from '../contexts/AuthContext';
import { User } from '../types/api';
import LoadingSpinner from '../components/UI/LoadingSpinner';
import Button from '../components/UI/Button';
import {
    MagnifyingGlassIcon,
    UserPlusIcon,
    UserMinusIcon,
} from '@heroicons/react/24/outline';

interface UserWithFollowStatus extends User {
    isFollowing?: boolean;
}

const Search: React.FC = () => {
    const navigate = useNavigate();
    const { user: currentUser } = useAuth();
    const queryClient = useQueryClient();

    const [searchQuery, setSearchQuery] = useState('');
    const [searchResults, setSearchResults] = useState<UserWithFollowStatus[]>([]);
    const [isSearching, setIsSearching] = useState(false);
    const [currentUserFollowing, setCurrentUserFollowing] = useState<number[]>([]);

    // Fetch current user's following list
    const {
        data: followingData,
    } = useQuery(
        ['following', currentUser?.user_id],
        () => currentUser ? socialAPI.getFollowing(currentUser.user_id) : null,
        {
            enabled: !!currentUser,
        }
    );

    // Update current user following list
    useEffect(() => {
        if (followingData?.data?.followings_ids) {
            setCurrentUserFollowing(followingData.data.followings_ids);
        }
    }, [followingData]);

    // Simulate search functionality (in real app, this would be a backend endpoint)
    const performSearch = async (query: string) => {
        if (!query.trim()) {
            setSearchResults([]);
            return;
        }

        setIsSearching(true);
        try {
            // Since there's no search endpoint, we'll simulate by filtering users
            // In a real app, you'd have a search API endpoint

            // For demo purposes, we'll create some sample users
            // In reality, this would be an API call to /api/v1/search/users?q=query
            const sampleUsers: User[] = [
                {
                    user_id: 1,
                    user_name: 'john_doe',
                    first_name: 'John',
                    last_name: 'Doe',
                    email: 'john@example.com',
                    date_of_birth: '1990-01-01',
                    profile_picture: '',
                    cover_picture: '',
                },
                {
                    user_id: 2,
                    user_name: 'jane_smith',
                    first_name: 'Jane',
                    last_name: 'Smith',
                    email: 'jane@example.com',
                    date_of_birth: '1992-05-15',
                    profile_picture: '',
                    cover_picture: '',
                },
                {
                    user_id: 3,
                    user_name: 'mike_wilson',
                    first_name: 'Mike',
                    last_name: 'Wilson',
                    email: 'mike@example.com',
                    date_of_birth: '1988-11-22',
                    profile_picture: '',
                    cover_picture: '',
                },
            ];

            // Filter users based on search query
            const filteredUsers = sampleUsers.filter(user =>
                user.user_name.toLowerCase().includes(query.toLowerCase()) ||
                user.first_name.toLowerCase().includes(query.toLowerCase()) ||
                user.last_name.toLowerCase().includes(query.toLowerCase()) ||
                user.email.toLowerCase().includes(query.toLowerCase())
            );

            // Add follow status
            const usersWithFollowStatus = filteredUsers.map(user => ({
                ...user,
                isFollowing: currentUserFollowing.includes(user.user_id),
            }));

            setSearchResults(usersWithFollowStatus);
        } catch (error) {
            console.error('Search error:', error);
            toast.error('Failed to search users');
        } finally {
            setIsSearching(false);
        }
    };

    // Debounced search
    useEffect(() => {
        const timer = setTimeout(() => {
            performSearch(searchQuery);
        }, 300);

        return () => clearTimeout(timer);
    }, [searchQuery, currentUserFollowing]);

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
                setSearchResults(prev => prev.map(user =>
                    user.user_id === userId
                        ? { ...user, isFollowing: action === 'follow' }
                        : user
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

    return (
        <div className="max-w-2xl mx-auto">
            {/* Header */}
            <div className="mb-6">
                <h1 className="text-2xl font-bold text-gray-900 mb-4">Search</h1>

                {/* Search Input */}
                <div className="relative">
                    <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                        <MagnifyingGlassIcon className="h-5 w-5 text-gray-400" />
                    </div>
                    <input
                        type="text"
                        placeholder="Search for users..."
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        className="block w-full pl-10 pr-3 py-3 border border-gray-300 rounded-lg leading-5 bg-white placeholder-gray-500 focus:outline-none focus:placeholder-gray-400 focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                    />
                    {isSearching && (
                        <div className="absolute inset-y-0 right-0 pr-3 flex items-center">
                            <LoadingSpinner size="sm" />
                        </div>
                    )}
                </div>
            </div>

            {/* Search Results */}
            <div className="space-y-4">
                {!searchQuery.trim() ? (
                    <div className="card p-8 text-center">
                        <MagnifyingGlassIcon className="w-12 h-12 text-gray-400 mx-auto mb-4" />
                        <h3 className="text-lg font-medium text-gray-900 mb-2">
                            Discover People
                        </h3>
                        <p className="text-gray-500">
                            Search for people by name, username, or email to connect with them.
                        </p>
                    </div>
                ) : searchResults.length === 0 && !isSearching ? (
                    <div className="card p-8 text-center">
                        <p className="text-gray-500">
                            No users found for "{searchQuery}"
                        </p>
                        <p className="text-gray-400 text-sm mt-1">
                            Try searching with a different term.
                        </p>
                    </div>
                ) : (
                    <div className="card">
                        <div className="px-4 py-3 border-b border-gray-200">
                            <h2 className="text-lg font-semibold text-gray-900">
                                Search Results ({searchResults.length})
                            </h2>
                        </div>
                        <div className="divide-y divide-gray-200">
                            {searchResults.map((user) => (
                                <div key={user.user_id} className="p-4 flex items-center justify-between hover:bg-gray-50">
                                    <div
                                        className="flex items-center space-x-3 cursor-pointer flex-1"
                                        onClick={() => handleUserClick(user.user_id)}
                                    >
                                        {user.profile_picture ? (
                                            <img
                                                src={user.profile_picture}
                                                alt={user.user_name}
                                                className="w-12 h-12 rounded-full object-cover"
                                            />
                                        ) : (
                                            <div className="w-12 h-12 rounded-full bg-primary-100 flex items-center justify-center">
                                                <span className="text-primary-600 font-medium">
                                                    {user.first_name[0]}{user.last_name[0]}
                                                </span>
                                            </div>
                                        )}
                                        <div className="flex-1 min-w-0">
                                            <h3 className="font-medium text-gray-900 truncate">
                                                {user.first_name} {user.last_name}
                                            </h3>
                                            <p className="text-gray-500 text-sm truncate">
                                                @{user.user_name}
                                            </p>
                                        </div>
                                    </div>

                                    {/* Follow/Unfollow Button */}
                                    {currentUser && user.user_id !== currentUser.user_id && (
                                        <Button
                                            size="sm"
                                            variant={user.isFollowing ? "secondary" : "primary"}
                                            onClick={() => handleFollowToggle(user.user_id, user.isFollowing || false)}
                                            disabled={followMutation.isLoading}
                                        >
                                            {user.isFollowing ? (
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
                    </div>
                )}
            </div>

            {/* Suggested Users (when no search query) */}
            {!searchQuery.trim() && (
                <div className="mt-8">
                    <div className="card p-6">
                        <h3 className="text-lg font-semibold text-gray-900 mb-4">
                            Suggestions for You
                        </h3>
                        <div className="text-center py-8">
                            <p className="text-gray-500">
                                User suggestions would appear here based on mutual connections and interests.
                            </p>
                            <p className="text-gray-400 text-sm mt-1">
                                Start following people to get personalized suggestions!
                            </p>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default Search; 
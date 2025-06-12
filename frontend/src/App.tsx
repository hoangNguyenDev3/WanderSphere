import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from 'react-query';
import { Toaster } from 'react-hot-toast';
import { AuthProvider, useAuth } from './contexts/AuthContext';
import { ThemeProvider } from './contexts/ThemeContext';
import Login from './pages/Login';
import Signup from './pages/Signup';
import Home from './pages/Home';
import Profile from './pages/Profile';
import EditProfile from './pages/EditProfile';
import PostDetail from './pages/PostDetail';
import FollowersList from './pages/FollowersList';
import FollowingList from './pages/FollowingList';
import Search from './pages/Search';
import CreatePost from './components/Post/CreatePost';
import Layout from './components/Layout/Layout';
import LoadingSpinner from './components/UI/LoadingSpinner';

const queryClient = new QueryClient({
    defaultOptions: {
        queries: {
            retry: 1,
            refetchOnWindowFocus: false,
        },
    },
});

// Protected Route Component
const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const { isAuthenticated, isLoading } = useAuth();

    if (isLoading) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-background">
                <LoadingSpinner />
            </div>
        );
    }

    return isAuthenticated ? <>{children}</> : <Navigate to="/login" replace />;
};

// Public Route Component (redirects to home if authenticated)
const PublicRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const { isAuthenticated, isLoading } = useAuth();

    if (isLoading) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-background">
                <LoadingSpinner />
            </div>
        );
    }

    return !isAuthenticated ? <>{children}</> : <Navigate to="/" replace />;
};

// Create Post Page wrapper
const CreatePostPage: React.FC = () => {
    return (
        <div className="max-w-2xl mx-auto">
            <h1 className="text-2xl font-bold text-foreground mb-6">Create New Post</h1>
            <CreatePost
                onPostCreated={() => {
                    // Handle post creation success
                    window.location.href = '/';
                }}
                onCancel={() => {
                    // Handle cancel
                    window.history.back();
                }}
            />
        </div>
    );
};

const App: React.FC = () => {
    return (
        <QueryClientProvider client={queryClient}>
            <ThemeProvider>
                <AuthProvider>
                    <Router>
                        <div className="min-h-screen bg-background text-foreground">
                            <Routes>
                                <Route
                                    path="/login"
                                    element={
                                        <PublicRoute>
                                            <Login />
                                        </PublicRoute>
                                    }
                                />
                                <Route
                                    path="/signup"
                                    element={
                                        <PublicRoute>
                                            <Signup />
                                        </PublicRoute>
                                    }
                                />
                                <Route
                                    path="/*"
                                    element={
                                        <ProtectedRoute>
                                            <Layout>
                                                <div className="pb-20 lg:pb-0">
                                                    <Routes>
                                                        <Route path="/" element={<Home />} />
                                                        <Route path="/create" element={<CreatePostPage />} />
                                                        <Route path="/search" element={<Search />} />
                                                        <Route path="/profile/edit" element={<EditProfile />} />
                                                        <Route path="/profile/:userId/followers" element={<FollowersList />} />
                                                        <Route path="/profile/:userId/following" element={<FollowingList />} />
                                                        <Route path="/profile/:userId" element={<Profile />} />
                                                        <Route path="/post/:postId" element={<PostDetail />} />
                                                        <Route path="*" element={<Navigate to="/" replace />} />
                                                    </Routes>
                                                </div>
                                            </Layout>
                                        </ProtectedRoute>
                                    }
                                />
                            </Routes>
                        </div>
                        <Toaster
                            position="top-center"
                            toastOptions={{
                                duration: 4000,
                                style: {
                                    background: 'hsl(var(--card))',
                                    color: 'hsl(var(--card-foreground))',
                                    border: '1px solid hsl(var(--border))',
                                    borderRadius: '12px',
                                    fontSize: '14px',
                                    fontWeight: '500',
                                },
                                success: {
                                    iconTheme: {
                                        primary: 'hsl(var(--primary))',
                                        secondary: 'hsl(var(--primary-foreground))',
                                    },
                                },
                                error: {
                                    iconTheme: {
                                        primary: 'hsl(var(--destructive))',
                                        secondary: 'hsl(var(--destructive-foreground))',
                                    },
                                },
                            }}
                        />
                    </Router>
                </AuthProvider>
            </ThemeProvider>
        </QueryClientProvider>
    );
};

export default App; 
import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { User, LoginRequest, CreateUserRequest } from '../types/api';
import { authAPI } from '../services/api';
import toast from 'react-hot-toast';

interface AuthContextType {
    user: User | null;
    isLoading: boolean;
    isAuthenticated: boolean;
    login: (data: LoginRequest) => Promise<boolean>;
    signup: (data: CreateUserRequest) => Promise<boolean>;
    logout: () => void;
    updateUser: (user: User) => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (!context) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
};

interface AuthProviderProps {
    children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
    const [user, setUser] = useState<User | null>(null);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        // Check if user is stored in localStorage on app start
        const storedUser = localStorage.getItem('user');
        if (storedUser) {
            try {
                setUser(JSON.parse(storedUser));
            } catch (error) {
                localStorage.removeItem('user');
            }
        }
        setIsLoading(false);
    }, []);

    const login = async (data: LoginRequest): Promise<boolean> => {
        try {
            setIsLoading(true);
            const response = await authAPI.login(data);

            if (response.data && response.status >= 200 && response.status < 300) {
                const userData = response.data.user;
                setUser(userData);
                localStorage.setItem('user', JSON.stringify(userData));
                toast.success('Welcome back!');
                return true;
            } else {
                const errorMessage = response.error?.message || 'Login failed';
                toast.error(errorMessage);
                return false;
            }
        } catch (error) {
            toast.error('Login failed. Please try again.');
            return false;
        } finally {
            setIsLoading(false);
        }
    };

    const signup = async (data: CreateUserRequest): Promise<boolean> => {
        try {
            setIsLoading(true);
            const response = await authAPI.signup(data);

            if (response.data && response.status >= 200 && response.status < 300) {
                toast.success('Account created successfully! Please log in.');
                return true;
            } else {
                const errorMessage = response.error?.message || 'Signup failed';
                toast.error(errorMessage);
                return false;
            }
        } catch (error) {
            toast.error('Signup failed. Please try again.');
            return false;
        } finally {
            setIsLoading(false);
        }
    };

    const logout = () => {
        setUser(null);
        localStorage.removeItem('user');
        toast.success('Logged out successfully');

        // Call logout API to clear server-side session
        authAPI.logout().catch(console.error);
    };

    const updateUser = (updatedUser: User) => {
        setUser(updatedUser);
        localStorage.setItem('user', JSON.stringify(updatedUser));
    };

    const value: AuthContextType = {
        user,
        isLoading,
        isAuthenticated: !!user,
        login,
        signup,
        logout,
        updateUser,
    };

    return (
        <AuthContext.Provider value={value}>
            {children}
        </AuthContext.Provider>
    );
}; 
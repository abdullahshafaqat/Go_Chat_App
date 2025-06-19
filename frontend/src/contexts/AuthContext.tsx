
import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { authApi, setAuthToken, clearAuth, getAuthToken } from '../api';
import { toast } from '@/hooks/use-toast';

interface User {
  id: string;
  username: string;
  email: string;
}

interface AuthContextType {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  login: (email: string, password: string) => Promise<void>;
  signup: (username: string, email: string, password: string) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
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
    const token = getAuthToken();
    console.log('Checking for existing token:', token);
    if (token) {
      // You might want to validate the token with the backend here
      // For now, we'll assume the token is valid if it exists
      setUser({ id: '1', username: 'User', email: 'user@example.com' });
    }
    setIsLoading(false);
  }, []);

  const login = async (email: string, password: string) => {
    try {
      setIsLoading(true);
      console.log('Attempting login with:', { email });
      const response = await authApi.login({ email, password });
      console.log('Login response:', response);
      
      const userData = response.user || { 
        id: response.user_id || '1', 
        username: response.username || 'User', 
        email: email 
      };
      
      setUser(userData);
      toast({
        title: "Login Successful",
        description: "Welcome back!",
      });
    } catch (error: any) {
      console.error('Login error:', error);
      toast({
        title: "Login Failed",
        description: error.response?.data?.message || "Invalid credentials",
        variant: "destructive",
      });
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  const signup = async (username: string, email: string, password: string) => {
    try {
      setIsLoading(true);
      console.log('Attempting signup with:', { username, email });
      await authApi.signup({ username, email, password });
      // Auto-login after successful signup
      await login(email, password);
      toast({
        title: "Account Created",
        description: "Welcome! You've been automatically logged in.",
      });
    } catch (error: any) {
      console.error('Signup error:', error);
      toast({
        title: "Signup Failed",
        description: error.response?.data?.message || "Failed to create account",
        variant: "destructive",
      });
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  const logout = () => {
    console.log('Logging out user');
    clearAuth();
    setUser(null);
    toast({
      title: "Logged Out",
      description: "You've been successfully logged out.",
    });
  };

  const value: AuthContextType = {
    user,
    isLoading,
    isAuthenticated: !!user,
    login,
    signup,
    logout,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

// src/contexts/AuthContext.tsx
import React, { createContext, useContext, useState, useEffect, ReactNode } from "react";
import { authApi, setAuthToken, clearAuth, getAuthToken } from "../services/api";
import { toast } from "@/hooks/use-toast";
import { webSocketService } from "../services/websocketservice";

interface User {
  id: number;
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
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};

export const AuthProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const token = getAuthToken();
    if (token) {
      // Get profile data or decode token to retrieve user info
      setUser({ id: 1, username: "User", email: "user@example.com" }); // Load real profile data here
    }
    setIsLoading(false);
  }, []);

  useEffect(() => {
    if (user) {
      const token = getAuthToken();
      if (token) {
        webSocketService.connect(token, (message) => {
          console.log("Incoming message:", message);
        });
      }
      return () => webSocketService.disconnect();
    }
  }, [user]);

  const login = async (email: string, password: string) => {
    try {
      setIsLoading(true);
      const response = await authApi.login({ email, password });
      const id = response.user_id ? parseInt(response.user_id, 10) : 1;

      const userData: User = {
        id,
        username: response.username || "User",
        email,
      };
      setUser(userData);

      toast({ title: "Login Successful", description: "Welcome back!" });
    } catch (error: any) {
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
      await authApi.signup({ username, email, password });
      await login(email, password); // auto-login after signup
      toast({ title: "Account Created", description: "Welcome! You've been logged in." });
    } catch (error: any) {
      toast({ title: "Signup Failed", description: error.response?.data?.message || "Error creating account", variant: "destructive" });
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  const logout = () => {
    webSocketService.disconnect();
    clearAuth();
    setUser(null);
    toast({ title: "Logged Out", description: "You have successfully logged out." });
  };

  return (
    <AuthContext.Provider value={{ user, isLoading, isAuthenticated: !!user, login, signup, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

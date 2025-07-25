// src/contexts/AuthContext.tsx
import React, { createContext, useContext, useState, useEffect, ReactNode } from "react";
import { authApi, setAuthToken, clearAuth, getAuthToken } from "../services/api";
import { toast } from "@/hooks/use-toast";

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
      try {
        const payloadBase64 = token.split(".")[1];
        const payload = JSON.parse(atob(payloadBase64));
        setUser({
          id: Number(payload.ID),
          username: payload.username || "",
          email: payload.email || "",
        });
      } catch (error) {
        console.error("[AuthProvider] Error decoding token:", error);
        clearAuth();
      }
    }
    setIsLoading(false);
  }, []);

  const login = async (email: string, password: string) => {
    try {
      setIsLoading(true);
      const response = await authApi.login({ email, password });

      // Save token & refresh token
      setAuthToken(response.access_token);  // also sets Axios headers
      localStorage.setItem('token', response.access_token);
      localStorage.setItem('refreshToken', response.refresh_token);

      const id = response.user_id ? parseInt(response.user_id, 10) : 0;
      setUser({ id, username: response.username || "User", email });

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
      toast({
        title: "Signup Failed",
        description: error.response?.data?.message || "Error creating account",
        variant: "destructive",
      });
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  const logout = () => {
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

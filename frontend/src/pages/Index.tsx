
import React from 'react';
import { useAuth } from '@/contexts/AuthContext';
import AuthForm from '@/components/auth/AuthForm';
import ChatLayout from '@/components/chat/ChatLayout';

const Index = () => {
  const { isAuthenticated, isLoading } = useAuth();

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-primary/10 via-background to-primary/5">
        <div className="text-center space-y-4">
          <div className="mx-auto h-16 w-16 bg-gradient-to-r from-primary to-purple-600 rounded-2xl flex items-center justify-center animate-bounce-in">
            <div className="h-8 w-8 bg-white rounded-lg animate-pulse"></div>
          </div>
          <h2 className="text-2xl font-bold bg-gradient-to-r from-primary to-purple-600 bg-clip-text text-transparent">
            ChatApp
          </h2>
          <p className="text-muted-foreground">Loading...</p>
        </div>
      </div>
    );
  }

  return isAuthenticated ? <ChatLayout /> : <AuthForm />;
};

export default Index;

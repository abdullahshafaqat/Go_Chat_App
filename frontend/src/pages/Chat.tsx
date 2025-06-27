import React, { useState, useEffect, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { useAuth } from '../contexts/AuthContext';
import { chatApi } from '../services/api';
import { toast } from '@/hooks/use-toast';
import { LogOut, MessageSquare, Send, Mail, Plus, Sparkles, RefreshCw, WifiIcon } from 'lucide-react';

interface Message {
  id: string;
  sender_id: string;
  recipient_id: string;
  content: string;
  message: string;
  timestamp: string;
}

const Chat: React.FC = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const [messages, setMessages] = useState<Message[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [showNewMessage, setShowNewMessage] = useState(false);
  const [newMessage, setNewMessage] = useState({
    recipient_id: '',
    content: ''
  });

  // 🚨 Redirect to login if user is not authenticated
  useEffect(() => {
    if (!user) {
      navigate('/login');
    }
  }, [user, navigate]);

  const fetchMessages = useCallback(async () => {
    try {
      setIsLoading(true);
      const fetchedMessages = await chatApi.getMessages();

      const normalized = fetchedMessages.map((msg: any) => ({
        ...msg,
        recipient_id: msg.receiver_id,
        content: msg.content || msg.message,
      }));

      setMessages(normalized);
    } catch (error: any) {
      toast({
        title: "Failed to load messages",
        description: error.response?.data?.message || "Something went wrong",
        variant: "destructive",
      });
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchMessages();
  }, [fetchMessages]);

  const handleSendMessage = async () => {
    if (!newMessage.content.trim() || !newMessage.recipient_id.trim()) {
      toast({
        title: "Invalid message",
        description: "Please enter both recipient ID and message content",
        variant: "destructive",
      });
      return;
    }

    try {
      await chatApi.sendMessage({
        content: newMessage.content,
        reciever_id: newMessage.recipient_id
      });

      setNewMessage({ recipient_id: '', content: '' });
      setShowNewMessage(false);

      toast({
        title: "Message sent",
        description: "Message sent via HTTP endpoint",
      });

      fetchMessages();
    } catch (error: any) {
      toast({
        title: "Failed to send message",
        description: error.response?.data?.message || "Something went wrong",
        variant: "destructive",
      });
    }
  };

  const formatTimestamp = (timestamp: string) => {
    try {
      const date = new Date(timestamp);
      return date.toLocaleString();
    } catch {
      return timestamp;
    }
  };

  const getUserInitials = (userId: string | number) => {
    const userIdStr = String(userId || 'U');
    return userIdStr.substring(0, 2).toUpperCase();
  };

  const handleLogout = () => {
    logout();
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-indigo-50 via-purple-50 to-pink-50 animate-fade-in">
      <header className="bg-white/80 backdrop-blur-xl shadow-xl border-b border-gradient-to-r from-indigo-200 to-purple-200 sticky top-0 z-50 animate-slide-in-right">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-20">
            <div className="flex items-center space-x-4">
              <div className="relative bg-gradient-to-r from-indigo-600 via-purple-600 to-pink-600 p-3 rounded-xl shadow-lg animate-float">
                <MessageSquare className="h-8 w-8 text-white animate-pulse" />
              </div>
              <div>
                <h1 className="text-3xl font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent">
                  Chat Universe
                </h1>
                <p className="text-sm text-gray-500 flex items-center">
                  <Sparkles className="h-3 w-3 text-yellow-500 animate-pulse" />
                  <span>Connect • Share • Inspire</span>
                  <WifiIcon className="h-3 w-3 text-green-500 ml-2" />
                  <span className="text-xs text-green-500">Online</span>
                </p>
              </div>
            </div>
            <div className="flex items-center space-x-6">
              <div className="flex items-center space-x-3 bg-gradient-to-r from-indigo-50 to-purple-50 px-6 py-3 rounded-xl border border-indigo-200">
                <Avatar className="h-12 w-12 border-2 border-indigo-300">
                  <AvatarFallback className="bg-gradient-to-r from-indigo-100 to-purple-100 text-indigo-700">
                    <Mail className="h-6 w-6" />
                  </AvatarFallback>
                </Avatar>
                <div className="flex flex-col">
                  {user?.email && (
                    <span className="text-sm font-semibold text-gray-900 hover:text-indigo-600 transition-colors duration-300">
                      {user.email}
                    </span>
                  )}
                  <span className="text-xs text-green-500 flex items-center">
                    <div className="w-2 h-2 rounded-full mr-2 animate-pulse bg-green-500"></div>
                    connected
                  </span>
                </div>
              </div>
              <Button
                variant="outline"
                size="sm"
                onClick={handleLogout}
                className="flex items-center space-x-2 border-red-200 text-red-600 hover:bg-red-50"
              >
                <LogOut className="h-4 w-4" />
                <span>Logout</span>
              </Button>
            </div>
          </div>
        </div>
      </header>

      <div className="max-w-5xl mx-auto p-6">
        <div className="bg-white/90 backdrop-blur-sm rounded-3xl shadow-2xl border border-indigo-100">
          <div className="p-8 border-b border-gradient-to-r from-indigo-100 to-purple-100">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-2xl font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent">
                Messages
              </h2>
              <Button
                variant="outline"
                size="sm"
                onClick={fetchMessages}
                disabled={isLoading}
                className="flex items-center space-x-2"
              >
                <RefreshCw className={`h-4 w-4 ${isLoading ? 'animate-spin' : ''}`} />
                <span>Refresh</span>
              </Button>
            </div>
            {isLoading ? (
              <div className="flex justify-center py-12">
                <div className="relative">
                  <div className="animate-spin rounded-full h-12 w-12 border-4 border-indigo-200 border-t-indigo-600"></div>
                  <div className="absolute inset-0 animate-ping rounded-full h-12 w-12 border-4 border-indigo-300"></div>
                </div>
              </div>
            ) : messages.length === 0 ? (
              <div className="text-center py-16 text-gray-500">
                <MessageSquare className="h-20 w-20 mx-auto text-gray-300 animate-float mb-6" />
                <h3 className="text-xl font-semibold mb-2 text-gray-700">No messages yet</h3>
                <p>Start a conversation and connect with others!</p>
              </div>
            ) : (
              <div className="space-y-6 max-h-96 overflow-y-auto custom-scrollbar">
                {messages.map((message) => (
                  <div
                    key={message.id}
                    className="flex items-start space-x-4 p-4 bg-gradient-to-r from-gray-50 to-indigo-50 rounded-2xl hover:shadow-lg"
                  >
                    <Avatar className="h-12 w-12">
                      <AvatarFallback className="bg-gradient-to-r from-indigo-500 to-purple-500 text-white font-bold">
                        {getUserInitials(message.sender_id)}
                      </AvatarFallback>
                    </Avatar>
                    <div className="flex-1">
                      <div className="flex items-center space-x-3 mb-2">
                        <span className="text-sm font-semibold text-gray-900">User {message.sender_id}</span>
                        <span className="text-xs text-gray-500"> to user {message.recipient_id}</span>
                        <span className="text-xs text-gray-400 bg-gray-100 px-2 py-1 rounded-full">
                          {formatTimestamp(message.timestamp)}
                        </span>
                      </div>
                      <p className="text-sm text-gray-700 bg-white p-3 rounded-lg">
                        {message.content || message.message}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>

          {showNewMessage && (
            <div className="p-8 border-b bg-gradient-to-r from-indigo-50 via-purple-50 to-pink-50">
              <h3 className="text-xl font-semibold mb-6 text-gray-800 flex items-center space-x-2">
                <Send className="h-5 w-5 text-indigo-600" />
                <span>Compose New Message</span>
              </h3>
              <div className="space-y-6">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">Recipient ID</label>
                  <Input
                    value={newMessage.recipient_id}
                    onChange={(e) => setNewMessage(prev => ({ ...prev, recipient_id: e.target.value }))}
                    placeholder="Enter recipient user ID"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">Message</label>
                  <Textarea
                    value={newMessage.content}
                    onChange={(e) => setNewMessage(prev => ({ ...prev, content: e.target.value }))}
                    placeholder="Type your message here..."
                    rows={4}
                    className="resize-none"
                  />
                </div>
                <div className="flex space-x-4">
                  <Button onClick={handleSendMessage} className="bg-gradient-to-r from-indigo-600 via-purple-600 to-pink-600">
                    <Send className="h-4 w-4" />
                    <span>Send Message</span>
                  </Button>
                  <Button variant="outline" onClick={() => setShowNewMessage(false)}>
                    Cancel
                  </Button>
                </div>
              </div>
            </div>
          )}
        </div>

        {!showNewMessage && (
          <div className="fixed bottom-8 right-8">
            <Button
              onClick={() => setShowNewMessage(true)}
              className="rounded-full h-16 w-16 shadow-2xl bg-gradient-to-r from-indigo-600 via-purple-600 to-pink-600"
              size="icon"
            >
              <Plus className="h-7 w-7" />
            </Button>
          </div>
        )}
      </div>
    </div>
  );
};

export default Chat;

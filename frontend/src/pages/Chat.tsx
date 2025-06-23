import React, { useState, useEffect, useCallback } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { useAuth } from '../contexts/AuthContext';
import { chatApi } from '../api';
import { toast } from '@/hooks/use-toast';
import { LogOut, MessageSquare, Send, User, Mail, Plus, Sparkles, RefreshCw } from 'lucide-react';
import { webSocketService } from '../services/websocketservice';

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
  const [messages, setMessages] = useState<Message[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [showNewMessage, setShowNewMessage] = useState(false);
  const [newMessage, setNewMessage] = useState({
    recipient_id: '',
    content: ''
  });

  useEffect(() => {
    const token = localStorage.getItem('token');
    if (!token || !user) return;

    const handleIncomingMessage = (message: any) => {
      setMessages(prev => [...prev, {
        id: Date.now().toString(),
        sender_id: message.sender_id,
        recipient_id: user.id,
        content: message.message,
        message: message.message,
        timestamp: new Date().toISOString()
      }]);
    };

    webSocketService.connect(token, handleIncomingMessage);

    return () => {
      webSocketService.disconnect();
    };
  }, [user]);

  const fetchMessages = useCallback(async () => {
    try {
      setIsLoading(true);
      console.log('Fetching messages...');
      const fetchedMessages = await chatApi.getMessages();
      console.log('Fetched messages:', fetchedMessages);
      setMessages(fetchedMessages);
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
  
    const receiverId = parseInt(newMessage.recipient_id, 10);
    if (isNaN(receiverId)) {
      throw new Error("Recipient ID must be a valid number");
    }

    webSocketService.sendMessage({
      receiver_id: receiverId,  // Now sending as number
      message: newMessage.content
    });

    setMessages(prev => [...prev, {
      id: Date.now().toString(),
      sender_id: user?.id || '',
      recipient_id: newMessage.recipient_id,
      content: newMessage.content,
      message: newMessage.content,
      timestamp: new Date().toISOString()
    }]);

    setNewMessage({ recipient_id: '', content: '' });
    setShowNewMessage(false);
    
    toast({
      title: "Message sent",
      description: "Your message has been sent successfully",
    });
  } catch (error: any) {
    console.error('Error sending message:', error);
    toast({
      title: "Failed to send message",
      description: error.message || "Something went wrong",
      variant: "destructive",
    });
  }
};

  const formatTimestamp = (timestamp: string) => {
    try {
      const date = new Date(timestamp);
      return date.toLocaleString();
    } catch (error) {
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
            <div className="flex items-center space-x-4 transform transition-all duration-300 hover:scale-105">
              <div className="relative bg-gradient-to-r from-indigo-600 via-purple-600 to-pink-600 p-3 rounded-xl shadow-lg hover:shadow-2xl transition-all duration-500 hover:rotate-12 animate-float">
                <MessageSquare className="h-8 w-8 text-white animate-pulse" />
                <div className="absolute -top-1 -right-1 w-4 h-4 bg-green-500 rounded-full animate-bounce"></div>
              </div>
              <div className="animate-fade-in" style={{ animationDelay: '200ms' }}>
                <h1 className="text-3xl font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent hover:from-purple-600 hover:to-pink-600 transition-all duration-500">
                  Chat Universe
                </h1>
                <p className="text-sm text-gray-500 flex items-center space-x-1">
                  <Sparkles className="h-3 w-3 text-yellow-500 animate-pulse" />
                  <span>Connect • Share • Inspire</span>
                </p>
              </div>
            </div>
            
            <div className="flex items-center space-x-6">
              <div className="flex items-center space-x-3 bg-gradient-to-r from-indigo-50 to-purple-50 px-6 py-3 rounded-xl border border-indigo-200 transition-all duration-300 hover:shadow-lg hover:bg-white transform hover:scale-105 animate-fade-in" style={{ animationDelay: '300ms' }}>
                <Avatar className="h-12 w-12 border-2 border-indigo-300 transition-all duration-300 hover:border-purple-400 hover:scale-110">
                  <AvatarFallback className="bg-gradient-to-r from-indigo-100 to-purple-100 text-indigo-700">
                    <Mail className="h-6 w-6" />
                  </AvatarFallback>
                </Avatar>
                <div className="flex flex-col">
                  <span className="text-sm font-semibold text-gray-900 hover:text-indigo-600 transition-colors duration-300">
                    {user?.email || 'user@example.com'}
                  </span>
                  <span className="text-xs text-gray-500 flex items-center">
                    <div className="w-2 h-2 bg-green-500 rounded-full mr-2 animate-pulse"></div>
                    <span className="animate-fade-in">Active now</span>
                  </span>
                </div>
              </div>
              
              <Button
                variant="outline"
                size="sm"
                onClick={handleLogout}
                className="flex items-center space-x-2 border-red-200 text-red-600 hover:bg-red-50 hover:border-red-300 transition-all duration-300 hover:scale-105 active:scale-95 hover:shadow-md animate-fade-in"
                style={{ animationDelay: '400ms' }}
              >
                <LogOut className="h-4 w-4" />
                <span>Logout</span>
              </Button>
            </div>
          </div>
        </div>
      </header>

      <div className="max-w-5xl mx-auto p-6 animate-scale-in">
        <div className="bg-white/90 backdrop-blur-sm rounded-3xl shadow-2xl hover:shadow-3xl transition-all duration-700 transform hover:-translate-y-2 border border-indigo-100">
          <div className="p-8 border-b border-gradient-to-r from-indigo-100 to-purple-100">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-2xl font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent hover:from-purple-600 hover:to-pink-600 transition-all duration-300 animate-fade-in">
                Messages
              </h2>
              <div className="flex items-center space-x-4">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={fetchMessages}
                  disabled={isLoading}
                  className="flex items-center space-x-2 transition-all duration-300 hover:scale-105"
                >
                  <RefreshCw className={`h-4 w-4 ${isLoading ? 'animate-spin' : ''}`} />
                  <span>Refresh</span>
                </Button>
                <div className="flex items-center space-x-2 animate-fade-in" style={{ animationDelay: '200ms' }}>
                  <div className="w-3 h-3 bg-blue-500 rounded-full"></div>
                  <span className="text-sm text-gray-500">Manual refresh</span>
                </div>
              </div>
            </div>
            
            {isLoading ? (
              <div className="flex justify-center py-12">
                <div className="relative">
                  <div className="animate-spin rounded-full h-12 w-12 border-4 border-indigo-200 border-t-indigo-600"></div>
                  <div className="absolute inset-0 animate-ping rounded-full h-12 w-12 border-4 border-indigo-300"></div>
                </div>
              </div>
            ) : messages.length === 0 ? (
              <div className="text-center py-16 text-gray-500 animate-fade-in">
                <div className="mb-6 transform transition-all duration-500 hover:scale-110">
                  <MessageSquare className="h-20 w-20 mx-auto text-gray-300 animate-float" />
                </div>
                <h3 className="text-xl font-semibold mb-2 text-gray-700">No messages yet</h3>
                <p className="text-gray-500">Start a conversation and connect with others!</p>
              </div>
            ) : (
              <div className="space-y-6 max-h-96 overflow-y-auto custom-scrollbar">
                {messages.map((message, index) => (
                  <div 
                    key={message.id} 
                    className="flex items-start space-x-4 p-4 bg-gradient-to-r from-gray-50 to-indigo-50 rounded-2xl transition-all duration-500 hover:from-indigo-50 hover:to-purple-50 hover:shadow-lg transform hover:scale-102 animate-fade-in hover:-translate-y-1"
                    style={{ animationDelay: `${index * 100}ms` }}
                  >
                    <Avatar className="h-12 w-12 transition-all duration-300 hover:scale-125 shadow-lg">
                      <AvatarFallback className="bg-gradient-to-r from-indigo-500 to-purple-500 text-white font-bold">
                        {getUserInitials(message.sender_id)}
                      </AvatarFallback>
                    </Avatar>
                    
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center space-x-3 mb-2">
                        <span className="text-sm font-semibold text-gray-900 hover:text-indigo-600 transition-colors duration-300">
                          User {message.sender_id}
                        </span>
                        <div className="flex items-center space-x-1 text-xs text-gray-500">
                          <span>to</span>
                          <span className="font-medium">User {message.recipient_id}</span>
                        </div>
                        <span className="text-xs text-gray-400 bg-gray-100 px-2 py-1 rounded-full">
                          {formatTimestamp(message.timestamp)}
                        </span>
                      </div>
                      <p className="text-sm text-gray-700 leading-relaxed bg-white p-3 rounded-lg shadow-sm">
                        {message.content || message.message}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>

          {showNewMessage && (
            <div className="p-8 border-b bg-gradient-to-r from-indigo-50 via-purple-50 to-pink-50 animate-fade-in">
              <h3 className="text-xl font-semibold mb-6 text-gray-800 flex items-center space-x-2">
                <Send className="h-5 w-5 text-indigo-600" />
                <span>Compose New Message</span>
              </h3>
              <div className="space-y-6">
                <div className="transform transition-all duration-300 hover:scale-105 animate-slide-in-right" style={{ animationDelay: '100ms' }}>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Recipient ID
                  </label>
                  <Input
                    value={newMessage.recipient_id}
                    onChange={(e) => setNewMessage(prev => ({ ...prev, recipient_id: e.target.value }))}
                    placeholder="Enter recipient user ID"
                    className="transition-all duration-300 focus:scale-105 hover:shadow-lg border-indigo-200 focus:border-indigo-500"
                  />
                </div>
                
                <div className="transform transition-all duration-300 hover:scale-105 animate-slide-in-right" style={{ animationDelay: '200ms' }}>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Message
                  </label>
                  <Textarea
                    value={newMessage.content}
                    onChange={(e) => setNewMessage(prev => ({ ...prev, content: e.target.value }))}
                    placeholder="Type your message here..."
                    rows={4}
                    className="transition-all duration-300 focus:scale-105 hover:shadow-lg resize-none border-indigo-200 focus:border-indigo-500"
                  />
                </div>
                
                <div className="flex space-x-4 animate-fade-in" style={{ animationDelay: '300ms' }}>
                  <Button 
                    onClick={handleSendMessage} 
                    className="flex items-center space-x-2 bg-gradient-to-r from-indigo-600 via-purple-600 to-pink-600 hover:from-indigo-700 hover:via-purple-700 hover:to-pink-700 transition-all duration-500 hover:scale-110 active:scale-95 hover:shadow-xl transform hover:-translate-y-1"
                  >
                    <Send className="h-4 w-4" />
                    <span>Send Message</span>
                  </Button>
                  <Button 
                    variant="outline" 
                    onClick={() => {
                      setShowNewMessage(false);
                      setNewMessage({ recipient_id: '', content: '' });
                    }}
                    className="transition-all duration-300 hover:scale-105 active:scale-95 hover:shadow-lg border-gray-300"
                  >
                    Cancel
                  </Button>
                </div>
              </div>
            </div>
          )}
        </div>

        {!showNewMessage && (
          <div className="fixed bottom-8 right-8 animate-fade-in">
            <Button
              onClick={() => setShowNewMessage(true)}
              className="rounded-full h-16 w-16 shadow-2xl hover:shadow-3xl bg-gradient-to-r from-indigo-600 via-purple-600 to-pink-600 hover:from-indigo-700 hover:via-purple-700 hover:to-pink-700 transition-all duration-500 hover:scale-125 active:scale-95 animate-float hover:animate-none group"
              size="icon"
            >
              <Plus className="h-7 w-7 transition-transform duration-300 group-hover:rotate-180" />
            </Button>
          </div>
        )}
      </div>
    </div>
  );
};

export default Chat;
// src/pages/Chat.tsx
import React, { useState, useEffect, useCallback } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { useAuth } from "../contexts/AuthContext";
import { chatApi } from "../services/api";
import { toast } from "@/hooks/use-toast";
import { LogOut, MessageSquare, Plus, Sparkles, Send } from "lucide-react";
import { webSocketService } from "../services/websocketservice";

interface Message {
  sender_id: number;
  receiver_id: number;
  message: string;
  timestamp: string;
}

const Chat: React.FC = () => {
  const { user, logout } = useAuth();
  const [messages, setMessages] = useState<Message[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [showNewMessage, setShowNewMessage] = useState(false);
  const [newMessage, setNewMessage] = useState({ receiver_id: "", message: "" });

  // Fetch messages initially
  const fetchMessages = useCallback(async () => {
    setIsLoading(true);
    try {
      const fetched = await chatApi.getMessages();
      setMessages(fetched);
    } catch (error: any) {
      toast({ title: "Error loading messages", description: error.message, variant: "destructive" });
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchMessages();
  }, [fetchMessages]);

  // Listen for new incoming messages
  useEffect(() => {
    webSocketService.setOnMessageCallback((incomingMessage) => {
      setMessages((prev) => [...prev, incomingMessage]);
    });
  }, []);

  const handleSendMessage = async () => {
    if (!newMessage.message.trim() || !newMessage.receiver_id.trim()) {
      toast({ title: "Error", description: "Please enter recipient ID and message.", variant: "destructive" });
      return;
    }
    const receiver_id = parseInt(newMessage.receiver_id, 10);
    if (isNaN(receiver_id)) {
      toast({ title: "Error", description: "Receiver ID must be a number.", variant: "destructive" });
      return;
    }

    // send via WS
    webSocketService.sendMessage({ receiver_id, message: newMessage.message });

    // store via API
    await chatApi.sendMessage({ content: newMessage.message, reciever_id: receiver_id });

    // optimistic UI update
    setMessages((prev) => [
      ...prev,
      { sender_id: user!.id, receiver_id, message: newMessage.message, timestamp: new Date().toISOString() },
    ]);

    setNewMessage({ receiver_id: "", message: "" });
    setShowNewMessage(false);
  };

  const formatTimestamp = (ts: string) => new Date(ts).toLocaleString();
  const getUserInitials = (userId: number) => String(userId).substring(0, 2).toUpperCase();

  return (
    <div className="min-h-screen bg-gradient-to-br from-indigo-50 via-purple-50 to-pink-50 animate-fade-in">
      <header className="bg-white/80 backdrop-blur-xl shadow-xl border-b border-gradient-to-r from-indigo-200 to-purple-200 sticky top-0 z-50 animate-slide-in-right">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 flex justify-between items-center h-20">
          <div className="flex items-center space-x-4 transform transition-all duration-300 hover:scale-105">
            <div className="relative bg-gradient-to-r from-indigo-600 via-purple-600 to-pink-600 p-3 rounded-xl shadow-lg hover:shadow-2xl transition-all duration-500 hover:rotate-12 animate-float">
              <MessageSquare className="h-8 w-8 text-white animate-pulse" />
              <div className="absolute -top-1 -right-1 w-4 h-4 bg-green-500 rounded-full animate-bounce"></div>
            </div>
            <div className="animate-fade-in" style={{ animationDelay: "200ms" }}>
              <h1 className="text-3xl font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent">Chat Universe</h1>
              <p className="text-sm text-gray-500 flex items-center space-x-1">
                <Sparkles className="h-3 w-3 text-yellow-500 animate-pulse" />
                <span>Connect • Share • Inspire</span>
              </p>
            </div>
          </div>
          <Button variant="outline" size="sm" onClick={logout} className="flex items-center space-x-2 border-red-200 text-red-600 hover:bg-red-50">
            <LogOut className="h-4 w-4" /> Logout
          </Button>
        </div>
      </header>

      <div className="max-w-7xl mx-auto p-6">
        <div className="bg-white/90 backdrop-blur-sm rounded-3xl shadow-2xl border border-indigo-100 p-8">
          <div className="flex items-center justify-between mb-6 border-b border-gradient-to-r from-indigo-100 to-purple-100 pb-4">
            <h2 className="text-2xl font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent">Messages</h2>
          </div>

          {isLoading ? (
            <div className="flex justify-center py-12">
              <div className="animate-spin rounded-full h-12 w-12 border-4 border-indigo-200 border-t-indigo-600"></div>
            </div>
          ) : messages.length === 0 ? (
            <div className="text-center py-16 text-gray-500">
              <MessageSquare className="h-20 w-20 mx-auto text-gray-300" />
              <h3 className="text-xl font-semibold mt-4">No messages yet</h3>
              <p>Start a conversation and connect with others!</p>
            </div>
          ) : (
            <div className="space-y-6 max-h-96 overflow-y-auto">
              {messages.map((message, index) => (
                <div key={`${message.timestamp}-${index}`} className="flex items-start space-x-4 p-4 bg-gradient-to-r from-gray-50 to-indigo-50 rounded-2xl hover:shadow-lg">
                  <Avatar className="h-12 w-12">
                    <AvatarFallback className="bg-gradient-to-r from-indigo-500 to-purple-500 text-white font-bold">
                      {getUserInitials(message.sender_id)}
                    </AvatarFallback>
                  </Avatar>
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center space-x-3 mb-2 text-xs text-gray-500">
                      <span>User {message.sender_id} → User {message.receiver_id}</span>
                      <span className="bg-gray-100 px-2 py-1 rounded-full">{formatTimestamp(message.timestamp)}</span>
                    </div>
                    <p className="text-sm text-gray-700 bg-white p-3 rounded-lg shadow-sm">{message.message}</p>
                  </div>
                </div>
              ))}
            </div>
          )}

          {showNewMessage && (
            <div className="p-8 mt-6 border-t bg-gradient-to-r from-indigo-50 via-purple-50 to-pink-50 rounded-xl shadow-inner">
              <h3 className="text-xl font-semibold mb-6 flex items-center space-x-2 text-indigo-600">
                <Send className="h-5 w-5" /><span>Compose New Message</span>
              </h3>
              <div className="space-y-4">
                <Input
                  value={newMessage.receiver_id}
                  onChange={(e) => setNewMessage((prev) => ({ ...prev, receiver_id: e.target.value }))}
                  placeholder="Recipient User ID"
                  className="shadow-sm"
                />
                <Textarea
                  value={newMessage.message}
                  onChange={(e) => setNewMessage((prev) => ({ ...prev, message: e.target.value }))}
                  placeholder="Your message"
                  className="shadow-sm"
                  rows={4}
                />
                <div className="flex space-x-4">
                  <Button onClick={handleSendMessage} className="bg-gradient-to-r from-indigo-600 to-pink-600 text-white shadow-lg">Send Message</Button>
                  <Button variant="outline" onClick={() => setShowNewMessage(false)}>Cancel</Button>
                </div>
              </div>
            </div>
          )}

          {!showNewMessage && (
            <div className="fixed bottom-8 right-8">
              <Button onClick={() => setShowNewMessage(true)} className="rounded-full h-16 w-16 bg-gradient-to-r from-indigo-600 via-purple-600 to-pink-600 text-white shadow-2xl hover:scale-110 transition-transform">
                <Plus className="h-7 w-7" />
              </Button>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default Chat;

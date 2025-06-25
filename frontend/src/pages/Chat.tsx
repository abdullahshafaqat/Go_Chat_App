import React, { useState, useEffect, useCallback } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { useAuth } from "../contexts/AuthContext";
import { chatApi } from "../services/api";
import { toast } from "@/hooks/use-toast";
import { LogOut, MessageSquare, Send, Mail, Plus, Sparkles, RefreshCw } from "lucide-react";
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

  // Fetch old messages
  const fetchMessages = useCallback(async () => {
    try {
      setIsLoading(true);
      const fetchedMessages = await chatApi.getMessages();
      setMessages(fetchedMessages);
    } catch (error: any) {
      toast({
        title: "Error",
        description: error?.message || "Error fetching messages",
        variant: "destructive",
      });
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchMessages();
  }, [fetchMessages]);

  // WebSocket connect
  useEffect(() => {
    const token = localStorage.getItem("token");
    if (!token || !user) return;

    webSocketService.connect(token, (message: Message) => {
      setMessages((prev) => [...prev, message]);
    });

    return () => webSocketService.disconnect();
  }, [user]);

  // Handle send message
  const handleSendMessage = async () => {
    if (!newMessage.message.trim() || !newMessage.receiver_id.trim()) {
      return toast({ title: "Error", description: "Recipient and message required", variant: "destructive" });
    }

    const receiver_id = parseInt(newMessage.receiver_id, 10);
    if (isNaN(receiver_id)) {
      return toast({ title: "Error", description: "Recipient must be a number.", variant: "destructive" });
    }

    try {
      // Save to DB
      await chatApi.sendMessage({ content: newMessage.message, reciever_id: receiver_id });

      // Real-time WS send
      webSocketService.sendMessage({ receiver_id, message: newMessage.message });

      setMessages((prev) => [
        ...prev,
        { sender_id: user!.id, receiver_id, message: newMessage.message, timestamp: new Date().toISOString() },
      ]);

      setNewMessage({ receiver_id: "", message: "" });
      setShowNewMessage(false);

      toast({ title: "Sent successfully!" });
    } catch (e) {
      toast({ title: "Error sending message.", variant: "destructive" });
      console.error(e);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-indigo-50 via-purple-50 to-pink-50 animate-fade-in">
      <header className="bg-white/80 backdrop-blur-xl shadow-xl border-b border-gradient-to-r from-indigo-200 to-purple-200 sticky top-0 z-50 animate-slide-in-right">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-20">
            <div className="flex items-center space-x-4">
              <MessageSquare className="h-8 w-8 text-purple-600" />
              <h1 className="text-3xl font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent">Chat Universe</h1>
            </div>
            <div className="flex items-center space-x-4">
              <div className="flex items-center space-x-2 bg-white px-4 py-2 rounded-xl shadow border border-indigo-200">
                <Avatar className="h-10 w-10 border-2 border-indigo-300">
                  <AvatarFallback className="bg-purple-100 text-purple-700"><Mail /></AvatarFallback>
                </Avatar>
                <div>
                  <span className="text-sm font-semibold text-gray-900">{user?.email}</span>
                  <span className="text-xs text-green-500 block">Online</span>
                </div>
              </div>
              <Button onClick={logout} variant="outline" className="text-red-500 border-red-200 hover:bg-red-50 flex items-center space-x-1">
                <LogOut size={16} /> <span>Logout</span>
              </Button>
            </div>
          </div>
        </div>
      </header>

      <main className="max-w-5xl mx-auto p-6">
        <section className="bg-white/90 backdrop-blur-md border border-indigo-100 rounded-2xl shadow-xl p-6">
          <h2 className="text-2xl font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent mb-4">Messages</h2>

          {isLoading ? (
            <div className="flex justify-center p-10">
              <div className="animate-spin border-4 border-indigo-200 border-t-indigo-600 h-12 w-12 rounded-full"></div>
            </div>
          ) : messages.length === 0 ? (
            <div className="text-center text-gray-500 py-10">
              <MessageSquare className="h-20 w-20 mx-auto text-gray-300" />
              <h3 className="text-xl mt-4">No messages yet</h3>
            </div>
          ) : (
            <div className="space-y-4 max-h-96 overflow-y-auto pr-2">
              {messages.map((message, index) => (
                <div key={index} className="bg-gradient-to-r from-indigo-50 to-purple-50 p-4 rounded-xl shadow-sm hover:shadow-md transition-all flex space-x-3">
                  <Avatar className="h-10 w-10 border-2 border-indigo-200">
                    <AvatarFallback className="bg-purple-500 text-white">{message.sender_id}</AvatarFallback>
                  </Avatar>
                  <div className="flex-1">
                    <div className="text-xs text-gray-500 mb-1">
                      User {message.sender_id} → {message.receiver_id} @ {new Date(message.timestamp).toLocaleString()}
                    </div>
                    <p className="bg-white p-3 rounded-lg text-gray-800 shadow-inner">{message.message}</p>
                  </div>
                </div>
              ))}
            </div>
          )}

          {showNewMessage ? (
            <div className="mt-8 bg-gradient-to-r from-indigo-50 via-purple-50 to-pink-50 p-6 rounded-xl border border-indigo-100">
              <h3 className="text-lg font-medium mb-4 flex items-center"><Send className="h-5 w-5 text-indigo-600 mr-2" /> Compose Message</h3>
              <Input
                placeholder="Recipient ID"
                value={newMessage.receiver_id}
                onChange={(e) => setNewMessage((m) => ({ ...m, receiver_id: e.target.value }))}
                className="mb-3"
              />
              <Textarea
                placeholder="Your message..."
                value={newMessage.message}
                onChange={(e) => setNewMessage((m) => ({ ...m, message: e.target.value }))}
                className="mb-3"
              />
              <div className="flex space-x-4">
                <Button onClick={handleSendMessage} className="bg-gradient-to-r from-indigo-600 to-pink-600 text-white">Send</Button>
                <Button variant="outline" onClick={() => setShowNewMessage(false)}>Cancel</Button>
              </div>
            </div>
          ) : (
            <Button
              onClick={() => setShowNewMessage(true)}
              className="fixed bottom-8 right-8 h-16 w-16 bg-gradient-to-r from-indigo-600 via-purple-600 to-pink-600 text-white rounded-full shadow-xl"
            >
              <Plus className="h-7 w-7" />
            </Button>
          )}
        </section>
      </main>
    </div>
  );
};

export default Chat;

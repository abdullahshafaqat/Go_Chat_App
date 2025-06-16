import React, { useState, useEffect, useRef } from 'react';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Button } from '@/components/ui/button';
import { RefreshCw, ArrowDown } from 'lucide-react';
import MessageBubble from './MessageBubble';
import MessageInput from './MessageInput';
import { chatApi } from '@/services/api';
import { toast } from '@/hooks/use-toast';
import { useAuth } from '@/contexts/AuthContext';

interface Message {
  id: string;
  content: string;
  sender_id: string;
  sender_username: string;
  receiver_id?: string;  // Added receiver_id to the interface
  timestamp: string;
  created_at?: string;
}

interface ChatWindowProps {
  recipientId?: string;  // Added recipientId as optional prop
}

const ChatWindow: React.FC<ChatWindowProps> = ({ recipientId }) => {
  const [messages, setMessages] = useState<Message[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isRefreshing, setIsRefreshing] = useState(false);
  const scrollAreaRef = useRef<HTMLDivElement>(null);
  const [showScrollButton, setShowScrollButton] = useState(false);
  const { user } = useAuth();

  const loadMessages = async (showRefreshToast = false) => {
    try {
      setIsRefreshing(true);
      let response;
      
      if (recipientId) {
        // Load conversation with specific recipient
        response = await chatApi.getConversation(recipientId);
      } else {
        // Load general messages
        response = await chatApi.getMessages();
      }
      
      const messageList = response.messages || response.data || response || [];
      setMessages(messageList);
      
      if (showRefreshToast) {
        toast({
          title: "Messages refreshed",
          description: "Chat has been updated with latest messages.",
        });
      }
    } catch (error) {
      console.error('Failed to load messages:', error);
      
    } finally {
      setIsLoading(false);
      setIsRefreshing(false);
    }
  };

  const scrollToBottom = () => {
    if (scrollAreaRef.current) {
      const scrollContainer = scrollAreaRef.current.querySelector('[data-radix-scroll-area-viewport]');
      if (scrollContainer) {
        scrollContainer.scrollTop = scrollContainer.scrollHeight;
      }
    }
  };

  const handleScroll = (e: React.UIEvent) => {
    const target = e.target as HTMLElement;
    const isNearBottom = target.scrollHeight - target.scrollTop - target.clientHeight < 100;
    setShowScrollButton(!isNearBottom);
  };

  useEffect(() => {
    loadMessages();
  }, [recipientId]);  // Reload when recipientId changes

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const handleMessageUpdate = async (messageId: string, newContent: string) => {
    try {
      setMessages(prevMessages =>
        prevMessages.map(msg =>
          msg.id === messageId ? { ...msg, content: newContent } : msg
        )
      );
      
      toast({
        title: "Message updated",
        description: "Your message has been updated successfully.",
      });
    } catch (error) {
      console.error('Failed to update message:', error);
      
      loadMessages();
    }
  };

  const handleMessageSent = () => {
    loadMessages();
    scrollToBottom();
  };

  const mockMessages: Message[] = [
    
    {
      id: '3',
      content: 'Try typing a message below to see the chat in action.',
      sender_id: 'system',
      sender_username: 'System',
      timestamp: new Date(Date.now() - 60000).toISOString(),
    }
  ];

  const displayMessages = messages.length > 0 ? messages : mockMessages;

  if (isLoading) {
    return (
      <div className="flex-1 flex items-center justify-center">
        <div className="text-center space-y-4">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto"></div>
          <p className="text-muted-foreground">Loading messages...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex-1 flex flex-col bg-chat-gradient">
      <div className="flex-1 relative">
        <ScrollArea
          ref={scrollAreaRef}
          className="h-full"
          onScrollCapture={handleScroll}
        >
          <div className="p-6 space-y-4">
            <div className="flex justify-center mb-6">
              <Button
                variant="outline"
                size="sm"
                onClick={() => loadMessages(true)}
                disabled={isRefreshing}
                className="text-xs hover:scale-105 transition-all"
              >
                <RefreshCw className={`h-3 w-3 mr-2 ${isRefreshing ? 'animate-spin' : ''}`} />
                Refresh Messages
              </Button>
            </div>
            
            {displayMessages.map((message) => (
              <MessageBubble
                key={message.id}
                message={message}
                onUpdate={handleMessageUpdate}
              />
            ))}
          </div>
        </ScrollArea>

        {showScrollButton && (
          <Button
            onClick={scrollToBottom}
            className="absolute bottom-4 right-4 h-10 w-10 rounded-full shadow-lg hover:scale-110 transition-all"
            size="sm"
          >
            <ArrowDown className="h-4 w-4" />
          </Button>
        )}
      </div>

      <MessageInput 
        onMessageSent={handleMessageSent}
        recipientId={recipientId}  // Pass recipientId to MessageInput
      />
    </div>
  );
};

export default ChatWindow;
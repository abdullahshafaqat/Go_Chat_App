import React, { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { useAuth } from '@/contexts/AuthContext';
import { Edit2, Check, X } from 'lucide-react';
import { chatApi } from '@/services/api';
import { toast } from '@/hooks/use-toast';

interface Message {
  id: string;
  content: string;
  sender_id: string;
  sender_username: string;
  timestamp: string;
}

interface MessageBubbleProps {
  message: Message;
  onUpdate: (messageId: string, newContent: string) => void;
}

const MessageBubble: React.FC<MessageBubbleProps> = ({ message, onUpdate }) => {
  const { user } = useAuth();
  const [isEditing, setIsEditing] = useState(false);
  const [editContent, setEditContent] = useState(message.content);
  const [isUpdating, setIsUpdating] = useState(false);

  const isOwnMessage = message.sender_id === user?.id;

  const handleSaveEdit = async () => {
    if (!editContent.trim()) {
      toast({
        title: "Message cannot be empty",
        variant: "destructive",
      });
      return;
    }

    if (editContent === message.content) {
      setIsEditing(false);
      return;
    }

    setIsUpdating(true);
    try {
      await chatApi.updateMessage(message.id, editContent.trim());
      onUpdate(message.id, editContent.trim());
      setIsEditing(false);
      toast({
        title: "Message updated",
        description: "Your message has been updated successfully.",
      });
    } catch (error) {
      toast({
        title: "Failed to update message",
        description: error.response?.data?.error || "Please try again.",
        variant: "destructive",
      });
      setEditContent(message.content); // Revert on error
    } finally {
      setIsUpdating(false);
    }
  };

  const handleCancelEdit = () => {
    setIsEditing(false);
    setEditContent(message.content);
  };

  const formatTime = (timestamp: string) => {
    return new Date(timestamp).toLocaleTimeString([], {
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  return (
    <div className={`flex mb-4 ${isOwnMessage ? 'justify-end' : 'justify-start'}`}>
      <div className={`flex max-w-[70%] ${isOwnMessage ? 'flex-row-reverse' : 'flex-row'}`}>
        {!isOwnMessage && (
          <Avatar className="h-8 w-8 mt-2 mr-3">
            <AvatarFallback className="text-xs">
              {message.sender_username?.charAt(0).toUpperCase() || 'U'}
            </AvatarFallback>
          </Avatar>
        )}
        
        <div className={`group ${isOwnMessage ? 'mr-3' : ''}`}>
          <div className={`flex items-center mb-1 ${isOwnMessage ? 'justify-end' : 'justify-start'}`}>
            <span className="text-xs text-muted-foreground font-medium">
              {isOwnMessage ? 'You' : message.sender_username}
            </span>
            <span className="text-xs text-muted-foreground ml-2">
              {formatTime(message.timestamp)}
            </span>
          </div>
          
          <div className={`relative rounded-2xl px-4 py-3 shadow-sm transition-all group-hover:shadow-md ${
            isOwnMessage 
              ? 'bg-primary text-primary-foreground' 
              : 'bg-secondary text-foreground'
          }`}>
            {isEditing ? (
              <div className="space-y-2">
                <Input
                  value={editContent}
                  onChange={(e) => setEditContent(e.target.value)}
                  className="border-0 bg-background/20 text-inherit"
                  disabled={isUpdating}
                  autoFocus
                />
                <div className="flex space-x-2">
                  <Button
                    size="sm"
                    onClick={handleSaveEdit}
                    disabled={isUpdating || !editContent.trim()}
                    className="h-6 px-2"
                  >
                    <Check className="h-3 w-3" />
                  </Button>
                  <Button
                    size="sm"
                    variant="ghost"
                    onClick={handleCancelEdit}
                    disabled={isUpdating}
                    className="h-6 px-2"
                  >
                    <X className="h-3 w-3" />
                  </Button>
                </div>
              </div>
            ) : (
              <>
                <p className="text-sm leading-relaxed break-words">{message.content}</p>
                {isOwnMessage && (
                  <Button
                    size="sm"
                    variant="ghost"
                    onClick={() => setIsEditing(true)}
                    className="absolute -top-2 -right-2 h-6 w-6 rounded-full bg-background/80 hover:bg-background opacity-0 group-hover:opacity-100 transition-all"
                  >
                    <Edit2 className="h-3 w-3" />
                  </Button>
                )}
              </>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default MessageBubble;
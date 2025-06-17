import React, { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';
import { Send, Loader2 } from 'lucide-react';
import { chatApi } from '@/services/api';
import { toast } from '@/hooks/use-toast';

interface MessageInputProps {
  onMessageSent: () => void;
  recipientId?: string;
}

const MessageInput: React.FC<MessageInputProps> = ({ onMessageSent, recipientId }) => {
  const [message, setMessage] = useState('');
  const [isSending, setIsSending] = useState(false);

  const handleSend = async () => {
    if (message.trim() && !isSending) {
      setIsSending(true);
      try {
        await chatApi.sendMessage({ 
          content: message.trim(),
          recipient_id: recipientId 
        });
        setMessage('');
        onMessageSent();
        toast({
          title: "Message sent",
          description: "Your message has been delivered.",
        });
      } catch (error) {
        toast({
          title: "Failed to send message",
          description: "Please try again.",
          variant: "destructive",
        });
      } finally {
        setIsSending(false);
      }
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  return (
    <div className="border-t border-border/40 p-4 bg-card/50 backdrop-blur-sm">
      <div className="flex space-x-3 items-end">
        <div className="flex-1">
          <Textarea
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            onKeyDown={handleKeyDown}
            placeholder="Type your message... (Press Enter to send)"
            className="min-h-[44px] max-h-32 resize-none border-border/60 focus:border-primary/60 transition-all focus:scale-[1.01]"
            disabled={isSending}
          />
        </div>
        <Button
          onClick={handleSend}
          disabled={!message.trim() || isSending}
          className="h-11 px-4 gradient-bg hover:scale-105 transition-all disabled:scale-100"
        >
          {isSending ? (
            <Loader2 className="h-4 w-4 animate-spin" />
          ) : (
            <Send className="h-4 w-4" />
          )}
        </Button>
      </div>
    </div>
  );
};

export default MessageInput;
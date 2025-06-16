
import React from 'react';
import ChatHeader from './ChatHeader';
import ChatWindow from './ChatWindow';

const ChatLayout: React.FC = () => {
  return (
    <div className="min-h-screen flex flex-col bg-background">
      <ChatHeader />
      <ChatWindow />
    </div>
  );
};

export default ChatLayout;

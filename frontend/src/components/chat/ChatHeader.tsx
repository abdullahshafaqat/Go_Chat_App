
import React from 'react';
import { Button } from '@/components/ui/button';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { useAuth } from '@/contexts/AuthContext';
import { useTheme } from '@/contexts/ThemeContext';
import { LogOut, Moon, Sun, Settings } from 'lucide-react';

const ChatHeader: React.FC = () => {
  const { user, logout } = useAuth();
  const { theme, toggleTheme } = useTheme();

  return (
    <div className="h-16 border-b border-border/40 px-6 flex items-center justify-between bg-card/50 backdrop-blur-sm">
      <div className="flex items-center space-x-4">
        <div className="h-8 w-8 bg-gradient-to-r from-primary to-purple-600 rounded-lg flex items-center justify-center">
          <div className="h-4 w-4 bg-white rounded-sm"></div>
        </div>
        <div>
          <h1 className="text-xl font-semibold">ChatApp</h1>
          <p className="text-xs text-muted-foreground">Connected</p>
        </div>
      </div>

      <div className="flex items-center space-x-2">
        <Button
          variant="ghost"
          size="sm"
          onClick={toggleTheme}
          className="h-9 w-9 rounded-lg hover:scale-110 transition-all"
        >
          {theme === 'light' ? (
            <Moon className="h-4 w-4" />
          ) : (
            <Sun className="h-4 w-4" />
          )}
        </Button>

        <div className="flex items-center space-x-3 ml-4">
          <Avatar className="h-8 w-8 ring-2 ring-primary/20">
            <AvatarFallback className="text-sm font-medium">
              {user?.username?.charAt(0).toUpperCase() || 'U'}
            </AvatarFallback>
          </Avatar>
          <div className="hidden sm:block">
            <p className="text-sm font-medium">{user?.username}</p>
            <p className="text-xs text-muted-foreground">{user?.email}</p>
          </div>
        </div>

        <Button
          variant="ghost"
          size="sm"
          onClick={logout}
          className="h-9 w-9 rounded-lg text-destructive hover:text-destructive hover:bg-destructive/10 hover:scale-110 transition-all ml-2"
        >
          <LogOut className="h-4 w-4" />
        </Button>
      </div>
    </div>
  );
};

export default ChatHeader;

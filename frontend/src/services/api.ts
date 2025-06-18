
import axios, { AxiosInstance, AxiosResponse } from 'axios';
import { toast } from '@/hooks/use-toast';

const BASE_URL = 'http://localhost:8003'; 

// Auth state management
let authToken: string | null = localStorage.getItem('token');
let isRefreshing = false;
let failedQueue: any[] = [];

const processQueue = (error: any, token: string | null = null) => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token);
    }
  });
  
  failedQueue = [];
};

const api: AxiosInstance = axios.create({
  baseURL: BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});


api.interceptors.request.use(
  (config) => {
    if (authToken) {
      config.headers.Authorization = `Bearer ${authToken}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);


api.interceptors.response.use(
  response => response,
  async error => {
    const originalRequest = error.config;

    if (error.response?.status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        }).then(token => {
          originalRequest.headers.Authorization = `Bearer ${token}`;
          return api(originalRequest);
        });
      }

      originalRequest._retry = true;
      isRefreshing = true;

      try {
        const refreshToken = localStorage.getItem('refreshToken');
        if (!refreshToken) throw new Error("No refresh token available");
        
        // Match your backend endpoint
        const response = await axios.post(`${BASE_URL}/refresh`, { 
          refresh_token: refreshToken 
        });
        
        // Match your backend response format
        const newAccessToken = response.data.access_token;
        setAuthToken(newAccessToken);
        processQueue(null, newAccessToken);
        
        originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;
        return api(originalRequest);
      } catch (refreshError) {
        processQueue(refreshError, null);
        clearAuth();
        window.location.href = '/login';
        toast({
          title: "Session Expired",
          description: "Please log in again",
          variant: "destructive",
        });
        return Promise.reject(refreshError);
      } finally {
        isRefreshing = false;
      }
    }

    return Promise.reject(error);
  }
);


export const setAuthToken = (token: string) => {
  authToken = token;
  localStorage.setItem('token', token);
};

export const clearAuth = () => {
  authToken = null;
  localStorage.removeItem('token');
  localStorage.removeItem('refreshToken');
};

export const getAuthToken = () => authToken;


export const authApi = {
  signup: async (data: { email: string; password: string; username: string }) => {
    const response = await api.post('/signup', data);
    return response.data;
  },
  
   login: async (data: { email: string; password: string }) => {
    const response = await api.post('/login', data);
    if (response.data.access_token) {  // Changed from 'token' to 'access_token'
      setAuthToken(response.data.access_token);
      if (response.data.refresh_token) {
        localStorage.setItem('refreshToken', response.data.refresh_token);
      }
    }
    return response.data;
  },
  
 refresh: async () => {  // Changed from refreshKey to refresh
    const refreshToken = localStorage.getItem('refreshToken');
    const response = await api.post('/refresh', {  // Changed from /refresh_key
      refresh_token: refreshToken 
    });
    return response.data;
  },
};

type Message = {
  id: string;
  sender_id: string;
  recipient_id: string;
  message: string;
  timestamp: string;
  // Add or adjust fields as needed to match your backend response
};

export const chatApi = {
    getMessages: async () => {
    const response = await api.get("/get_message");
    // Map response to your frontend structure if needed
    return response.data.messages.map((msg: any) => ({
      ...msg,
      content: msg.message // Map backend 'message' field to frontend 'content'
    }));
  },
   getConversation: async (recipientId: string): Promise<Message[]> => {
    const response = await api.get(`/conversation/${recipientId}`);
    return response.data;
  },


  
   sendMessage: async (data: { 
  content: string; 
  recipient_id?: string | number 
}) => {
  const response = await api.post('/messages', {
    content: data.content,
    receiver_id: data.recipient_id ? Number(data.recipient_id) : undefined // FIXED!
  });
  return response.data;
},

  
  updateMessage: async (messageId: string, content: string) => {
    try {
      const response = await api.put(`/update/${messageId}`, {
        message: content // Match backend request structure
      });
      return response.data;
    } catch (error) {
      console.error("Failed to update message:", error);
      throw error;
    }
  },
};

export default api;
export const userApi = {
  getUser: async (userId: string) => {
    const response = await api.get(`/user/${userId}`);
    return response.data;
  },
  
  updateUser: async (userId: string, data: { username?: string; email?: string }) => {
    const response = await api.put(`/user/${userId}`, data);
    return response.data;
  },
  
  deleteUser: async (userId: string) => {
    const response = await api.delete(`/user/${userId}`);
    return response.data;
  },
};
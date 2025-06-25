
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
        
        const response = await axios.post(`${BASE_URL}/refresh`, { 
          refresh_token: refreshToken 
        });
        
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
    if (response.data.access_token) {
      setAuthToken(response.data.access_token);
      if (response.data.refresh_token) {
        localStorage.setItem('refreshToken', response.data.refresh_token);
      }
    }
    return response.data;
  },
  
  refresh: async () => {
    const refreshToken = localStorage.getItem('refreshToken');
    const response = await api.post('/refresh', {
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
};

export const chatApi = {
 getMessages: async () => {
  const response = await api.get("/get_message");
  return response.data.messages.map((msg: any) => ({
    id: msg.id,
    sender_id: msg.sender_id,
    receiver_id: msg.receiver_id,
    message: msg.message, // <- mapping this to "content"
    timestamp: msg.timestamp,
  }));
},


  getConversation: async (receiverId: string | number): Promise<Message[]> => {
    const response = await api.get(`/conversation/${receiverId}`);
    return response.data.map((msg: any) => ({
      id: msg.id,
      sender_id: msg.sender_id,
      receiver_id: msg.receiver_id,
      content: msg.message,
      timestamp: msg.timestamp,
    }));
  },
  
  sendMessage: async (data: { 
    content: string; 
    reciever_id?: string | number 
  }) => {
    const response = await api.post('/send_messages', {
      content: data.content,
      receiver_id: data.reciever_id ? Number(data.reciever_id) : undefined
    });
    return response.data;
  },
  
  updateMessage: async (messageId: string, content: string) => {
    try {
      const response = await api.put(`/update/${messageId}`, {
        message: content
      });
      return response.data;
    } catch (error) {
      console.error("Failed to update message:", error);
      throw error;
    }
  },
};

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

export default api;
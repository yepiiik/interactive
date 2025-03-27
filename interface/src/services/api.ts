import axios from 'axios';
import {
  AuthResponse,
  CreatePollRequest,
  CreateRoomRequest,
  JoinRoomRequest,
  Poll,
  Room,
  User,
  VoteRequest,
} from '../types';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api';

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add token to requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Auth endpoints
export const auth = {
  register: async (email: string, password: string, name: string): Promise<AuthResponse> => {
    const response = await api.post<AuthResponse>('/auth/register', { email, password, name });
    return response.data;
  },

  login: async (email: string, password: string): Promise<AuthResponse> => {
    const response = await api.post<AuthResponse>('/auth/login', { email, password });
    return response.data;
  },

  googleLogin: () => {
    window.location.href = `${API_URL}/auth/google`;
  },
};

// Room endpoints
export const rooms = {
  create: async (data: CreateRoomRequest): Promise<Room> => {
    const response = await api.post<Room>('/rooms', data);
    return response.data;
  },

  get: async (id: string): Promise<Room> => {
    const response = await api.get<Room>(`/rooms/${id}`);
    return response.data;
  },

  join: async (data: JoinRoomRequest): Promise<Room> => {
    const response = await api.post<Room>('/rooms/join', data);
    return response.data;
  },
};

// Poll endpoints
export const polls = {
  create: async (data: CreatePollRequest): Promise<Poll> => {
    const response = await api.post<Poll>('/polls', data);
    return response.data;
  },

  vote: async (id: number, data: VoteRequest): Promise<void> => {
    await api.post(`/polls/${id}/vote`, data);
  },

  getResults: async (id: number): Promise<{ poll: Poll; results: any[] }> => {
    const response = await api.get(`/polls/${id}/results`);
    return response.data;
  },
};

// WebSocket connection
export const createWebSocket = (roomId: string): WebSocket => {
  const token = localStorage.getItem('token');
  const wsUrl = `${process.env.REACT_APP_WS_URL || 'ws://localhost:8080'}/ws/rooms/${roomId}`;
  return new WebSocket(wsUrl + (token ? `?token=${token}` : ''));
}; 
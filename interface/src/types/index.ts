export interface User {
  id: number;
  email: string;
  name: string;
  googleId?: string;
}

export interface Room {
  id: string;
  name: string;
  hostId: number;
  host: User;
  inviteCode: string;
  isActive: boolean;
  participants: User[];
  createdAt: string;
  updatedAt: string;
}

export interface Poll {
  id: number;
  roomId: string;
  question: string;
  options: Option[];
  duration: number;
  startTime: string;
  endTime: string;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface Option {
  id: number;
  pollId: number;
  text: string;
  isCorrect: boolean;
  votes: Vote[];
  createdAt: string;
  updatedAt: string;
}

export interface Vote {
  id: number;
  userId: number;
  pollId: number;
  optionId: number;
  timeTaken: number;
  createdAt: string;
}

export interface WebSocketMessage {
  type: 'vote' | 'start_poll' | 'end_poll';
  payload: any;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface CreateRoomRequest {
  name: string;
}

export interface JoinRoomRequest {
  inviteCode: string;
}

export interface CreatePollRequest {
  roomId: string;
  question: string;
  options: string[];
  duration: number;
  correctId: number;
}

export interface VoteRequest {
  optionId: number;
  timeTaken: number;
} 
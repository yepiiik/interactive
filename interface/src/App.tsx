import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Route, Routes, Navigate } from 'react-router-dom';
import { LoginForm } from './components/auth/LoginForm';
import { CreateRoomForm } from './components/room/CreateRoomForm';
import { CreatePollForm } from './components/poll/CreatePollForm';
import { PollVote } from './components/poll/PollVote';
import { Room } from './types';
import { rooms, polls, createWebSocket } from './services/api';

const App: React.FC = () => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [currentRoom, setCurrentRoom] = useState<Room | null>(null);
  const [ws, setWs] = useState<WebSocket | null>(null);

  useEffect(() => {
    const token = localStorage.getItem('token');
    if (token) {
      setIsAuthenticated(true);
    }
  }, []);

  useEffect(() => {
    if (currentRoom && isAuthenticated) {
      const websocket = createWebSocket(currentRoom.id);
      websocket.onmessage = (event) => {
        const message = JSON.parse(event.data);
        handleWebSocketMessage(message);
      };
      setWs(websocket);

      return () => {
        websocket.close();
      };
    }
  }, [currentRoom, isAuthenticated]);

  const handleLogin = (token: string) => {
    localStorage.setItem('token', token);
    setIsAuthenticated(true);
  };

  const handleCreateRoom = (room: Room) => {
    setCurrentRoom(room);
  };

  const handleCreatePoll = (poll: any) => {
    // Handle poll creation
  };

  const handleVote = (vote: any) => {
    // Handle vote submission
  };

  const handleWebSocketMessage = (message: any) => {
    switch (message.type) {
      case 'start_poll':
        // Handle poll start
        break;
      case 'vote':
        // Handle vote
        break;
      case 'end_poll':
        // Handle poll end
        break;
    }
  };

  if (!isAuthenticated) {
    return <LoginForm onSuccess={handleLogin} />;
  }

  return (
    <Router>
      <div className="min-h-screen bg-gray-100">
        <nav className="bg-white shadow-sm">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="flex justify-between h-16">
              <div className="flex">
                <div className="flex-shrink-0 flex items-center">
                  <h1 className="text-xl font-bold text-indigo-600">Polling App</h1>
                </div>
              </div>
              <div className="flex items-center">
                <button
                  onClick={() => {
                    localStorage.removeItem('token');
                    setIsAuthenticated(false);
                    setCurrentRoom(null);
                  }}
                  className="ml-4 px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700"
                >
                  Logout
                </button>
              </div>
            </div>
          </div>
        </nav>

        <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
          <Routes>
            <Route
              path="/"
              element={
                currentRoom ? (
                  <Navigate to={`/room/${currentRoom.id}`} />
                ) : (
                  <CreateRoomForm onSuccess={handleCreateRoom} />
                )
              }
            />
            <Route
              path="/room/:id"
              element={
                currentRoom ? (
                  <div>
                    <h2 className="text-2xl font-bold mb-4">{currentRoom.name}</h2>
                    <CreatePollForm roomId={currentRoom.id} onSuccess={handleCreatePoll} />
                    {/* Add PollVote component here when a poll is active */}
                  </div>
                ) : (
                  <Navigate to="/" />
                )
              }
            />
          </Routes>
        </main>
      </div>
    </Router>
  );
};

export default App; 
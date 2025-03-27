import React, { useState } from 'react';
import { rooms } from '../../services/api';
import { Room } from '../../types';

interface CreateRoomFormProps {
  onSuccess: (room: Room) => void;
}

export const CreateRoomForm: React.FC<CreateRoomFormProps> = ({ onSuccess }) => {
  const [name, setName] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const room = await rooms.create({ name });
      onSuccess(room);
    } catch (err) {
      setError('Failed to create room');
    }
  };

  return (
    <div className="max-w-md mx-auto mt-8">
      <form onSubmit={handleSubmit} className="space-y-6">
        {error && (
          <div className="rounded-md bg-red-50 p-4">
            <div className="text-sm text-red-700">{error}</div>
          </div>
        )}
        <div>
          <label htmlFor="room-name" className="block text-sm font-medium text-gray-700">
            Room Name
          </label>
          <div className="mt-1">
            <input
              type="text"
              name="room-name"
              id="room-name"
              required
              className="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
              value={name}
              onChange={(e) => setName(e.target.value)}
            />
          </div>
        </div>

        <div>
          <button
            type="submit"
            className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
          >
            Create Room
          </button>
        </div>
      </form>
    </div>
  );
}; 
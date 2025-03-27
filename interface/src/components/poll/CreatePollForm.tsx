import React, { useState } from 'react';
import { polls } from '../../services/api';
import { Poll } from '../../types';

interface CreatePollFormProps {
  roomId: string;
  onSuccess: (poll: Poll) => void;
}

export const CreatePollForm: React.FC<CreatePollFormProps> = ({ roomId, onSuccess }) => {
  const [question, setQuestion] = useState('');
  const [options, setOptions] = useState(['', '', '', '']);
  const [correctOption, setCorrectOption] = useState(1);
  const [duration, setDuration] = useState(30);
  const [error, setError] = useState('');

  const handleOptionChange = (index: number, value: string) => {
    const newOptions = [...options];
    newOptions[index] = value;
    setOptions(newOptions);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const poll = await polls.create({
        roomId,
        question,
        options: options.filter(opt => opt.trim() !== ''),
        duration,
        correctId: correctOption,
      });
      onSuccess(poll);
    } catch (err) {
      setError('Failed to create poll');
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
          <label htmlFor="question" className="block text-sm font-medium text-gray-700">
            Question
          </label>
          <div className="mt-1">
            <input
              type="text"
              name="question"
              id="question"
              required
              className="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
              value={question}
              onChange={(e) => setQuestion(e.target.value)}
            />
          </div>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700">Options</label>
          <div className="mt-2 space-y-2">
            {options.map((option, index) => (
              <div key={index} className="flex items-center space-x-2">
                <input
                  type="radio"
                  name="correct-option"
                  checked={correctOption === index + 1}
                  onChange={() => setCorrectOption(index + 1)}
                  className="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300"
                />
                <input
                  type="text"
                  className="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                  value={option}
                  onChange={(e) => handleOptionChange(index, e.target.value)}
                  placeholder={`Option ${index + 1}`}
                />
                {index > 1 && (
                  <button
                    type="button"
                    onClick={() => {
                      const newOptions = options.filter((_, i) => i !== index);
                      setOptions(newOptions);
                    }}
                    className="text-red-600 hover:text-red-800"
                  >
                    Remove
                  </button>
                )}
              </div>
            ))}
            {options.length < 4 && (
              <button
                type="button"
                onClick={() => setOptions([...options, ''])}
                className="text-indigo-600 hover:text-indigo-800"
              >
                Add Option
              </button>
            )}
          </div>
        </div>

        <div>
          <label htmlFor="duration" className="block text-sm font-medium text-gray-700">
            Duration (seconds)
          </label>
          <div className="mt-1">
            <input
              type="number"
              name="duration"
              id="duration"
              min="5"
              max="300"
              required
              className="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
              value={duration}
              onChange={(e) => setDuration(parseInt(e.target.value))}
            />
          </div>
        </div>

        <div>
          <button
            type="submit"
            className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
          >
            Create Poll
          </button>
        </div>
      </form>
    </div>
  );
}; 
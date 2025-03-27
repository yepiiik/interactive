import React, { useState, useEffect } from 'react';
import { polls } from '../../services/api';
import { Poll, VoteRequest } from '../../types';

interface PollVoteProps {
  poll: Poll;
  onVote: (vote: VoteRequest) => void;
}

export const PollVote: React.FC<PollVoteProps> = ({ poll, onVote }) => {
  const [selectedOption, setSelectedOption] = useState<number | null>(null);
  const [timeRemaining, setTimeRemaining] = useState(poll.duration);
  const [hasVoted, setHasVoted] = useState(false);

  useEffect(() => {
    const timer = setInterval(() => {
      setTimeRemaining((prev) => {
        if (prev <= 0) {
          clearInterval(timer);
          return 0;
        }
        return prev - 1;
      });
    }, 1000);

    return () => clearInterval(timer);
  }, []);

  const handleVote = async (optionId: number) => {
    if (hasVoted) return;

    const timeTaken = poll.duration - timeRemaining;
    const vote: VoteRequest = {
      optionId,
      timeTaken,
    };

    try {
      await polls.vote(poll.id, vote);
      setHasVoted(true);
      onVote(vote);
    } catch (err) {
      console.error('Failed to submit vote:', err);
    }
  };

  if (!poll.isActive) {
    return (
      <div className="text-center text-gray-600">
        Poll has ended
      </div>
    );
  }

  return (
    <div className="max-w-md mx-auto mt-8">
      <div className="text-center mb-4">
        <h2 className="text-xl font-semibold">{poll.question}</h2>
        <p className="text-gray-600">Time remaining: {timeRemaining}s</p>
      </div>

      <div className="space-y-4">
        {poll.options.map((option) => (
          <button
            key={option.id}
            onClick={() => handleVote(option.id)}
            disabled={hasVoted}
            className={`w-full p-4 rounded-lg text-left transition-colors ${
              selectedOption === option.id
                ? 'bg-indigo-100 border-indigo-500'
                : 'bg-white border-gray-200 hover:bg-gray-50'
            } border-2 ${
              hasVoted ? 'cursor-not-allowed opacity-50' : 'cursor-pointer'
            }`}
          >
            {option.text}
          </button>
        ))}
      </div>

      {hasVoted && (
        <div className="mt-4 text-center text-green-600">
          Vote submitted!
        </div>
      )}
    </div>
  );
}; 
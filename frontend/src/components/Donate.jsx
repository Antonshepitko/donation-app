import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

const Donate = () => {
  const [streamer, setStreamer] = useState('');
  const [donorName, setDonorName] = useState('');
  const [amount, setAmount] = useState('');
  const [currency, setCurrency] = useState('USD');
  const [message, setMessage] = useState('');
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    const res = await fetch('http://45.144.52.58:5000/api/donate', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ streamer, donorName, amount: parseFloat(amount), currency, message }),
    });
    const data = await res.json();
    if (res.status === 201) {
      alert('Donation sent!');
      navigate('/');
    } else {
      alert(data.error || 'Error');
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100">
      <div className="bg-white p-8 rounded-lg shadow-lg w-96">
        <h2 className="text-2xl font-bold mb-6">Donate</h2>
        <form onSubmit={handleSubmit}>
          <input
            type="text"
            placeholder="Streamer Username"
            value={streamer}
            onChange={(e) => setStreamer(e.target.value)}
            className="w-full p-2 mb-4 border rounded"
          />
          <input
            type="text"
            placeholder="Your Name"
            value={donorName}
            onChange={(e) => setDonorName(e.target.value)}
            className="w-full p-2 mb-4 border rounded"
          />
          <input
            type="number"
            placeholder="Amount"
            value={amount}
            onChange={(e) => setAmount(e.target.value)}
            className="w-full p-2 mb-4 border rounded"
          />
          <select
            value={currency}
            onChange={(e) => setCurrency(e.target.value)}
            className="w-full p-2 mb-4 border rounded"
          >
            <option>USD</option>
            <option>EUR</option>
            <option>RUB</option>
          </select>
          <textarea
            placeholder="Message"
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            className="w-full p-2 mb-4 border rounded"
          />
          <button type="submit" className="w-full bg-green-500 text-white p-2 rounded">
            Send Donation
          </button>
        </form>
      </div>
    </div>
  );
};

export default Donate;
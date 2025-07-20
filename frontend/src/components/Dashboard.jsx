import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

const Dashboard = () => {
  const [donations, setDonations] = useState([]);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchDonations = async () => {
      const token = localStorage.getItem('token');
      if (!token) return navigate('/');
      const res = await fetch('http://45.144.52.58:5000/api/donations', {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (res.status === 401) {
        localStorage.removeItem('token');
        navigate('/');
      } else {
        const data = await res.json();
        setDonations(data);
      }
    };
    fetchDonations();
  }, [navigate]);

  return (
    <div className="min-h-screen bg-gray-100 p-8">
      <h1 className="text-3xl font-bold mb-6">Your Donations</h1>
      <div className="space-y-4">
        {donations.map((don) => (
          <div key={don._id} className="bg-white p-4 rounded-lg shadow-md flex justify-between items-start">
            <div>
              <h3 className="font-bold text-lg">{don.donorName}</h3>
              <p className="text-gray-600">{don.message}</p>
            </div>
            <span className="text-green-500 font-bold">{don.amount} {don.currency}</span>
          </div>
        ))}
      </div>
      <button
        onClick={() => { localStorage.removeItem('token'); navigate('/'); }}
        className="mt-4 bg-red-500 text-white p-2 rounded"
      >
        Logout
      </button>
    </div>
  );
};

export default Dashboard;
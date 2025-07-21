import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { motion } from 'framer-motion';

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

  // Функция удаления пользователя (временная)
  const handleDeleteUser = async () => {
    const username = prompt('Enter username to delete:');
    if (!username) return;
    const res = await fetch(`http://45.144.52.58:5000/api/user/${username}`, {
      method: 'DELETE',
    });
    const data = await res.json();
    alert(data.message || data.error);
    localStorage.removeItem('token');
    navigate('/');
  };

  return (
    <div className="min-h-screen bg-gradient-to-r from-gray-100 to-gray-200 p-8">
      <h1 className="text-4xl font-bold mb-8 text-center text-gray-800">Your Donations</h1>
      <div className="space-y-6 max-w-2xl mx-auto">
        {donations.map((don, index) => (
          <motion.div 
            key={don._id}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.3, delay: index * 0.1 }}
            className="bg-white p-6 rounded-xl shadow-lg flex justify-between items-center border-l-4 border-green-500 hover:shadow-xl transition"
          >
            <div>
              <h3 className="font-bold text-xl text-gray-800">{don.donorName}</h3>
              <p className="text-gray-600">{don.message}</p>
            </div>
            <span className="text-green-600 font-bold text-lg">{don.amount} {don.currency}</span>
          </motion.div>
        ))}
      </div>
      <div className="flex justify-center mt-8 space-x-4">
        <motion.button 
          whileHover={{ scale: 1.05 }}
          onClick={() => { localStorage.removeItem('token'); navigate('/'); }}
          className="bg-red-500 text-white p-3 rounded-lg hover:bg-red-600 transition"
        >
          Logout
        </motion.button>
        <motion.button 
          whileHover={{ scale: 1.05 }}
          onClick={handleDeleteUser}
          className="bg-orange-500 text-white p-3 rounded-lg hover:bg-orange-600 transition"
        >
          Delete User (Temp)
        </motion.button>
      </div>
    </div>
  );
};

export default Dashboard;
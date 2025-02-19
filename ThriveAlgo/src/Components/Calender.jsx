import React, { useState, useEffect } from 'react';
import axios from 'axios';
const Calendar = () => {
  const [currentDate, setCurrentDate] = useState(new Date());
  const [holidays, setHolidays] = useState([]);
  const [hoveredDate, setHoveredDate] = useState(null);
  const [showModal, setShowModal] = useState(false);
  const [selectedDate, setSelectedDate] = useState(null);
  const [newHoliday, setNewHoliday] = useState({ name: ''});
  const months = ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"];
  const daysInMonth = new Date(currentDate.getFullYear(), currentDate.getMonth() + 1, 0).getDate();
  const firstDayOfMonth = new Date(currentDate.getFullYear(), currentDate.getMonth(), 1).getDay();
  const API_URL = import.meta.env.VITE_API_URL;


  useEffect(() => {
    fetchHolidays();
  }, []);
  const fetchHolidays = async () => {
    try {
      const response = await axios.get(`${API_URL}/api/holidays`);
      setHolidays(response.data);
      console.log(response.data)
      // return 0: id: '', date: '2025-12-25T00:00:00Z', title: 'Christmas'}
    } catch (error) {
      console.error('Error fetching holidays:', error);
    }
  };
  const updateMonth = (offset) => {
    setCurrentDate(prev => new Date(prev.getFullYear(), prev.getMonth() + offset, 1));
  };
  const handleAddHoliday = (date) => {
    setSelectedDate(date);
    setShowModal(true);
  };
  const handleSubmitHoliday = async () => {
    if (newHoliday.name && selectedDate) {
      try {
        const response = await axios.post(`${API_URL}/api/holidays`, {
          title: newHoliday.name,
          date: selectedDate
        });
        setHolidays([...holidays, response.data]);
        setShowModal(false);
        setNewHoliday({ name: ''});
      } catch (error) {
        console.error('Error adding holiday:', error);
      }
    }
  };
  const handleRemoveHoliday = async (holidayId) => {
    try {
      console.log(holidayId)
      await axios.delete(`${API_URL}/api/holidays/${holidayId}`);
      setHolidays(holidays.filter(h => h.id !== holidayId));
    } catch (error) {
      console.error('Error removing holiday:', error);
    }
  };
  const renderCalendarDays = () => {
    const days = [];
    for (let i = 0; i < firstDayOfMonth; i++) {
      days.push(<div key={`empty-${i}`} className="h-24 border border-gray-200"></div>);
    }
    for (let day = 1; day <= daysInMonth; day++) {
      const date = `${currentDate.getFullYear()}-${String(currentDate.getMonth() + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`;
      const holiday = Array.isArray(holidays) ? holidays.find(h => h.date === date) : null;
      days.push(
        <div
          key={day}
          className="h-24 border border-gray-200 p-2 relative"
          onMouseEnter={() => setHoveredDate(date)}
          onMouseLeave={() => setHoveredDate(null)}
        >
          <div className="font-semibold">{day}</div>
          {holiday && (
            <div className="flex items-center justify-between text-sm text-blue-600 mt-1">
              <span>{holiday.title}</span>
              
              <button onClick={() => handleRemoveHoliday(holiday.id)} className="text-red-500 hover:text-red-700">x</button>
            </div>
          )}
          {hoveredDate === date && !holiday && (
            <button onClick={() => handleAddHoliday(date)} className="absolute bottom-2 right-2 bg-blue-500 text-white px-2 py-1 rounded text-sm">Add Holiday</button>
          )}
        </div>
      );
    }
    return days;
  };

  return (
    <div className="container mx-auto p-4">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-2xl font-bold">{months[currentDate.getMonth()]} {currentDate.getFullYear()}</h2>
        <div className="space-x-2">
          <button onClick={() => updateMonth(-1)} className="px-4 py-2 bg-gray-200 rounded">Previous</button>
          <button onClick={() => updateMonth(1)} className="px-4 py-2 bg-gray-200 rounded">Next</button>
        </div>
      </div>
      <div className="grid grid-cols-7 gap-0">
        {['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'].map(day => (
          <div key={day} className="text-center font-semibold p-2 border border-gray-200">{day}</div>
        ))}
        {renderCalendarDays()}
      </div>
      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center">
          <div className="bg-white p-6 rounded-lg w-96">
            <h3 className="text-lg font-semibold mb-4">Add Holiday</h3>
            <input
              type="text"
              placeholder="Holiday Name"
              className="w-full p-2 border rounded mb-4"
              value={newHoliday.name}
              onChange={(e) => setNewHoliday({ ...newHoliday, name: e.target.value })}
            />
            <div className="flex justify-end space-x-2">
              <button onClick={() => setShowModal(false)} className="px-4 py-2 bg-gray-200 rounded">Cancel</button>
              <button onClick={handleSubmitHoliday} className="px-4 py-2 bg-blue-500 text-white rounded">Save</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Calendar;
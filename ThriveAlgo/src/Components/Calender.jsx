import React, { useState, useEffect, useMemo, useCallback } from 'react';
import axios from 'axios';

const Calendar = () => {
  const [currentDate, setCurrentDate] = useState(new Date());
  const [holidays, setHolidays] = useState([]);
  const [hoveredDate, setHoveredDate] = useState(null);
  const [showModal, setShowModal] = useState(false);
  const [selectedDate, setSelectedDate] = useState(null);
  const [newHoliday, setNewHoliday] = useState({ name: '' });
  const [error, setError] = useState(null);
  const API_URL = 'http://localhost:8080';

  const months = useMemo(() => ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"], []);

  // Memoize these calculations to prevent unnecessary recalculations
  const { daysInMonth, firstDayOfMonth } = useMemo(() => ({
    daysInMonth: new Date(currentDate.getFullYear(), currentDate.getMonth() + 1, 0).getDate(),
    firstDayOfMonth: new Date(currentDate.getFullYear(), currentDate.getMonth(), 1).getDay()
  }), [currentDate]);

  const fetchHolidays = useCallback(async () => {
    setError(null);
    try {
      const response = await axios.get(`${API_URL}/api/holidays`);
      setHolidays(response.data);
    } catch (error) {
      console.error('Error fetching holidays:', error);
      setError('Failed to load holidays. Please try again later.');
    }
  }, []);

  useEffect(() => {
    fetchHolidays();
  }, [fetchHolidays]);

  const updateMonth = useCallback((offset) => {
    setCurrentDate(prev => new Date(prev.getFullYear(), prev.getMonth() + offset, 1));
  }, []);

  const handleAddHoliday = useCallback((date) => {
    setSelectedDate(date);
    setShowModal(true);
  }, []);

  // Helper function to standardize date formats for comparison
  const formatDateForComparison = useCallback((dateStr) => {
    if (!dateStr) return null;
    // Handle ISO strings (from API) by extracting YYYY-MM-DD portion
    if (dateStr.includes('T')) {
      return dateStr.split('T')[0];
    }
    return dateStr;
  }, []);

  const handleSubmitHoliday = useCallback(async () => {
    if (newHoliday.name && selectedDate) {
      setError(null);
      try {
        const response = await axios.post(`${API_URL}/api/holidays`, {
          title: newHoliday.name,
          date: selectedDate
        });
        setHolidays(prev => [...prev, response.data]);
        setShowModal(false);
        setNewHoliday({ name: '' });
      } catch (error) {
        console.error('Error adding holiday:', error);
        setError('Failed to add holiday. Please try again.');
      }
    }
  }, [newHoliday.name, selectedDate]);

  const handleRemoveHoliday = useCallback(async (holidayId, e) => {
    e.stopPropagation(); // Prevent triggering parent events
    setError(null);
    try {
      await axios.delete(`${API_URL}/api/holidays/${holidayId}`);
      setHolidays(prev => prev.filter(h => h.id !== holidayId));
    } catch (error) {
      console.error('Error removing holiday:', error);
      setError('Failed to remove holiday. Please try again.');
    }
  }, []);

  const renderCalendarDays = useMemo(() => {
    const days = [];
    
    // Add empty cells for days before the first day of the month
    for (let i = 0; i < firstDayOfMonth; i++) {
      days.push(<div key={`empty-${i}`} className="h-24 border border-gray-200" aria-hidden="true"></div>);
    }

    // Generate cells for each day of the month
    for (let day = 1; day <= daysInMonth; day++) {
      const date = `${currentDate.getFullYear()}-${String(currentDate.getMonth() + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`;
      
      // Find holiday for this date, accounting for different date formats
      const holiday = Array.isArray(holidays) ? 
        holidays.find(h => formatDateForComparison(h.date) === date) : null;

      days.push(
        <div
          key={day}
          className="h-24 border border-gray-200 p-2 relative"
          onMouseEnter={() => setHoveredDate(date)}
          onMouseLeave={() => setHoveredDate(null)}
          role="gridcell"
          tabIndex="0"
          aria-label={`${months[currentDate.getMonth()]} ${day}, ${currentDate.getFullYear()}${holiday ? `, Holiday: ${holiday.title}` : ''}`}
        >
          <div className="font-semibold">{day}</div>
          {holiday && (
            <div className="flex items-center justify-between text-sm text-blue-600 mt-1">
              <span>{holiday.title}</span>
              <button 
                onClick={(e) => handleRemoveHoliday(holiday.id, e)} 
                className="text-red-500 hover:text-red-700"
                aria-label={`Remove holiday: ${holiday.title}`}
              >
                ×
              </button>
            </div>
          )}
          {hoveredDate === date && !holiday && (
            <button 
              onClick={() => handleAddHoliday(date)} 
              className="absolute bottom-2 right-2 bg-blue-500 text-white px-2 py-1 rounded text-sm"
              aria-label={`Add holiday on ${months[currentDate.getMonth()]} ${day}`}
            >
              Add Holiday
            </button>
          )}
        </div>
      );
    }
    return days;
  }, [currentDate, daysInMonth, firstDayOfMonth, holidays, hoveredDate, months, handleAddHoliday, handleRemoveHoliday, formatDateForComparison]);

  return (
    <div className="container mx-auto p-4">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-2xl font-bold">{months[currentDate.getMonth()]} {currentDate.getFullYear()}</h2>
        <div className="space-x-2">
          <button 
            onClick={() => updateMonth(-1)} 
            className="px-4 py-2 bg-gray-200 rounded"
            aria-label="Previous month"
          >
            Previous
          </button>
          <button 
            onClick={() => updateMonth(1)} 
            className="px-4 py-2 bg-gray-200 rounded"
            aria-label="Next month"
          >
            Next
          </button>
        </div>
      </div>

      {error && (
        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
          {error}
        </div>
      )}

      <div className="grid grid-cols-7 gap-0" role="grid" aria-label="Calendar">
        {['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'].map(day => (
          <div key={day} className="text-center font-semibold p-2 border border-gray-200" role="columnheader">{day}</div>
        ))}
        {renderCalendarDays}
      </div>

      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" role="dialog" aria-modal="true" aria-labelledby="modal-title">
          <div className="bg-white p-6 rounded-lg w-96">
            <h3 id="modal-title" className="text-lg font-semibold mb-4">Add Holiday</h3>
            <input
              type="text"
              placeholder="Holiday Name"
              className="w-full p-2 border rounded mb-4"
              value={newHoliday.name}
              onChange={(e) => setNewHoliday({ ...newHoliday, name: e.target.value })}
              aria-label="Holiday name"
            />
            <div className="flex justify-end space-x-2">
              <button 
                onClick={() => setShowModal(false)} 
                className="px-4 py-2 bg-gray-200 rounded"
              >
                Cancel
              </button>
              <button 
                onClick={handleSubmitHoliday} 
                className="px-4 py-2 bg-blue-500 text-white rounded"
              >
                Save
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Calendar;
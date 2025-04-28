import React, { useState, useEffect } from 'react';

const CacheWrapper = ({ cacheKey, fetchData, children }) => {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const cachedData = localStorage.getItem(cacheKey);
    if (cachedData) {
      setData(JSON.parse(cachedData));
      setLoading(false);
    } else {
      fetchData().then(fetchedData => {
        localStorage.setItem(cacheKey, JSON.stringify(fetchedData));
        setData(fetchedData);
        setLoading(false);
      });
    }
  }, [cacheKey, fetchData]);

  if (loading) {
    return <div>Loading...</div>;
  }

  return children(data);
};

export default CacheWrapper;

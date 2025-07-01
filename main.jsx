import React from 'react';
import { createRoot } from 'react-dom/client';
import Gallery from './Gallery.jsx';
import Profile from './Profile.jsx';

function App() {
  return (
    <div>
      <Profile />
      <Gallery />
    </div>
  );
}

const root = createRoot(document.getElementById('root'));
root.render(<App />);

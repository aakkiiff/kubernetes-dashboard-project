import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';

import Sidebar from './components/Sidebar';
import Home from './pages/Home';
import Namespace from './pages/Namespace';
import Pod from './pages/Pod';
import Node from './pages/Node';
import './App.css';
import 'bootstrap/dist/css/bootstrap.min.css';
function App() {
  return (
    <Router>
      <div className="App">
  
        <Sidebar />
        <div className="content">
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/node" element={<Node />} />
            <Route path="/namespace" element={<Namespace />} />
            <Route path="/pod" element={<Pod />} />

          </Routes>
        </div>
      </div>
    </Router>
  );
}

export default App;

import React from 'react';
import { Link } from 'react-router-dom';
import './Sidebar.css';

function Sidebar() {
  return (
    <div className="sidebar">
      <h2 className="sidebar-title">
        <Link to="/" style={{ textDecoration: 'none', color: '#2c3e50' }}>
          K8's Dashboard
        </Link>
      </h2>
      <ul className="sidebar-links">
        <li><Link to="/node">Node</Link></li>
        <li><Link to="/namespace">Namespace</Link></li>
        <li><Link to="/pod">Pod</Link></li>
      </ul>
    </div>
  );
}

export default Sidebar;

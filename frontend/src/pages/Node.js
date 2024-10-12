import React, { useState, useEffect } from 'react';
import { Button, Container, Table, Form } from 'react-bootstrap';

function NodeList() {
  const [nodes, setNodes] = useState([]);
  const [error, setError] = useState('');
  const [refreshInterval, setRefreshInterval] = useState(5000); // Default to 5 seconds

  // Function to fetch nodes from the API
  const fetchNodes = async () => {
    try {
      const response = await fetch('/api/node');
      if (!response.ok) {
        throw new Error('Network response was not ok');
      }
      const data = await response.json();
      setNodes(data.nodes); // Assuming response is an array of nodes
      setError(''); // Clear any previous errors
    } catch (err) {
      setError('Failed to fetch nodes. Please try again.');
      console.error('Fetch error:', err);
    }
  };

  // Fetch nodes when the component mounts and set interval for refreshing
  useEffect(() => {
    fetchNodes();

    // Set up interval to fetch nodes based on selected refresh interval
    const intervalId = setInterval(fetchNodes, refreshInterval);

    // Clear the interval on component unmount
    return () => clearInterval(intervalId);
  }, [refreshInterval]); // Depend on refreshInterval

  return (
    <Container className="mt-4">
      <h1 className="mb-4">Node List</h1>

      {/* Dropdown for selecting refresh interval */}
      <Form.Group controlId="refreshInterval" className="mb-4">
        <Form.Label>Select Refresh Interval</Form.Label>
        <Form.Control
          as="select"
          value={refreshInterval / 1000} // Convert to seconds for display
          onChange={(e) => setRefreshInterval(e.target.value * 1000)} // Convert back to milliseconds
        >
          <option value={1}>1 second</option>
          <option value={2}>2 seconds</option>
          <option value={5}>5 seconds</option>
          <option value={10}>10 seconds</option>
          <option value={30}>30 seconds</option>
        </Form.Control>
      </Form.Group>

      {error && <p style={{ color: 'red' }}>{error}</p>}

      {nodes.length > 0 ? (
        <Table striped bordered hover>
          <thead>
            <tr>
              <th>#</th>
              <th>Node Name</th>
            </tr>
          </thead>
          <tbody>
            {nodes.map((node, index) => (
              <tr key={index}>
                <td>{index + 1}</td>
                <td>{node}</td>
              </tr>
            ))}
          </tbody>
        </Table>
      ) : (
        <p>No nodes found.</p>
      )}
    </Container>
  );
}

export default NodeList;

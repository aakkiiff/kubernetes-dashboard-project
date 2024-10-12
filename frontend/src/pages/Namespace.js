import React, { useState, useEffect } from "react";
import { Button, Container, Table, Modal, Form } from "react-bootstrap";

function NamespaceList() {
  const [namespaces, setNamespaces] = useState([]);
  const [error, setError] = useState("");
  const [newNamespace, setNewNamespace] = useState("");
  const [showModal, setShowModal] = useState(false);
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [namespaceToDelete, setNamespaceToDelete] = useState("");
  const [refreshInterval, setRefreshInterval] = useState(5000); // Default to 5 seconds
  const [updatePodName, setUpdatePodName] = useState(""); // State for the pod name being updated
  const [updatePodImage, setUpdatePodImage] = useState(""); // State for the image being updated
  const [showUpdateModal, setShowUpdateModal] = useState(false); // State to control the update modal visibility
  
  // Function to fetch namespaces from the API
  const fetchNamespaces = async () => {
    try {
      const response = await fetch("/api/namespace");
      if (!response.ok) {
        throw new Error("Network response was not ok");
      }
      const data = await response.json();
      setNamespaces(data.namespaces); // Assuming response is an array of namespaces
      setError(""); // Clear any previous errors
    } catch (err) {
      setError("Failed to fetch namespaces. Please try again.");
      console.error("Fetch error:", err);
    }
  };

  // Function to create a new namespace
  const createNamespace = async (e) => {
    e.preventDefault();
    try {
      const response = await fetch("/api/namespace", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ name: newNamespace }),
      });

      if (!response.ok) {
        throw new Error("Network response was not ok");
      }

      // Refresh the namespace list after creation
      fetchNamespaces();
      setNewNamespace(""); // Clear the input
      setShowModal(false); // Close the modal
    } catch (err) {
      setError("Failed to create namespace. Please try again.");
      console.error("Create namespace error:", err);
    }
  };

  // Function to delete a namespace
  const deleteNamespace = async () => {
    try {
      const response = await fetch("/api/namespace", {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ name: namespaceToDelete }),
      });

      if (!response.ok) {
        throw new Error("Network response was not ok");
      }

      // Refresh the namespace list after deletion
      fetchNamespaces();
      setNamespaceToDelete(""); // Clear the selected namespace
      setShowDeleteModal(false); // Close the delete confirmation modal
    } catch (err) {
      setError("Failed to delete namespace. Please try again.");
      console.error("Delete namespace error:", err);
    }
  };

  // Fetch namespaces when the component mounts and set interval for refreshing
  useEffect(() => {
    fetchNamespaces();

    // Set up interval to fetch namespaces based on selected refresh interval
    const intervalId = setInterval(fetchNamespaces, refreshInterval);

    // Clear the interval on component unmount
    return () => clearInterval(intervalId);
  }, [refreshInterval]); // Depend on refreshInterval

  return (
    <Container className="mt-4">
      <h1 className="mb-4">Namespace List</h1>

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

      <Button
        variant="success"
        onClick={() => setShowModal(true)}
        className="mb-4"
        style={{ padding: "10px 20px" }} // Add padding to the button
      >
        Create New Namespace
      </Button>
      {error && <p style={{ color: "red" }}>{error}</p>}

      {namespaces.length > 0 ? (
        <Table striped bordered hover>
          <thead>
            <tr>
              <th>#</th>
              <th>Namespace Name</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {namespaces.map((namespace, index) => (
              <tr key={index}>
                <td>{index + 1}</td>
                <td>{namespace}</td>
                <td>
                  <Button
                    variant="danger"
                    onClick={() => {
                      setNamespaceToDelete(namespace);
                      setShowDeleteModal(true);
                    }}
                  >
                    Delete
                  </Button>
                </td>
              </tr>
            ))}
          </tbody>
        </Table>
      ) : (
        <p>No namespaces found.</p>
      )}

      {/* Modal for creating new namespace */}
      <Modal show={showModal} onHide={() => setShowModal(false)}>
        <Modal.Header closeButton>
          <Modal.Title>Create New Namespace</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <Form onSubmit={createNamespace}>
            <Form.Group controlId="namespaceName">
              <Form.Label>Namespace Name</Form.Label>
              <Form.Control
                type="text"
                placeholder="Enter namespace name"
                value={newNamespace}
                onChange={(e) => setNewNamespace(e.target.value)}
                required
              />
            </Form.Group>
            <Button variant="primary" type="submit">
              Create
            </Button>
          </Form>
        </Modal.Body>
      </Modal>

      {/* Modal for confirming deletion */}
      <Modal show={showDeleteModal} onHide={() => setShowDeleteModal(false)}>
        <Modal.Header closeButton>
          <Modal.Title>Confirm Deletion</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <p>
            Are you sure you want to delete the namespace "{namespaceToDelete}"?
          </p>
        </Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={() => setShowDeleteModal(false)}>
            Cancel
          </Button>
          <Button variant="danger" onClick={deleteNamespace}>
            Delete
          </Button>
        </Modal.Footer>
      </Modal>
    </Container>
  );
}

export default NamespaceList;

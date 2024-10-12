import React, { useState, useEffect, useCallback } from "react";
import { Button, Container, Table, Form, Row, Col, Modal } from "react-bootstrap";

function PodPage() {
  const [namespaces, setNamespaces] = useState([]); 
  const [selectedNamespace, setSelectedNamespace] = useState(""); 
  const [pods, setPods] = useState([]); 
  const [refreshInterval, setRefreshInterval] = useState(5000); 
  const [error, setError] = useState(""); 
  const [podName, setPodName] = useState(""); 
  const [podImage, setPodImage] = useState(""); 
  const [showModal, setShowModal] = useState(false); 
  
  const fetchNamespaces = async () => {
    try {
      const response = await fetch("/api/namespace");
      if (!response.ok) {
        throw new Error("Failed to fetch namespaces.");
      }
      const data = await response.json();
      setNamespaces(data.namespaces || []); 
    } catch (err) {
      console.error("Error fetching namespaces:", err);
      setError("Error fetching namespaces.");
    }
  };

  const fetchPods = useCallback(async () => {
    if (!selectedNamespace) return;
    try {
      const response = await fetch(`/api/pods/${selectedNamespace}`);
      if (!response.ok) {
        throw new Error("Failed to fetch pods.");
      }
      const data = await response.json();
      setPods(data.pods || []); 
    } catch (err) {
      console.error("Error fetching pods:", err);
      setError("Error fetching pods.");
    }
  }, [selectedNamespace]);

  const handleCreatePod = async (e) => {
    e.preventDefault();
    try {
      const response = await fetch("/api/pod", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          name: podName,
          namespacename: selectedNamespace,
          image: podImage,
        }),
      });

      if (!response.ok) {
        throw new Error("Failed to create pod.");
      }

      setPodName("");
      setPodImage("");
      setShowModal(false);
      fetchPods();
    } catch (err) {
      console.error("Error creating pod:", err);
      setError("Error creating pod.");
    }
  };

  const handleDeletePod = async (podName) => {
    try {
      const response = await fetch("/api/pod", {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          name: podName,
          namespacename: selectedNamespace,
        }),
      });

      if (!response.ok) {
        throw new Error("Failed to delete pod.");
      }

      fetchPods();
    } catch (err) {
      console.error("Error deleting pod:", err);
      setError("Error deleting pod.");
    }
  };

  useEffect(() => {
    fetchNamespaces();
  }, []);

  useEffect(() => {
    if (selectedNamespace) {
      fetchPods(); 
      const intervalId = setInterval(fetchPods, refreshInterval); 
      return () => clearInterval(intervalId); 
    }
  }, [fetchPods, refreshInterval]);

  return (
    <Container className="mt-4">
      <h1 className="mb-4">Pods in Namespace</h1>

      {/* First Line: Refresh interval and Namespace dropdown */}
      <Row className="mb-4">
        <Col md={6}>
          <Form.Group controlId="refreshInterval">
            <Form.Label>Refresh Interval</Form.Label>
            <Form.Control
              as="select"
              value={refreshInterval / 1000}
              onChange={(e) => setRefreshInterval(e.target.value * 1000)}
            >
              <option value={1}>1 second</option>
              <option value={5}>5 seconds</option>
              <option value={10}>10 seconds</option>
              <option value={30}>30 seconds</option>
            </Form.Control>
          </Form.Group>
        </Col>
        <Col md={6}>
          <Form.Group controlId="namespaceSelect">
            <Form.Label>Select Namespace</Form.Label>
            <Form.Control
              as="select"
              value={selectedNamespace}
              onChange={(e) => setSelectedNamespace(e.target.value)}
            >
              <option value="">-- Select Namespace --</option>
              {namespaces.map((namespace, index) => (
                <option key={index} value={namespace}>
                  {namespace}
                </option>
              ))}
            </Form.Control>
          </Form.Group>
        </Col>
      </Row>

      {/* Second Line: Create Pod Button */}
      <Row className="mb-4">
        <Col>
          <Button 
            variant="success" 
            disabled={!selectedNamespace}
            onClick={() => setShowModal(true)} 
          >
            Create Pod
          </Button>
        </Col>
      </Row>

      {/* Third Section: List of Pods */}
      {error && <p style={{ color: "red" }}>{error}</p>}
      {pods.length > 0 ? (
        <Table striped bordered hover>
          <thead>
            <tr>
              <th>#</th>
              <th>Pod Name</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {pods.map((pod, index) => (
              <tr key={index}>
                <td>{index + 1}</td>
                <td>{pod}</td>
                <td>
                  <Button 
                    variant="danger" 
                    onClick={() => handleDeletePod(pod)} 
                    className="ml-2"
                  >
                    Delete
                  </Button>
                </td>
              </tr>
            ))}
          </tbody>
        </Table>
      ) : (
        <p>No pods found for the selected namespace.</p>
      )}

      {/* Modal for Creating Pod */}
      <Modal show={showModal} onHide={() => setShowModal(false)}>
        <Modal.Header closeButton>
          <Modal.Title>Create a New Pod</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <Form onSubmit={handleCreatePod}>
            <Form.Group controlId="podName">
              <Form.Label>Pod Name</Form.Label>
              <Form.Control
                type="text"
                placeholder="Enter pod name"
                value={podName}
                onChange={(e) => setPodName(e.target.value)}
                required
              />
            </Form.Group>
            <Form.Group controlId="podImage">
              <Form.Label>Pod Image</Form.Label>
              <Form.Control
                type="text"
                placeholder="Enter image (e.g., user/image:tag)"
                value={podImage}
                onChange={(e) => setPodImage(e.target.value)}
                required
              />
            </Form.Group>
            <Form.Group controlId="namespaceName">
              <Form.Label>Namespace</Form.Label>
              <Form.Control
                type="text"
                value={selectedNamespace}
                disabled
                readOnly
              />
            </Form.Group>
            <Button variant="primary" type="submit">
              Create Pod
            </Button>
          </Form>
        </Modal.Body>
      </Modal>

    </Container>
  );
}

export default PodPage;

import React, { useEffect, useState } from 'react'
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Card from 'react-bootstrap/Card'
import './style.css'
import { Alert, Badge, Button, Container } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Link } from 'react-router-dom';

const Modules = () => {
  const [modules, setModules] = useState<any>([])

  useEffect(() => {
    fetch("http://localhost:8080/modules")
      .then(res => res.json())
      .then(res => setModules(res))

  }, [])

  return (
    <div className='container mt-4'>
      <div style={{display: 'flex', justifyContent: 'space-between'}}>
        <h1 className='text-white'>Modules</h1>
        <Button variant='secondary' style={{background: '#373739'}}>
          <FontAwesomeIcon icon='plus' />{' '}New module
        </Button>

      </div>
      <hr style={{color: '#fff'}} />
      {modules.map((module: any) => (
        <Col key={module.name} className="mb-3">
          <Card style={{background: 'transparent', border: '2px solid #1c1c1e', color: '#fff'}}>
            <Card.Header style={{background: '#1c1c1e'}}>
              <FontAwesomeIcon icon="folder" />  {module.name}
            </Card.Header>
            <Card.Body>
              <div className='font-weight-light'>
                <div><strong>Repository: </strong> {module.repositoryPath}</div>
                <div><strong>Deployment path: </strong> {module.deploymentPath}</div>
                <div><strong>Secret: </strong> {module.secretRef}</div>
                <div><strong>Template: </strong> {module.templateType}</div>
              </div>
            </Card.Body>
          </Card>
        </Col>
      ))}
    </div>
  )
}

export default Modules
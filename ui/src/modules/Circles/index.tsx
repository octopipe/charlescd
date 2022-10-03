import React, { useEffect, useState } from 'react'
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Card from 'react-bootstrap/Card'
import './style.css'
import { Alert, Badge, Button, Container } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Link } from 'react-router-dom';
import CircleModules from '../CircleModules';

const colors = {
  "": "secondary",
  'Healthy': 'success',
  'Progressing': 'primary',
  'Degraded': 'danger'
} as any

const Circles = () => {
  const [circles, setCircles] = useState<any>([])
  const [currentCircle, setCurrentCircle] = useState<any>({})
  const [showAddModule, setShowAddModule] = useState(false)

  const getCircleStatusByModules = (modules: any) => {
    const dangerModules = Object.keys(modules)
      .filter(module => modules[module].health !== "Healthy")
    
    return dangerModules.length <= 0 ? "Healthy" : modules[dangerModules[0]]["health"]
  }

  useEffect(() => {
    fetch("http://localhost:8080/circles")
      .then(res => res.json())
      .then(res => setCircles(res))

  }, [])

  return (
    <div className='container mt-4'>
      <div style={{display: 'flex', justifyContent: 'space-between', }}>
        <h1 className='text-white'>Circles</h1>
        <Button variant='secondary' style={{background: '#373739'}}>
          <FontAwesomeIcon icon='plus' />{' '}New circle
        </Button>

      </div>
      <hr style={{color: '#fff'}} />
      <Row xs={3} md={3} className="g-4">
        {circles.map((circle: any) => (
          <Col key={circle.name}>
            <Card style={{background: '#1c1c1e', color: '#fff'}}>
              <Card.Body>
                <p>
                  <div className='mb-2' style={{display: 'flex', justifyContent: 'space-between', alignItems: 'center'}}>
                    <h4><FontAwesomeIcon icon={["far", 'circle']} /> <Link className='text-decoration-none text-white' to={`./${circle.name}`}>{circle.name}</Link></h4>
                    <div style={{display: 'flex'}}>
                      {circle?.isDefault && (
                        <Badge bg="primary">Default</Badge>
                      )}

                      {getCircleStatusByModules(circle?.modules || {}) !== 'Healthy' && (
                        <Badge bg="danger">{getCircleStatusByModules(circle?.modules)}</Badge>
                      )}
                    </div>
                  </div>
                  <p>
                    {circle?.description}
                  </p>
                </p>
                <CircleModules circle={circle} />
              </Card.Body>
            </Card>
          </Col>
        ))}
      </Row>
    </div>
  )
}


export default Circles
import React, { useEffect, useState } from 'react'
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Card from 'react-bootstrap/Card'
import './style.css'
import { Alert, Badge, Button, Container } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Link } from 'react-router-dom';

const fakeCircles = [
  {
    name: 'circle-default',
    isDefault: true,
    description: '',
    status: 'Healthy',
    modules: [
      {
        name: 'guestbook-ui',
        status: 'Healthy',
        error: ""
      },
      {
        name: 'guestbook-nginx',
        status: 'Healthy',
        error: ""
      },
      {
        name: 'guestbook-server',
        status: 'Healthy',
        error: ""
      },
    ]
  },
  {
    name: 'circle-2',
    description: '',
    status: 'Degraded',
    modules: [
      {
        name: 'guestbook-server',
        status: 'Degraded',
        error: 'Back-off pulling image \"naaaaginx:1.16.1\"'
      },
    ]
  },
  {
    name: 'circle-sample',
    description: '',
    status: 'Healthy',
    modules: [
      {
        name: 'guestbook-ui',
        status: 'Healthy',
        error: ""
      },
      {
        name: 'guestbook-nginx',
        status: 'Healthy',
        error: ""
      }
    ]
  },
]

const colors = {
  "": "secondary",
  'Healthy': 'success',
  'Progressing': 'primary',
  'Degraded': 'danger'
} as any

const Circles = () => {
  const [circles, setCircles] = useState<any>([])

  useEffect(() => {
    fetch("http://localhost:8080/circles")
      .then(res => res.json())
      .then(res => setCircles(res))

  }, [])

  return (
    <div className='my-4'>
      <div style={{display: 'flex', justifyContent: 'space-between'}}>
        <h1 className='text-white'>Circles</h1>
        <Button variant='secondary' style={{background: '#373739'}}>
          <FontAwesomeIcon icon='plus' />{' '}New circle
        </Button>

      </div>
      <hr style={{color: '#fff'}} />
      <Row xs={3} md={3} className="g-4">
        {circles.map((circle: any) => (
          <Col key={circle.name}>
            <Card style={{background: '#2c2c2e', color: '#fff'}}>
              <Card.Body>
                <p>
                  <div className='mb-2' style={{display: 'flex', justifyContent: 'space-between'}}>
                    <Link to={`./${circle.name}/diagram`}>{circle.name}</Link>
                    <div style={{display: 'flex'}}>
                      {circle?.isDefault && (
                        <Badge bg="primary">Default</Badge>
                      )}

                      {circle.status !== 'Healthy' && (
                        <Badge bg="danger">Danger</Badge>
                      )}
                    </div>
                  </div>
                  <p>
                    {circle?.description}
                  </p>
                </p>
                {Object.keys(circle.modules).map((name: any) => (
                  <Card
                    bg={colors[circle.modules[name].status]}
                    className='mb-2' 
                    style={{background: 'transparent'}}
                  >
                    <Card.Body>
                      {name}
                    </Card.Body>
                  </Card>
                ))}
                <div className="d-grid gap-2">
                  <Button className='mt-2' variant='secondary' style={{background: '#373739'}}>
                    <FontAwesomeIcon icon='plus' />{' '}Add module
                  </Button>
                </div>
              </Card.Body>
            </Card>
          </Col>
        ))}
      </Row>
    </div>
  )
}


export default Circles
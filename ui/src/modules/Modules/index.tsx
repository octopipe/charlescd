import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import React from 'react'
import { Button, Card, Col, Row } from 'react-bootstrap'
import './style.css'

const Modules = () => {
  return (
    <div className='my-4'>
    <div style={{display: 'flex', justifyContent: 'space-between'}}>
      <h1 className='text-white'>Modules</h1>
      <Button variant='secondary' style={{background: '#373739'}}>
        <FontAwesomeIcon icon='plus' />{' '}New module 
      </Button>
    </div>
    <hr style={{color: '#fff'}} />
    <Row xs={3} md={4} className="g-4">
      {Array.from({ length: 4 }).map((_, idx) => (
        <Col>
          <Card style={{background: '#2c2c2e', color: '#fff'}}>
            <Card.Header>Module name</Card.Header>
            <Card.Body>
            </Card.Body>
          </Card>
        </Col>
      ))}
    </Row>
  </div>
  )
}

export default Modules
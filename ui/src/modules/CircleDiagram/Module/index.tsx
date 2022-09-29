import React, { useEffect, useState } from "react";
import { Card, Dropdown } from "react-bootstrap";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import './style.css'

const colors = {
  "": "secondary",
  'Healthy': 'success',
  'Progressing': 'primary',
  'Degraded': 'danger'
} as any

const Module = ({ name, status, onUpdated, circle }: any) => {

  const handleRemoveModule = () => {
    const modules = circle.modules.filter((mod: any) => mod.moduleRef !== name)

    circle.modules = modules
    
    fetch("http://localhost:8080/circles/" + circle?.name, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(circle),
    })
      .then(res => res.json())
      .then(res => onUpdated(res))
  }

  return (
    <Card
      bg={colors[status]}
      className='mb-2 circle_module' 
      style={{background: 'transparent'}}
    >
      <Card.Body style={{display: 'flex', justifyContent: 'space-between'}}>
        {name}
        <Dropdown>
          <Dropdown.Toggle className="circle_module_toggle" variant="success" id="dropdown-basic">
            <FontAwesomeIcon icon="ellipsis-vertical" />
          </Dropdown.Toggle>
          <Dropdown.Menu>
            <Dropdown.Item onClick={handleRemoveModule}>Remove</Dropdown.Item>
            <Dropdown.Item>Move To</Dropdown.Item>
          </Dropdown.Menu>
        </Dropdown>
      </Card.Body>
    </Card>
  )
}

export default Module
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

const Module = ({ module, name, onRemove, onEdit }: any) => {
  return (
    <Card
      bg={colors[module?.status]}
      className='mb-2 circle_module' 
      style={{background: 'transparent'}}
    >
      <Card.Body style={{display: 'flex', justifyContent: 'space-between'}}>
        {name}
        <Dropdown>
          <Dropdown.Toggle className="circle_module_toggle" id="dropdown-basic">
            <FontAwesomeIcon icon="ellipsis-vertical" />
          </Dropdown.Toggle>
          <Dropdown.Menu>
            <Dropdown.Item onClick={() => onRemove(name)}>Remove</Dropdown.Item>
            <Dropdown.Item>Move To</Dropdown.Item>
            <Dropdown.Item onClick={() => onEdit(name)}>Edit</Dropdown.Item>
          </Dropdown.Menu>
        </Dropdown>
      </Card.Body>
    </Card>
  )
}

export default Module
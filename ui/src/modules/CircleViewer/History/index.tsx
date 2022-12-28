import React from 'react'
import { ListGroup } from 'react-bootstrap'
import { CircleModel } from '../../../core/types/circle'
import './style.scss'

interface Props {
  circle?: CircleModel 
}


const CircleHistory = ({ circle }: Props) => {
  return (
    <>
      <ListGroup>
        {circle?.status.history?.map(history => (
          <ListGroup.Item>{history.message}</ListGroup.Item>
        ))}
      </ListGroup>
    </>
  )
}

export default CircleHistory
import React, { useCallback, useEffect, useState } from 'react'
import { Card, Col, Container, Row } from 'react-bootstrap'
import { useSelector } from 'react-redux'
import useFetch from 'use-http'
import { useAppSelector } from '../../core/hooks/redux'
import CircleModules from '../CircleModules'
import { ReactComponent as EmptyCircles } from '../../core/assets/svg/empty-circles.svg'
import './style.scss'
import { CircleItem } from './types'
import Placeholder from '../../core/components/Placeholder'
import { Link } from 'react-router-dom'


const Circles = () => {
  const currentWorkspace = useAppSelector(state => state.main.currentWorkspace)
  const [circles, setCircles] = useState<CircleItem[]>([])
  const { response, get } = useFetch()

  const loadCircles = async () => {
    const circles = await get(`/workspaces/${currentWorkspace}/circles`)
    if (response.ok) setCircles(circles || [])
  }

  useEffect(() => {
    if (currentWorkspace == "")
      return
    loadCircles()
  }, [currentWorkspace])

  return (
    <Container className='circles'>
      {circles.length <= 0 && (
        <Placeholder text="There are no circles in this workspace">
          <EmptyCircles />
        </Placeholder>
      )}
      <Row>
      
      {circles?.map(circle => (
        <Col xs={3}>
          <Card className='circles__item'>
            <Card.Body>
              <Card.Title><Link to={circle.name}>{circle.name}</Link></Card.Title>
              <CircleModules modules={circle.modules} />
            </Card.Body>
          </Card>
        </Col>
      ))}
      </Row>
    </Container>
  )
}

export default Circles
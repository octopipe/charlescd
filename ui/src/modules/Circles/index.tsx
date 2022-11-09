import React, { useCallback, useEffect, useState } from 'react'
import { Button, Card, Col, Container, Form, FormControl, Row } from 'react-bootstrap'
import { useSelector } from 'react-redux'
import useFetch from 'use-http'
import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { useAppSelector } from '../../core/hooks/redux'
import CircleModules from '../CircleModules'
import { ReactComponent as EmptyCircles } from '../../core/assets/svg/empty-circles.svg'
import './style.scss'
import { CircleItem } from './types'
import Placeholder from '../../core/components/Placeholder'
import { Link, Navigate, useNavigate, useParams } from 'react-router-dom'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

const data = [
  {
    name: 'Page A',
    uv: 4000,
    pv: 2400,
    amt: 2400,
  },
  {
    name: 'Page B',
    uv: 3000,
    pv: 1398,
    amt: 2210,
  },
  {
    name: 'Page C',
    uv: 2000,
    pv: 9800,
    amt: 2290,
  },
  {
    name: 'Page D',
    uv: 2780,
    pv: 3908,
    amt: 2000,
  },
  {
    name: 'Page E',
    uv: 1890,
    pv: 4800,
    amt: 2181,
  },
  {
    name: 'Page F',
    uv: 2390,
    pv: 3800,
    amt: 2500,
  },
  {
    name: 'Page G',
    uv: 3490,
    pv: 4300,
    amt: 2100,
  },
];

const Circles = () => {
  const navigate = useNavigate()
  const { workspaceId } = useParams()
  const [circles, setCircles] = useState<CircleItem[]>([])
  const { response, get } = useFetch()

  const loadCircles = async () => {
    
    const circles = await get(`/workspaces/${workspaceId}/circles`)
    if (response.ok) setCircles(circles || [])
  }

  useEffect(() => {
    loadCircles()
  }, [workspaceId])

  return (
    <div className='circles'>
      <div className='circles__navbar'>
        <div>
          <Form.Group className="mb-3" controlId="formBasicEmail">
            <Form.Control type="text" placeholder='Search circle...' />
          </Form.Group>
        </div>
        <div className='circles__navbar__controller'>
          <Button className='mx-2' onClick={() => navigate(`/workspaces/${workspaceId}/circles/create`)}>
            <FontAwesomeIcon icon="add"  />
          </Button>
          <Button>
            <FontAwesomeIcon icon="rotate" />
          </Button>
        </div>
      </div>

      <Container className='circles__content'>
        {circles.length <= 0 && (
          <Placeholder text="There are no circles in this workspace">
            <EmptyCircles />
          </Placeholder>
        )}
        <Row>
        {circles?.map(circle => (
          <Col xs={4} key={circle.name}>
            <Card className='circles__content__item mb-4'>
              <Card.Body>
                <Card.Title className='d-flex justify-content-between'>
                  <Link className='text-decoration-none text-white' to={circle.name}>
                    {circle.name}
                  </Link>
                  <div className='d-flex'>
                    <FontAwesomeIcon icon="eye" size='sm' className='me-3' />
                    <FontAwesomeIcon icon="trash" size='sm' />
                  </div>
                </Card.Title>
                <CircleModules modules={circle.modules} />
              </Card.Body>
              <div style={{ width: '100%', height: 300 }}>
                <ResponsiveContainer>
                  <AreaChart width={300} height={250} data={data}
                    margin={{ top: 10, right: 0, left: 0, bottom: 0 }}>
                    <defs>
                      <linearGradient id="colorUv" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="5%" stopColor="#8884d8" stopOpacity={0.8}/>
                        <stop offset="95%" stopColor="#8884d8" stopOpacity={0}/>
                      </linearGradient>
                      <linearGradient id="colorPv" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="5%" stopColor="#82ca9d" stopOpacity={0.8}/>
                        <stop offset="95%" stopColor="#82ca9d" stopOpacity={0}/>
                      </linearGradient>
                    </defs>

                    <Tooltip />
                    <Area type="monotone" dataKey="uv" stroke="#8884d8" fillOpacity={1} fill="url(#colorUv)" />
                    <Area type="monotone" dataKey="pv" stroke="#82ca9d" fillOpacity={1} fill="url(#colorPv)" />
                  </AreaChart>
                </ResponsiveContainer>
              </div>
            </Card>
          </Col>
        ))}
        </Row>
      </Container>
    </div>
  )
}

export default Circles
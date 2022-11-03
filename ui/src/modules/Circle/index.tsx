import React, { useCallback, useEffect, useState } from 'react'
import { Card, Col, Container, Row } from 'react-bootstrap'
import { useSelector } from 'react-redux'
import useFetch from 'use-http'
import { useAppSelector } from '../../core/hooks/redux'
import CircleModules from '../CircleModules'
import { ReactComponent as EmptyCircles } from '../../core/assets/svg/empty-circles.svg'
import './style.scss'
import Placeholder from '../../core/components/Placeholder'
import { Outlet, useParams } from 'react-router-dom'
import CircleSidebar from './Sidebar'
import { Circle, CircleModule } from './types'
import { CircleItemModule } from '../CircleModules/types'


const CircleMain = () => {
  const { workspaceId, circleName } = useParams()
  const [circle, setCircle] = useState<Circle>()
  const { response, get } = useFetch()

  const loadCircle = async () => {
    const circle = await get(`/workspaces/${workspaceId}/circles/${circleName}`)
    if (response.ok) setCircle(circle)
  }

  useEffect(() => {
    loadCircle()
  }, [workspaceId])  

  return (
    <div className='circle'>
      {circle && <CircleSidebar circle={circle} /> }
      <div className='circle__content'>
        <Outlet />
      </div>
    </div>
  )
}

export default CircleMain
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
import { Outlet, useParams } from 'react-router-dom'


const Circle = () => {
  const currentWorkspace = useAppSelector(state => state.main.currentWorkspace)
  const { name: circleName } = useParams()
  const [circle, setCircle] = useState<CircleItem[]>()
  const { response, get } = useFetch()

  const loadCircle = async () => {
    const circle = await get(`/workspaces/${currentWorkspace}/circles/${circleName}`)
    if (response.ok) setCircle(circle)
  }

  useEffect(() => {
    if (currentWorkspace == "")
      return

    loadCircle()
  }, [currentWorkspace])

  return (
    <div className='circle'>
      Circle
      <Outlet />
    </div>
  )
}

export default Circle
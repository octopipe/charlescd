import React, { useCallback, useEffect, useState } from 'react'
import useFetch from 'use-http'
import './style.scss'
import { CircleItem, CirclePagination } from './types'
import Placeholder from '../../core/components/Placeholder'
import { Link, Navigate, useLocation, useNavigate, useParams, useSearchParams } from 'react-router-dom'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import CirclesSidebar from './Sidebar';
import Circle from './Circle';

const CirclesMain = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const [searchParams, setSearchParams] = useSearchParams();
  const { workspaceId } = useParams()
  const [circles, setCircles] = useState<CirclePagination>({continue: '', items: []})
  const [activeCircleIds, setActiveCirclesIds] = useState<string[]>([])
  const { response, get } = useFetch()

  const loadCircles = async () => {
    
    const circles = await get(`/workspaces/${workspaceId}/circles`)
    if (response.ok) setCircles(circles || [])
  }

  useEffect(() => {
    loadCircles()
  }, [workspaceId])

  useEffect(() => {
    let currentActiveCirclesIds: string[] = []
    const keys = searchParams.forEach((value, key) => currentActiveCirclesIds.push(key))
    setActiveCirclesIds(currentActiveCirclesIds)
  }, [location])


  const handleCircleClick = (circleId: string) => {
    setSearchParams(i => {
      if (!i.has(circleId)) {
        i.append(circleId, "1")
      } else {
        i.delete(circleId)
      }

      return i
    })
  }

  return (
    <div className='circles'>
      <CirclesSidebar circles={circles} onCircleClick={handleCircleClick} />
      <div className='circles__content'>
        {activeCircleIds.map(id => (
          <Circle circleId={id} />
        ))}
     
      </div>
    </div>
  )
}

export default CirclesMain
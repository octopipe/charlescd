import React, { useCallback, useEffect, useState } from 'react'
import useFetch, { CachePolicies } from 'use-http'
import './style.scss'
import { CircleItem, CirclePagination } from './types'
import Placeholder from '../../core/components/Placeholder'
import { Link, Navigate, useLocation, useNavigate, useParams, useSearchParams } from 'react-router-dom'
import CirclesSidebar from './Sidebar';
import Circle from './Circle';
import { Circle as CircleType } from './Circle/types'

const createCircleId = 'untitled'

const CirclesMain = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const [searchParams, setSearchParams] = useSearchParams();
  const { workspaceId } = useParams()
  const [circles, setCircles] = useState<CirclePagination>({continue: '', items: []})
  const [activeCircleIds, setActiveCirclesIds] = useState<string[]>([])
  const { response, get, post, delete: deleteMethod } = useFetch({cachePolicy: CachePolicies.NO_CACHE})

  const loadCircles = async () => {
    const circles = await get(`/workspaces/${workspaceId}/circles`)
    if (response.ok) setCircles(circles || [])
  }

  useEffect(() => {
    loadCircles()
  }, [workspaceId])

  useEffect(() => {
    let currentActiveCirclesIds: string[] = []
    searchParams.forEach((value, key) => currentActiveCirclesIds.push(key))
    setActiveCirclesIds(currentActiveCirclesIds)
  }, [location])


  const handleCircleClick = (circleId: string) => {
    setSearchParams(i => {
      if (!i.has(circleId)) {
        i.append(circleId, "R")
      } else {
        i.delete(circleId)
      }

      return i
    })
  }

  const handleCircleCreateClick = () => {
    
    setSearchParams(i => {
      if (!i.has(createCircleId)) {
        i.append(createCircleId, "C")
      }

      return i
    })
  }

  const handleCloseCircle = (circleId: string) => {
    setSearchParams(i => {
      i.delete(circleId)
      return i
    })
  }

  const handleUpdateCircle = (circleId: string) => {
    setSearchParams(i => {
      i.set(circleId, "U")
      return i
    })
  }

  const handleDeleteCircle = async (circleId: string) => {
    await deleteMethod(`/workspaces/${workspaceId}/circles/${circleId}`)
    await loadCircles()
    if (response.ok) {
      setSearchParams(i => {
        i.delete(circleId)
        return i
      })
    }
  }

  const handleSaveCircle = async (circle: CircleType) => {
    const newCircle = await post(`/workspaces/${workspaceId}/circles`, circle)
    if (response.ok) {
      setSearchParams(i => {
        if (!i.has(newCircle.id)) {
          i.append(newCircle.id, "R")
        }

        i.delete(createCircleId)
  
        return i
      })
      await loadCircles()
    }
  }

  return (
    <div className='circles'>
      <CirclesSidebar
        circles={circles}
        onCircleClick={handleCircleClick}
        onCircleCreateClick={handleCircleCreateClick}
      />
      <div className='circles__content'>
        {activeCircleIds.map(id => (
          <Circle
            key={id}
            circleId={id}
            circleOp={searchParams.get(id) || 'R'}
            onClose={handleCloseCircle}
            onUpdate={handleUpdateCircle}
            onSave={handleSaveCircle}
            onDelete={handleDeleteCircle}
          />
        ))}
      </div>
    </div>
  )
}

export default CirclesMain
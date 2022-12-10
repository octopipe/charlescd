import React, { useCallback, useEffect, useState } from 'react'
import { CirclePagination } from './types'
import { useLocation, useParams, useSearchParams } from 'react-router-dom'
import CirclesSidebar from './Sidebar';
import Circle from './Circle';
import { Circle as CircleType } from './Circle/types'
import Placeholder from '../../core/components/Placeholder'
import { ReactComponent as EmptyCirclesSVG } from '../../core/assets/svg/empty-circles.svg'
import { setBreadcrumbItems } from '../Main/mainSlice'
import './style.scss'
import { useAppDispatch } from '../../core/hooks/redux';
import useFetch from '../../core/hooks/fetch';

const createCircleId = 'untitled'

const CirclesMain = () => {
  const dispatch = useAppDispatch()
  const { workspaceId } = useParams()
  const [circles, setCircles] = useState<CirclePagination>()
  const [searchParams, setSearchParams] = useSearchParams();
  const [activeCircleIds, setActiveCirclesIds] = useState<string[]>([])
  const { loading: loadingCircles, fetch: get } = useFetch<CirclePagination>()
  const { fetch } = useFetch()

  useEffect(() => {
    get(`/workspaces/${workspaceId}/circles`, 'GET').then(res => setCircles(res))
    dispatch(setBreadcrumbItems([
      { name: 'Circles', to: `/workspaces/${workspaceId}/circles` },
    ]))
  }, [])

  useEffect(() => {
    let currentActiveCirclesIds: string[] = []
    searchParams.forEach((value, key) => currentActiveCirclesIds.push(key))
    setActiveCirclesIds(currentActiveCirclesIds)
  }, [searchParams])


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
    await fetch(`/workspaces/${workspaceId}/circles/${circleId}`, 'DELETE')
    await get(`/workspaces/${workspaceId}/circles`, 'GET')
  }

  const handleSaveCircle = async (circle: CircleType) => {
    const newCircle = await fetch(`/workspaces/${workspaceId}/circles`, 'POST', circle)
    await get(`/workspaces/${workspaceId}/circles`, 'GET')
    setSearchParams(i => {
      if (!i.has(newCircle.id)) {
        i.append(newCircle.id, "R")
      }

      i.delete(createCircleId)

      return i
    })
  }

  return (
    <div className='circles'>
      <CirclesSidebar
        circles={circles}
        loading={loadingCircles}
        onCircleClick={handleCircleClick}
        onCircleCreateClick={handleCircleCreateClick}
      />
      <div className={activeCircleIds.length > 0 ? 'circles__content' : 'circles__content-empty'}>
        {activeCircleIds.length <= 0 && (
          <div className='container'>
            <Placeholder text='No circle selected'>
              <EmptyCirclesSVG />
            </Placeholder>
          </div>
        )}
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
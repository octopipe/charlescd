import React, { useCallback, useEffect, useState } from 'react'
import { CirclePagination, CIRCLE_ROUTING_STRATEGY, CIRCLE_VIEW_MODE } from '../../core/types/circle'
import { useLocation, useParams, useSearchParams } from 'react-router-dom'
import Placeholder from '../../core/components/Placeholder'
import { ReactComponent as EmptyCirclesSVG } from '../../core/assets/svg/empty-circles.svg'
import { setBreadcrumbItems } from '../Main/mainSlice'
import './style.scss'
import { useAppDispatch, useAppSelector } from '../../core/hooks/redux';
import useFetch from '../../core/hooks/fetch';
import AppSidebar from '../../core/components/AppSidebar';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { fetchCircles } from './circlesSlice';
import { FETCH_STATUS } from '../../core/utils/fetch'
import Viewer from '../../core/components/Viewer'
import CircleViewer from '../CircleViewer'
import { ProgressBar } from 'react-bootstrap'

const createCircleId = 'untitled'

const COLOR = ['#7e57c2', '#29b6f6', '#66bb6a', '#ff9800', '#03a9f4']

const CirclesMain = () => {
  const dispatch = useAppDispatch()
  const { workspace } = useAppSelector(state => state.main)
  const { list, status } = useAppSelector(state => state.circles)
  const { workspaceId } = useParams()
  const [searchParams, setSearchParams] = useSearchParams();
  const [activeCircleIds, setActiveCirclesIds] = useState<string[]>([])

  useEffect(() => {
    dispatch(fetchCircles(workspaceId))
  }, [])

  useEffect(() => {
    let currentActiveCirclesIds: string[] = []
    searchParams.forEach((value, key) => currentActiveCirclesIds.push(key))
    setActiveCirclesIds(currentActiveCirclesIds)
  }, [searchParams])


  const handleCircleClick = (circleId: string) => {
    setSearchParams(i => {
      if (!i.has(circleId)) {
        i.append(circleId, CIRCLE_VIEW_MODE.VIEW)
      } else {
        i.delete(circleId)
      }

      return i
    })
  }

  const handleCircleCreateClick = () => {
    setSearchParams(i => {
      if (!i.has(createCircleId)) {
        i.append(createCircleId, CIRCLE_VIEW_MODE.CREATE)
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
      i.set(circleId, CIRCLE_VIEW_MODE.UPDATE)
      return i
    })
  }

  const handleDeleteCircle = async (circleId: string) => {
    await fetch(`/workspaces/${workspaceId}/circles/${circleId}`, { method: 'DELETE' })
    dispatch(fetchCircles(workspaceId))
  }

  const handleChangeViewMode = (circleId: string, viewMode: CIRCLE_VIEW_MODE) => {
    setSearchParams(i => {
      i.set(circleId, viewMode)
      return i
    })
  }

  return (
    <div className='circles'>
      <AppSidebar>
        <AppSidebar.Header>
          <AppSidebar.HeaderItem onClick={handleCircleCreateClick}>
            <FontAwesomeIcon icon="plus-circle" className="me-1" /> Create circle
          </AppSidebar.HeaderItem>
        </AppSidebar.Header>
        <AppSidebar.List loading={status === FETCH_STATUS.LOADING}>
          {list.items && list.items.length > 0 && list?.items.map(item => (
            <AppSidebar.ListItem
              key={item.id}
              isActive={searchParams.has(item.id)}
              icon={['far', 'circle']}
              activeIcon="circle"
              text={item.name}
              onClick={() => handleCircleClick(item.id)}
            />
          ))}
        </AppSidebar.List>
        {workspace?.routingStrategy === CIRCLE_ROUTING_STRATEGY.CANARY && (
          <div className='circles__progress'>
            <ProgressBar>
              {list?.items?.map((item, idx) => (
                <ProgressBar variant="000" now={item.routing.canary?.weight} key={1}
                  style={{backgroundColor: COLOR[idx]}}
                />
              ))}
            </ProgressBar>
          </div>
        )}
      </AppSidebar>
      <div className={activeCircleIds.length > 0 ? 'circles__content' : 'circles__content-empty'}>
        {activeCircleIds.length <= 0 && (
          <div className='container'>
            <Placeholder text='No circle selected'>
              <EmptyCirclesSVG />
            </Placeholder>
          </div>
        )}
        {activeCircleIds.map(id => (
          <CircleViewer
            key={id}
            circleId={id}
            viewMode={searchParams.get(id) as CIRCLE_VIEW_MODE}
            onClose={handleCloseCircle}
            onChangeViewMode={handleChangeViewMode}
          />
        ))}
      </div>
    </div>
  )
}

export default CirclesMain
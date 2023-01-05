import React, { useEffect } from 'react'
import { useParams, useSearchParams } from 'react-router-dom'
import { circleApi } from '../../core/api/circle'
import Spinner from '../../core/components/Spinner'
import Viewer, { ViewerTabsOption } from '../../core/components/Viewer'
import { useAppDispatch, useAppSelector } from '../../core/hooks/redux'
import { Circle, CIRCLE_VIEW_MODE } from '../../core/types/circle'
import { FETCH_STATUS } from '../../core/utils/fetch'
import { fetchCircles } from '../Circles/circlesSlice'
import CircleTree from '../CircleTree'
import { fetchCircle, fetchCircleCreate, fetchCircleSync, fetchCircleUpdate, removeCircleViewer } from './circleViewerSlice'
import CircleForm from './Form'
import CircleHistory from './History'
import CircleMetrics from './Metrics'

interface Props {
  circleId: string
  viewMode: CIRCLE_VIEW_MODE
  onClose: (circleId: string) => void
  onChangeViewMode: (circleId: string, viewMode: CIRCLE_VIEW_MODE) => void
}



const CircleViewer = ({ circleId, viewMode, onClose, onChangeViewMode }: Props) => {
  const dispatch = useAppDispatch()
  const { workspaceId } = useParams()
  const { circleViewer } = useAppSelector(state => state)
  const [searchParams, setSearchParams] = useSearchParams();

  const OPTIONS: ViewerTabsOption[] = [
    { icon: 'trash', name: 'Remove', onAction: () => { } },
    { icon: 'rotate',  name: 'Sync', onAction: () => { dispatch(fetchCircleSync({ workspaceId, circleId })) } }
  ]

  useEffect(() => {
    if (viewMode !== CIRCLE_VIEW_MODE.CREATE)
      dispatch(fetchCircle({workspaceId, circleId}))
  }, [])

  const handleClose = () => {
    onClose(circleId)
    dispatch(removeCircleViewer({ circleId }))
  }

  const handleSave = async (form: Circle) => {
    await dispatch(fetchCircleCreate({workspaceId, data: form}))
    dispatch(fetchCircles(workspaceId))
    onClose(circleId)
  }

  const handleUpdate = async (form: Circle) => {
    await dispatch(fetchCircleUpdate({workspaceId, circleId, data: form}))
    dispatch(fetchCircles(workspaceId))
    dispatch(fetchCircle({workspaceId, circleId}))
  }

  if (circleViewer[circleId]) {
    return (
      <Viewer>
        {circleViewer[circleId]?.item.status === FETCH_STATUS.LOADING && <Spinner />}
        { circleViewer[circleId]?.item.status === FETCH_STATUS.SUCCEEDED && (
          <>
            <Viewer.Tabs hasOptions={viewMode !== CIRCLE_VIEW_MODE.CREATE} options={OPTIONS}>
              <Viewer.TabsItem
                id={circleId}
                icon="circle"
                text={circleViewer[circleId]?.item.data.name}
                isActive={viewMode === CIRCLE_VIEW_MODE.VIEW}
                hasClose={true}
                onClick={() => onChangeViewMode(circleId, CIRCLE_VIEW_MODE.VIEW)}
                onClose={handleClose}
              />
              <Viewer.TabsItem
                id={circleId}
                icon="diagram-project"
                text="Tree"
                isActive={viewMode === CIRCLE_VIEW_MODE.TREE}
                onClick={() => onChangeViewMode(circleId, CIRCLE_VIEW_MODE.TREE)}
              />
             <Viewer.TabsItem
                id={circleId}
                icon="chart-simple"
                text="Metrics"
                isActive={viewMode === CIRCLE_VIEW_MODE.METRICS}
                onClick={() => onChangeViewMode(circleId, CIRCLE_VIEW_MODE.METRICS)}
              /> 
            </Viewer.Tabs>
            <Viewer.Content>
              {viewMode === CIRCLE_VIEW_MODE.VIEW && circleViewer[circleId]?.item && <CircleForm circle={circleViewer[circleId]?.item.data} viewMode={viewMode} onSave={handleSave} onUpdate={handleUpdate}  />}
              {viewMode === CIRCLE_VIEW_MODE.TREE && circleViewer[circleId]?.item && <CircleTree circleId={circleId}  />}
              {viewMode === CIRCLE_VIEW_MODE.METRICS && circleViewer[circleId]?.item && <CircleMetrics circle={circleViewer[circleId]?.item.data}  />}
            </Viewer.Content>
          </>
        )}
      </Viewer>
    )
  }

  return (
    <Viewer>
      <>
        <Viewer.Tabs>
          <Viewer.TabsItem
            id={circleId}
            icon="circle"
            text={'untitled'}
            isActive={true}
            hasClose={true}
            onClick={() => {}}
            onClose={handleClose}
          />
        </Viewer.Tabs>
        <Viewer.Content>
          {viewMode !== CIRCLE_VIEW_MODE.TREE && <CircleForm circle={circleViewer[circleId]?.item.data} viewMode={viewMode} onSave={handleSave} onUpdate={handleUpdate}  />}
        </Viewer.Content>
      </>
    </Viewer>
  )
}

export default CircleViewer
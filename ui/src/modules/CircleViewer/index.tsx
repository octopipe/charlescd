import React, { useEffect } from 'react'
import { useParams } from 'react-router-dom'
import Spinner from '../../core/components/Spinner'
import Viewer from '../../core/components/Viewer'
import { useAppDispatch, useAppSelector } from '../../core/hooks/redux'
import { CIRCLE_VIEW_MODE } from '../../core/types/circle'
import { FETCH_STATUS } from '../../core/utils/fetch'
import CircleTree from '../CircleTree'
import { fetchCircle, removeCircleViewer } from './circleViewerSlice'
import CircleForm from './Form'

interface Props {
  circleId: string
  viewMode: CIRCLE_VIEW_MODE
  onClose: (circleId: string) => void
  onChangeViewMode: (circleId: string, viewMode: CIRCLE_VIEW_MODE) => void
}

const CircleViewer = ({ circleId, viewMode, onClose, onChangeViewMode }: Props) => {
  const { workspaceId } = useParams()
  const dispatch = useAppDispatch()
  const { circleViewer } = useAppSelector(state => state)

  useEffect(() => {
    if (viewMode !== CIRCLE_VIEW_MODE.CREATE)
      dispatch(fetchCircle({workspaceId, circleId}))
  }, [])

  const handleClose = () => {
    onClose(circleId)
    dispatch(removeCircleViewer({ circleId }))

  }

  const isCreate = () => viewMode === CIRCLE_VIEW_MODE.CREATE

  if (circleViewer[circleId]) {
    return (
      <Viewer>
        {circleViewer[circleId]?.status === FETCH_STATUS.LOADING && <Spinner />}
        { circleViewer[circleId]?.status === FETCH_STATUS.SUCCEEDED && (
          <>
            <Viewer.Tabs hasOptions={viewMode !== CIRCLE_VIEW_MODE.CREATE} options={[{ name: 'Remove', onAction: () => {} }]}>
              <Viewer.TabsItem
                id={circleId}
                icon="circle"
                text={circleViewer[circleId]?.item.name}
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
            </Viewer.Tabs>
            <Viewer.Content>
              {viewMode === CIRCLE_VIEW_MODE.VIEW && circleViewer[circleId]?.item && <CircleForm circle={circleViewer[circleId]?.item} viewMode={viewMode} onSave={() => {}}  />}
              {viewMode === CIRCLE_VIEW_MODE.TREE && circleViewer[circleId]?.item && <CircleTree circleId={circleId}  />}
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
          {viewMode !== CIRCLE_VIEW_MODE.TREE && <CircleForm circle={circleViewer[circleId]?.item} viewMode={viewMode} onSave={() => {}}  />}
        </Viewer.Content>
      </>
    </Viewer>
  )
}

export default CircleViewer
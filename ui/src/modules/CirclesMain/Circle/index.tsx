import React, { useCallback, useEffect, useState } from 'react'
import './style.scss'
import { Circle as CircleType, CircleEnrivonment, CircleModel, CircleRouting, CircleRoutingCustomMatch, CircleRoutingSegment } from './types'
import { useParams, useSearchParams } from 'react-router-dom'
import { useAppSelector } from '../../../core/hooks/redux'
import Alert from '../../../core/components/Alert'
import CircleContent from './Content'
import CircleTree from './Tree'
import CircleTabs from './Tabs'
import DynamicContainer from '../../../core/components/DynamicContainer'
import useFetch from '../../../core/hooks/fetch'

interface Props {
  circleId: string
  circleOp: string
  onClose: (id: string) => void
  onUpdate: (id: string) => void
  onSave: (circle: CircleType) => void
  onDelete: (circleId: string) => void
}

enum TABS {
  CONTENT = 'content',
  TREE = 'tree'
}

const Circle = ({ circleId, circleOp, onClose, onSave, onDelete }: Props) => {
  const { workspaceId } = useParams()
  const [circle, setCircle] = useState<CircleModel>()
  const { loading, fetch } = useFetch()
  const [ showDeleteAlert, toggleDeleteAlert ] = useState(false)
  const [activeTab, setActiveTab] = useState(TABS.CONTENT)


  useEffect(() => {
    if (circleOp === "C") {
     return
    }

    fetch(`/workspaces/${workspaceId}/circles/${circleId}`, 'GET').then(res => setCircle(res))
  }, [])

  const handleDelete = (circleId: string) => {
    onDelete(circleId)
    toggleDeleteAlert(false)
  }

  return (
    <DynamicContainer loading={loading} className='circle'>
      <CircleTabs
        circleId={circleId}
        activeTab={activeTab}
        circleOp={circleOp}
        circle={circle}
        onClose={id => onClose(id)}
        onChange={tab => setActiveTab(tab)}
        onDelete={handleDelete}
      />
      {activeTab === TABS.CONTENT && <CircleContent circle={circle} circleOp={circleOp} onSave={onSave} />}
      {activeTab === TABS.TREE && <CircleTree circleId={circleId} /> }
      <Alert show={showDeleteAlert} action={() => handleDelete(circleId)} onClose={() => toggleDeleteAlert(false)} />
    </DynamicContainer>
  )
}

export default Circle
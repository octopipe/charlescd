import React, { useCallback, useEffect, useState } from 'react'
import useFetch from 'use-http'
import './style.scss'
import { Circle as CircleType, CircleEnrivonment, CircleModel, CircleRouting, CircleRoutingCustomMatch, CircleRoutingSegment } from './types'
import { useParams, useSearchParams } from 'react-router-dom'

import "ace-builds/src-noconflict/mode-json";
import "ace-builds/src-noconflict/theme-monokai";
import { useAppSelector } from '../../../core/hooks/redux'
import Alert from '../../../core/components/Alert'
import CircleContent from './Content'
import CircleTree from './Tree'
import CircleTabs from './Tabs'

interface Props {
  circleId: string
  circleOp: string
  onClose: (id: string) => void
  onUpdate: (id: string) => void
  onSave: (circle: CircleType) => void
  onDelete: (circleId: string) => void
}

const initialEnviroments = [
  { key: 'KEY_EXAMPLE', value: 'VALUE_EXAMPLE' }
]

const initialCustomMatch = { headers: { 'x-header-example': '1111' } }

const initialSegments = [
  { key: 'email', op: 'EQUAL', value: 'email@mail.com' }
]

enum TABS {
  CONTENT = 'content',
  TREE = 'tree'
}

const Circle = ({ circleId, circleOp, onClose, onSave, onDelete }: Props) => {
  const [searchParams] = useSearchParams();
  const { routingStrategy } = useAppSelector(state => state.main)
  const { workspaceId } = useParams()
  const [circle, setCircle] = useState<CircleModel>()
  const { response, get } = useFetch()
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [matchStrategy, setMatchStrategy] = useState('customMatch')
  const [customMatch, setCustomMatch] = useState<CircleRoutingCustomMatch>(initialCustomMatch)
  const [segments, setSegments] = useState<CircleRoutingSegment[]>(initialSegments)
  const [environments, setEnvironments] = useState<CircleEnrivonment[]>(initialEnviroments)
  const [ showDeleteAlert, toggleDeleteAlert ] = useState(false)
  const [activeTab, setActiveTab] = useState(TABS.CONTENT)
  const [modules, setModules] = useState([])

  const loadCircle = async () => {
    const circle = await get(`/workspaces/${workspaceId}/circles/${circleId}`)
    if (response.ok) setCircle(circle || [])
  }

  useEffect(() => {
    if (circleOp !== "C")
      loadCircle()
  }, [workspaceId])

  useEffect(() => {
    setName(circle?.name || '')
    setDescription(circle?.description || '')
    setCustomMatch(circle?.routing?.match?.customMatch || initialCustomMatch)
    setEnvironments(circle?.environments || initialEnviroments)
  }, [circle])

  const handleClickSave = () => {
    let routing: CircleRouting = { strategy: routingStrategy }
    if (routingStrategy === 'match') {
      if (matchStrategy === 'customMatch') {
        routing = {
          ...routing,
          match: { customMatch} ,
        }
      } else {
        routing = {
          ...routing,
          match: { segments }
        }
      }
    }

    const newCircle = {
      name,
      description,
      environments,
      modules: [],
      routing,
    }

    console.log(newCircle)

    onSave(newCircle)
  }

  const handleDelete = (circleId: string) => {
    onDelete(circleId)
    toggleDeleteAlert(false)
  }

  const CustomToggle = React.forwardRef<any, any>(({ children, onClick }, ref) => (
    <a
      ref={ref}
      onClick={(e) => {
        e.preventDefault();
        onClick(e);
      }}
      className="circle-modules__item__menu"
    >
      {children}
    </a>
  ));

  return (
    <div className='circle'>
      <CircleTabs
        circleId={circleId}
        activeTab={activeTab}
        circleOp={circleOp}
        circle={circle}
        onClose={onClose}
        onChange={tab => setActiveTab(tab)}
        onDelete={handleDelete}
      />
      {activeTab === TABS.CONTENT && <CircleContent circleId={circleId} circleOp={circleOp} onSave={handleClickSave} />}
      {activeTab === TABS.TREE && <CircleTree show={true} circleId={circleId} onClose={() => setActiveTab(TABS.CONTENT)} /> }
      <Alert show={showDeleteAlert} action={() => handleDelete(circleId)} onClose={() => toggleDeleteAlert(false)} />
    </div>
  )
}

export default Circle
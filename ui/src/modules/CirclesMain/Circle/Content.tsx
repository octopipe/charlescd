import React, { useEffect, useState } from 'react'
import { Nav } from 'react-bootstrap'
import useFetch from 'use-http'
import './style.scss'
import { Circle as CircleType, CircleEnrivonment, CircleModel, CircleRouting, CircleRoutingCustomMatch, CircleRoutingSegment } from './types'
import { useParams, useSearchParams } from 'react-router-dom'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import CircleModules from '../../CircleModules'
import AceEditor from "react-ace";

import "ace-builds/src-noconflict/mode-json";
import "ace-builds/src-noconflict/theme-monokai";
import ViewInput from '../../../core/components/ViewInput'
import { useAppSelector } from '../../../core/hooks/redux'

interface Props {
  circleId: string
  circleOp: string
  onUpdate?: (id: string) => void
  onSave: (circle: CircleType) => void
}

const initialEnviroments = [
  { key: 'KEY_EXAMPLE', value: 'VALUE_EXAMPLE' }
]

const initialCustomMatch = { headers: { 'x-header-example': '1111' } }

const initialSegments = [
  { key: 'email', op: 'EQUAL', value: 'email@mail.com' }
]

const CircleContent = ({ circleId, circleOp, onSave }: Props) => {
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

  const isCreate = () => {
    return searchParams.get(circleId) === 'C'
  }

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


  return (
    <div className='circle__content'>
      {circleOp === 'C' && (
        <div className='circle__save'>
          <div onClick={handleClickSave}>
            <FontAwesomeIcon icon="check" color='#4caf50' className="me-1" /> Save circle
          </div>
        </div>
      )}
      <div className='circle__content'>
        <div className='circle__content__title'>
          <FontAwesomeIcon icon={["far", "circle"]} className="me-2" />
          <ViewInput
            value={name}
            edit={isCreate()}
            canEdit={isCreate()}
            onChange={setName}
            placeholder="Circle name"
          />
        </div>
        <div className='circle__content__description'>
          <FontAwesomeIcon icon="align-justify" className="me-2" />
          <ViewInput
            value={description}
            edit={isCreate()}
            onChange={setDescription}
            as="textarea"
            placeholder="Circle description"
          />
        </div>
        <div className='circle__content__section'>
          <div className='circle__content__section__title'>
            <FontAwesomeIcon icon="route" className="me-2" /> Routing
          </div>
          <Nav variant='pills' defaultActiveKey={matchStrategy} onSelect={key => setMatchStrategy(key || 'customMatch')} className='mb-3'>
            <Nav.Item><Nav.Link eventKey="customMatch">Custom match</Nav.Link></Nav.Item>
            <Nav.Item><Nav.Link eventKey="segments">Segments</Nav.Link></Nav.Item>
          </Nav>
          {matchStrategy === 'customMatch' ? (
            <div className='circle__content__section__custom-match'>
              <AceEditor
                width='100%'
                height='200px'
                fontSize={14}
                mode="json"
                theme="monokai"
                value={JSON.stringify(customMatch, null, 2)}
                onChange={value => setCustomMatch(JSON.parse(value))}
              />
            </div>
          ) : (
            <div className='circle__content__section__segments'>
              <AceEditor
                width='100%'
                height='200px'
                fontSize={14}
                mode="json"
                theme="monokai"
                value={JSON.stringify(segments, null, 2)}
                onChange={value => setSegments(JSON.parse(value))}
              />
            </div>
          )}
          
          
        </div>
        <div className='circle__content__section'>
          <div className='circle__content__section__title'>
            <FontAwesomeIcon icon="folder" className="me-2" /> Modules
          </div>
          {circle && <CircleModules circle={circle} />}
        </div>
        <div className='circle__content__section'>
          <div className='circle__content__section__title'>
            <FontAwesomeIcon icon="folder" className="me-2" /> Environments
          </div>
          <AceEditor
            width='100%'
            height='200px'
            fontSize={14}
            mode="json"
            theme="monokai"
            value={JSON.stringify(environments, null, 2)}
            onChange={value => setEnvironments(JSON.parse(value))}
          />
        </div>
      </div>
    </div>
  )
}

export default CircleContent
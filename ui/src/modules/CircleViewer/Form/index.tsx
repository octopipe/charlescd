import React, { useCallback, useEffect, useState } from 'react'
import { Nav } from 'react-bootstrap'
import './style.scss'
import { Circle as CircleType, CircleEnrivonment, CircleModel, CircleRouting, CircleRoutingCustomMatch, CircleRoutingSegment } from '../../../core/types/circle'
import { useParams, useSearchParams } from 'react-router-dom'
import CircleModules from '../../CircleModules'
import ViewInput from '../../../core/components/ViewInput'
import { useAppSelector } from '../../../core/hooks/redux'
import FloatingButton from '../../../core/components/FloatingButton'
import Editor from '../../../core/components/Editor'
import { CIRCLE_VIEW_MODE } from '../../../core/types/circle'

interface Props {
  circle?: CircleModel
  viewMode: CIRCLE_VIEW_MODE
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

const CircleForm = ({ circle, viewMode, onSave }: Props) => {
  const { workspace } = useAppSelector(state => state.main)
  const [hasChange, setHasChange] = useState(false)
  const [form, setForm] = useState({
    name: circle?.name || '',
    description: circle?.description || '',
    customMatch: JSON.stringify(circle?.routing?.match?.customMatch, null, 2) || JSON.stringify(initialCustomMatch, null, 2),
    segments: JSON.stringify(initialSegments, null, 2),
    environments: JSON.stringify(circle?.environments, null, 2) || JSON.stringify(initialEnviroments, null, 2),
  })
  const [matchStrategy, setMatchStrategy] = useState('customMatch')

  const handleChange = (input: any) => {
    setHasChange(true)
    setForm(state => ({...state, ...input}))
  }

  const isCreate = () => {
    return viewMode === CIRCLE_VIEW_MODE.CREATE
  }

  const handleClickSave = () => {
    let routing: CircleRouting = { strategy: workspace.routingStrategy }
    if (workspace.routingStrategy === 'match') {
      if (matchStrategy === 'customMatch') {
        routing = {
          ...routing,
          match: { customMatch: JSON.parse(form.customMatch) } ,
        }
      } else {
        routing = {
          ...routing,
          match: { segments: JSON.parse(form.segments) }
        }
      }
    }

    const newCircle = {
      name: form.name,
      description: form.description,
      environments: JSON.parse(form.environments),
      modules: [],
      routing,
    }

    onSave(newCircle)
  }


  return (
    <div className='circle__content'>
      {(isCreate() || hasChange) && (
        <FloatingButton
          icon="check"
          iconColor='white'
          text="Save circle"
          onClick={handleClickSave}
        />
      )}
      <ViewInput.Text
        icon={["far", "circle"]}
        label='Name'
        value={form.name}
        edit={isCreate()}
        canEdit={false}
        onChange={value => handleChange({ name: value })}
        placeholder="Circle name"
      />
      <ViewInput.Text
        icon="align-justify"
        label='Description'
        value={form.description}
        edit={isCreate()}
        onChange={value => handleChange({ description: value })}
        as="textarea"
        placeholder="Circle description"
      />
      <ViewInput
        label="Routing"
        icon="route"
      >
        <Nav variant='pills' defaultActiveKey={matchStrategy} onSelect={key => setMatchStrategy(key || 'customMatch')} className='mb-3'>
          <Nav.Item><Nav.Link eventKey="customMatch">Custom match</Nav.Link></Nav.Item>
          <Nav.Item><Nav.Link eventKey="segments">Segments</Nav.Link></Nav.Item>
        </Nav>
        {matchStrategy === 'customMatch' ? (
          <div className='circle__content__section__custom-match'>
            <Editor
              height='200px'
              value={form.customMatch}
              onChange={value => handleChange({ customMatch: value })}
            />
          </div>
        ) : (
          <div className='circle__content__section__segments'>
            <Editor
              height='200px'
              value={form.segments}
              onChange={value => handleChange({ segments: value })}
            />
          </div>
        )}
      </ViewInput>
      <ViewInput
        label="Modules"
        icon="folder"
      >
        {circle && <CircleModules circle={circle} />}
      </ViewInput>
      <ViewInput
        label="Modules"
        icon="folder"
      >
        <Editor
          height='200px'
          value={form.environments}
          onChange={value => handleChange({ environments: value })}
        />
      </ViewInput>
    </div>
  )
}

export default CircleForm
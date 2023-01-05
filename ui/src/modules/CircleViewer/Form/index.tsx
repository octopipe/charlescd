import React, { useCallback, useEffect, useState } from 'react'
import { Nav } from 'react-bootstrap'
import './style.scss'
import { Circle, Circle as CircleType, CircleEnrivonment, CircleModel, CircleModule, CircleRouting, CircleRoutingSegment, CIRCLE_ROUTING_STRATEGY } from '../../../core/types/circle'
import { useParams, useSearchParams } from 'react-router-dom'
import CircleModules from '../../CircleModules'
import ViewInput from '../../../core/components/ViewInput'
import { useAppSelector } from '../../../core/hooks/redux'
import FloatingButton from '../../../core/components/FloatingButton'
import Editor from '../../../core/components/Editor'
import { CIRCLE_VIEW_MODE } from '../../../core/types/circle'
import { useForm } from 'react-hook-form'
import FormControl from '../../../core/components/FormControl'
import { formValidations } from '../../../core/form/validate'

interface Props {
  circle?: CircleModel
  viewMode: CIRCLE_VIEW_MODE
  onUpdate: (circle: CircleType) => void
  onSave: (circle: CircleType) => void
}

const initialEnviroments = [
  { key: 'KEY_EXAMPLE', value: 'VALUE_EXAMPLE' }
]

const initialCustomMatch = { headers: { 'x-header-example': '1111' } }

const initialSegments = [
  { key: 'email', op: 'EQUAL', value: 'email@mail.com' }
]

const initialMatchRouting = {
  strategy: CIRCLE_ROUTING_STRATEGY,
  match: initialCustomMatch,
}

const getInitialRouting = (routingStrategy: string, routing?: CircleRouting) => {
  if (routingStrategy === CIRCLE_ROUTING_STRATEGY.MATCH) {
    return { strategy: CIRCLE_ROUTING_STRATEGY.MATCH, match: JSON.stringify(routing?.match || initialCustomMatch, null, 2) }
  }

  if (routingStrategy === CIRCLE_ROUTING_STRATEGY.SEGMENTS) {
    return { strategy: CIRCLE_ROUTING_STRATEGY.SEGMENTS, segments: JSON.stringify(routing?.segments || initialSegments, null, 2) }
  }

  return { strategy: CIRCLE_ROUTING_STRATEGY.CANARY, canary: { weight: 0 } }
}

interface CircleForm {
  name: string
  description: string
  modules: CircleModule[]
  environments: string
  routingMatch: string
  routingSegments: string
  weight: number
}

const CircleForm = ({ circle, viewMode, onSave, onUpdate }: Props) => {
  const { workspace } = useAppSelector(state => state.main)
  const [hasChange, setHasChange] = useState(false)
  const [matchStrategy, setMatchStrategy] = useState<CIRCLE_ROUTING_STRATEGY>(CIRCLE_ROUTING_STRATEGY.MATCH)
  const { control, formState: {errors}, watch, handleSubmit, setValue } = useForm<CircleForm>({
    defaultValues: {
      name: circle?.name || '',
      description: circle?.description || '',
      routingMatch: JSON.stringify(circle?.routing?.match || initialCustomMatch, null, 2),
      routingSegments: JSON.stringify(circle?.routing?.segments || initialSegments, null, 2),
      weight: 0,
      modules: circle?.modules || [],
      environments: JSON.stringify(circle?.environments || initialEnviroments, null, 2)
    }
  })

  useEffect(() => {
    const subscription = watch((value, { name, type }) => setHasChange(true));
    return () => subscription.unsubscribe();
  }, [watch])

  const isCreate = () => {
    return viewMode === CIRCLE_VIEW_MODE.CREATE
  }

  const handleClickSave = (form: CircleForm) => {
    let routing: CircleRouting = { strategy: CIRCLE_ROUTING_STRATEGY.MATCH, match: JSON.parse(form.routingMatch) }
  
    if (workspace.routingStrategy === CIRCLE_ROUTING_STRATEGY.SEGMENTS) {
      routing = { strategy: CIRCLE_ROUTING_STRATEGY.SEGMENTS, segments: JSON.parse(form.routingSegments) }
    }
  
    if (workspace.routingStrategy === CIRCLE_ROUTING_STRATEGY.SEGMENTS) {
      routing = { strategy: CIRCLE_ROUTING_STRATEGY.CANARY, canary: { weight: form.weight } }
    }

    let curretCircle: Circle = {
      name: form.name,
      description: form.description,
      modules: form.modules,
      environments: JSON.parse(form.environments),
      routing: routing
    }
    
    
    if (isCreate()) {
      console.log(curretCircle)
      onSave(curretCircle)
      setHasChange(false)
    } else {
      onUpdate(curretCircle)
      setHasChange(false)
    }
  }

  const handleDeleteModule = (module: CircleModule) => {
    if (!circle) {
      return
    }

    const modules = circle?.modules?.filter(m => m.name !== module.name) || []
    let curretCircle: Circle = {
      name: circle.name,
      description: circle.description,
      modules: modules,
      environments: circle.environments,
      routing: circle.routing
    }
    onUpdate(curretCircle)
  }


  return (
    <div className='circle__content'>
      {(isCreate() || hasChange) && (
        <FloatingButton
          icon="check"
          iconColor='white'
          text="Save circle"
          onClick={handleSubmit(handleClickSave)}
        />
      )}

      <FormControl name="name" control={control} rules={{required: 'Name is required'}} errors={errors}>
        <ViewInput.Text
          icon={["far", "circle"]}
          label='Name'
          edit={isCreate()}
          canEdit={false}
          placeholder="Circle name"
        />
      </FormControl>

      <FormControl name="description" control={control} errors={errors}>
        <ViewInput.Text
          icon="align-justify"
          label='Description'
          edit={isCreate()}
          as="textarea"
          placeholder="Circle description"
        />
      </FormControl>
      
      {workspace.routingStrategy !== CIRCLE_ROUTING_STRATEGY.CANARY && (
        <ViewInput
          label="Routing"
          icon="route"
        >
          <Nav variant='pills' defaultActiveKey={matchStrategy} onSelect={key => setMatchStrategy(key as CIRCLE_ROUTING_STRATEGY)} className='mb-3'>
            <Nav.Item><Nav.Link eventKey={CIRCLE_ROUTING_STRATEGY.MATCH}>Custom match</Nav.Link></Nav.Item>
            <Nav.Item><Nav.Link eventKey={CIRCLE_ROUTING_STRATEGY.SEGMENTS}>Segments</Nav.Link></Nav.Item>
          </Nav>
          <div>
            {matchStrategy === CIRCLE_ROUTING_STRATEGY.MATCH && (
              <div className='circle__content__section__custom-match'>
                <FormControl name="routingMatch" rules={{validate: formValidations.jsonValidate}} control={control} errors={errors}>
                  <Editor
                    height='200px'
                  />
                </FormControl>
              </div>
            )}
            {matchStrategy === CIRCLE_ROUTING_STRATEGY.SEGMENTS && (
              <div className='circle__content__section__segments'>
                <FormControl name="routingSegments" rules={{validate: formValidations.jsonValidate}} control={control} errors={errors}>
                  <Editor
                    height='200px'
                  />
                </FormControl>
              </div>
            )}
          </div>
        </ViewInput>
      )}
      {workspace.routingStrategy === 'CANARY' && (
        <FormControl name="weight" rules={{ min: 0, max: 100 }} control={control} errors={errors}>
          <ViewInput.Text
            icon="percent"
            label='Weight'
            edit={isCreate()}
            placeholder="Circle canary weight"
            type='number'
          />
        </FormControl>
      )}
      <ViewInput
        label="Modules"
        icon="folder"
      >
        <CircleModules onDelete={handleDeleteModule} onChangeModules={modules => setValue('modules', modules)} circle={circle} />
      </ViewInput>
      <ViewInput
        label="Environments"
        icon="folder"
      >
        <FormControl name="environments" rules={{validate: formValidations.jsonValidate}} control={control} errors={errors}>
          <Editor
            height='200px'
          />
        </FormControl>
      </ViewInput>
      <ViewInput
        label='Metrics'
        icon="chart-simple"
      >

      </ViewInput>
    </div>
  )
}

export default CircleForm
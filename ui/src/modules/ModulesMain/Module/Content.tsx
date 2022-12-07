import React, { useEffect, useState } from 'react'
import { Nav } from 'react-bootstrap'
import useFetch from 'use-http'
import './style.scss'
import { Module as ModuleType, ModuleModel } from '../types'
import { useParams, useSearchParams } from 'react-router-dom'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import AceEditor from "react-ace";

import "ace-builds/src-noconflict/mode-json";
import "ace-builds/src-noconflict/theme-monokai";
import ViewInput from '../../../core/components/ViewInput'
import { useAppSelector } from '../../../core/hooks/redux'
import FloatingButton from '../../../core/components/FloatingButton'

interface Props {
  moduleId: string
  moduleOp: string
  onUpdate?: (id: string) => void
  onSave: (module: ModuleType) => void
}

const initialEnviroments = [
  { key: 'KEY_EXAMPLE', value: 'VALUE_EXAMPLE' }
]

const initialCustomMatch = { headers: { 'x-header-example': '1111' } }

const initialSegments = [
  { key: 'email', op: 'EQUAL', value: 'email@mail.com' }
]

const ModuleContent = ({ moduleId, moduleOp, onSave }: Props) => {
  const [searchParams] = useSearchParams();
  const { workspaceId } = useParams()
  const [module, setModule] = useState<ModuleModel>()
  const { response, get } = useFetch()
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [url, setUrl] = useState('')
  const [path, setPath] = useState('')
  const [templateType, setTemplateType] = useState('')
  const [visibility, setVisibility] = useState('')


  const loadModule = async () => {
    const module = await get(`/workspaces/${workspaceId}/modules/${moduleId}`)
    if (response.ok) setModule(module || [])
  }

  useEffect(() => {
    if (moduleOp !== "C")
      loadModule()
  }, [workspaceId])

  useEffect(() => {
    setName(module?.name || '')
    setDescription(module?.description || '')
    setUrl(module?.url || '')
    setPath(module?.path || '')
    setTemplateType(module?.templateType || '')
    setVisibility(module?.visibility || '')
  }, [module])

  const isCreate = () => {
    return searchParams.get(moduleId) === 'C'
  }

  const handleClickSave = () => {
    const newModule = {
      name,
      description,
      path,
      url,
      templateType,
      visibility,
    }

    onSave(newModule)
  }


  return (
    <div className='module__content'>
      {moduleOp === 'C' && (
        <FloatingButton
          icon="check"
          iconColor='#4caf50'
          text="Save module"
          onClick={handleClickSave}
        />
      )}
      <div className='module__content'>
        <ViewInput.Text
          icon="folder"
          label='Name'
          value={name}
          edit={isCreate()}
          canEdit={false}
          onChange={setName}
          placeholder="Module name"
        />
        <ViewInput.Text
          icon="align-justify"
          label='Description'
          value={description}
          edit={isCreate()}
          onChange={setDescription}
          as="textarea"
          placeholder="Module description"
        />
        <ViewInput.Text
          icon={["fab", "git-alt"]}
          label='Url'
          value={url}
          edit={isCreate()}
          onChange={setUrl}
          placeholder="Module name"
        />
        <ViewInput.Text
          icon="folder"
          label='Path'
          value={path}
          edit={isCreate()}
          onChange={setPath}
          placeholder="Module name"
        />
      </div>
    </div>
  )
}

export default ModuleContent
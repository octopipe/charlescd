import React, { useCallback, useEffect, useState } from 'react'
import { Module as ModuleType, ModuleModel } from '../types'
import { useParams, useSearchParams } from 'react-router-dom'
import Alert from '../../../core/components/Alert'
import ModuleContent from './Content'
import ModuleTabs from './Tabs'
import './style.scss'
import useFetch from '../../../core/hooks/fetch'

interface Props {
  moduleId: string
  moduleOp: string
  onClose: (id: string) => void
  onUpdate: (id: string) => void
  onSave: (module: ModuleType) => void
  onDelete: (moduleId: string) => void
}



enum TABS {
  CONTENT = 'content',
  TREE = 'tree'
}

const Module = ({ moduleId, moduleOp, onClose, onSave, onDelete }: Props) => {
  const [searchParams] = useSearchParams();
  const { workspaceId } = useParams()
  const [module, setModule] = useState<ModuleModel>()
  const { fetch } = useFetch()
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [ showDeleteAlert, toggleDeleteAlert ] = useState(false)
  const [activeTab, setActiveTab] = useState(TABS.CONTENT)
  const [modules, setModules] = useState([])

  useEffect(() => {
    fetch(`/workspaces/${workspaceId}/modules/${moduleId}`).then(res => setModule(module))
  }, [])

  useEffect(() => {
    setName(module?.name || '')
    setDescription(module?.description || '')
  }, [module])

  const handleDelete = (moduleId: string) => {
    onDelete(moduleId)
    toggleDeleteAlert(false)
  }

  const CustomToggle = React.forwardRef<any, any>(({ children, onClick }, ref) => (
    <a
      ref={ref}
      onClick={(e) => {
        e.preventDefault();
        onClick(e);
      }}
      className="module-modules__item__menu"
    >
      {children}
    </a>
  ));

  return (
    <div className='module'>
      <ModuleTabs
        moduleId={moduleId}
        activeTab={activeTab}
        moduleOp={moduleOp}
        module={module}
        onClose={onClose}
        onChange={tab => setActiveTab(tab)}
        onDelete={handleDelete}
      />
      {activeTab === TABS.CONTENT && <ModuleContent moduleId={moduleId} moduleOp={moduleOp} onSave={onSave} />}
      <Alert show={showDeleteAlert} action={() => handleDelete(moduleId)} onClose={() => toggleDeleteAlert(false)} />
    </div>
  )
}

export default Module
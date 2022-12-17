import React, { useEffect, useState } from 'react'
import './style.scss'
import { Module as ModuleType, ModulePagination } from './types'
import { useLocation, useParams, useSearchParams } from 'react-router-dom'
import Module from './Module'
import { useAppDispatch } from '../../core/hooks/redux'
import { setBreadcrumbItems } from '../Main/mainSlice'
import useFetch from '../../core/hooks/fetch';
import AppSidebar from '../../core/components/AppSidebar';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

const createModuleId = 'untitled'

const ModulesMain = () => {
  const location = useLocation()
  const [searchParams, setSearchParams] = useSearchParams();
  const { workspaceId } = useParams()
  const [modules, setModules] = useState<ModulePagination>({continue: '', items: []})
  const [activeModuleIds, setActiveModulesIds] = useState<string[]>([])
  const dispatch = useAppDispatch()
  const { fetch,  loading } = useFetch()

  const loadModules = async () => {
    const modules = await fetch(`/workspaces/${workspaceId}/modules`)
    setModules(modules || [])
  }

  useEffect(() => {
    loadModules()
    dispatch(setBreadcrumbItems([
      { name: 'Modules', to: `/workspaces/${workspaceId}/modules` },
    ]))
  }, [])

  useEffect(() => {
    let currentActiveModulesIds: string[] = []
    searchParams.forEach((value, key) => currentActiveModulesIds.push(key))
    setActiveModulesIds(currentActiveModulesIds)
  }, [location])


  const handleModuleClick = (moduleId: string) => {
    setSearchParams(i => {
      if (!i.has(moduleId)) {
        i.append(moduleId, "R")
      } else {
        i.delete(moduleId)
      }

      return i
    })
  }

  const handleModuleCreateClick = () => {
    
    setSearchParams(i => {
      if (!i.has(createModuleId)) {
        i.append(createModuleId, "C")
      }

      return i
    })
  }

  const handleCloseModule = (moduleId: string) => {
    setSearchParams(i => {
      i.delete(moduleId)
      return i
    })
  }

  const handleUpdateModule = (moduleId: string) => {
    setSearchParams(i => {
      i.set(moduleId, "U")
      return i
    })
  }

  const handleDeleteModule = async (moduleId: string) => {
    await fetch(`/workspaces/${workspaceId}/modules/${moduleId}`, { method: 'DELETE' })
    await loadModules()
    setSearchParams(i => {
      i.delete(moduleId)
      return i
    })
  }

  const handleSaveModule = async (module: ModuleType) => {
    const newModule = await fetch(`/workspaces/${workspaceId}/modules`, {method: 'POST', data: module})
    setSearchParams(i => {
      if (!i.has(newModule.id)) {
        i.append(newModule.id, "R")
      }

      i.delete(createModuleId)

      return i
    })
  }

  return (
    <div className='modules'>
      <AppSidebar>
        <AppSidebar.Header>
          <AppSidebar.HeaderItem onClick={handleModuleCreateClick}>
            <FontAwesomeIcon icon="plus-circle" className="me-1" /> Create circle
          </AppSidebar.HeaderItem>
        </AppSidebar.Header>
        <AppSidebar.List loading={loading}>
          {modules&& modules.items.length > 0 && modules?.items.map(item => (
            <AppSidebar.ListItem
              key={item.id}
              isActive={searchParams.has(item.id)}
              icon="folder"
              activeIcon="folder"
              text={item.name}
              onClick={() => handleModuleClick(item.id)}
            />
          ))}
        </AppSidebar.List>
      </AppSidebar>
      <div className='modules__content'>
        {activeModuleIds.map(id => (
          <Module
            key={id}
            moduleId={id}
            moduleOp={searchParams.get(id) || 'R'}
            onClose={handleCloseModule}
            onUpdate={handleUpdateModule}
            onSave={handleSaveModule}
            onDelete={handleDeleteModule}
          />
        ))}
      </div>
    </div>
  )
}

export default ModulesMain
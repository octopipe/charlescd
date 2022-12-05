import React, { useCallback, useEffect, useState } from 'react'
import useFetch, { CachePolicies } from 'use-http'
import './style.scss'
import { Module as ModuleType, ModulePagination } from './types'
import Placeholder from '../../core/components/Placeholder'
import { Link, Navigate, useLocation, useNavigate, useParams, useSearchParams } from 'react-router-dom'
import ModulesSidebar from './Sidebar';
import Module from './Module'

const createModuleId = 'untitled'

const ModulesMain = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const [searchParams, setSearchParams] = useSearchParams();
  const { workspaceId } = useParams()
  const [modules, setModules] = useState<ModulePagination>({continue: '', items: []})
  const [activeModuleIds, setActiveModulesIds] = useState<string[]>([])
  const { response, get, post, delete: deleteMethod } = useFetch({cachePolicy: CachePolicies.NO_CACHE})

  const loadModules = async () => {
    const modules = await get(`/workspaces/${workspaceId}/modules`)
    if (response.ok) setModules(modules || [])
  }

  useEffect(() => {
    loadModules()
  }, [workspaceId])

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
    await deleteMethod(`/workspaces/${workspaceId}/modules/${moduleId}`)
    await loadModules()
    if (response.ok) {
      setSearchParams(i => {
        i.delete(moduleId)
        return i
      })
    }
  }

  const handleSaveModule = async (module: ModuleType) => {
    const newModule = await post(`/workspaces/${workspaceId}/modules`, module)
    if (response.ok) {
      setSearchParams(i => {
        if (!i.has(newModule.id)) {
          i.append(newModule.id, "R")
        }

        i.delete(createModuleId)
  
        return i
      })
      await loadModules()
    }
  }

  return (
    <div className='modules'>
      <ModulesSidebar
        modules={modules}
        onModuleClick={handleModuleClick}
        onModuleCreateClick={handleModuleCreateClick}
      />
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
import React, { useCallback, useEffect, useState } from 'react';
import { generatePath, matchRoutes, Outlet, useLocation, useNavigate } from 'react-router-dom';
import useFetch from 'use-http'
import MainNavbar from './Navbar';
import MainSidebar from './Sidebar';
import './style.scss'

const routes = [
  { path: '/workspaces/:workspaceId' },
  { path: '/workspaces/:workspaceId/circles' }
]

const Main = () => {
  const location = useLocation()
  const navigate = useNavigate()
  const [workspaces, setWorkspaces] = useState<any[]>([])
  const { response, get } = useFetch()

  const loadWorkspaces = useCallback(async () => {
    const workspaces = await get('/workspaces')
    if (response.ok) setWorkspaces(workspaces)
  }, [get, response])

  useEffect(() => {
    loadWorkspaces()
  }, [])

  const goToWorkspacePage = (workspaceId: string) => {
    const [{ route }] = matchRoutes(routes, location) || []
    navigate(generatePath(route?.path || '' , { workspaceId }))
  }

  useEffect(() => {
    if (workspaces?.length <= 0)
      return

    goToWorkspacePage(workspaces[0].id)
  }, [workspaces])

  return (
    <div className='main'>
      <MainNavbar
        workspaces={workspaces || []}
        onSelectWorkspace={(workspaceId: any) => goToWorkspacePage(workspaceId)}
      />
      <div className='main__content'>
        <MainSidebar />
        <Outlet />
      </div>
    </div>
  )
}

export default Main 
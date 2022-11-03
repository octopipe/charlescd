import React, { useCallback, useEffect, useState } from 'react';
import { generatePath, matchRoutes, Outlet, useLocation, useNavigate } from 'react-router-dom';
import useFetch from 'use-http'
import Placeholder from '../../core/components/Placeholder';
import { ReactComponent as EmptyWorkspaces } from '../../core/assets/svg/empty-workspaces.svg'
import MainNavbar from './Navbar';
import MainSidebar from './Sidebar';
import './style.scss'
import { Container } from 'react-bootstrap';
import { useAppDispatch } from '../../core/hooks/redux';
import { setDeployStrategy } from './mainSlice';

const routes = [
  { path: '' },
  { path: '/workspaces/:workspaceId' },
  { path: '/workspaces/:workspaceId/circles' }
]

const EmptyWorkspacesPlaceholser = () => (
  <Container>
    <Placeholder text="You don't have any workspace...">
      <EmptyWorkspaces />
    </Placeholder>
  </Container>
)

const Main = () => {
  const dispatch = useAppDispatch()
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

  const goToWorkspacePage = (workspace: any) => {
    const matches = matchRoutes(routes, location) || []
    dispatch(setDeployStrategy(workspace.deployStrategy))

    if (matches?.length <= 0 || routes[0]?.path === '') {
      navigate(`/workspaces/${workspace.id}`)
      return
    }


    navigate(generatePath(routes[0]?.path || '' , { workspaceId: workspace.id }))
  }

  useEffect(() => {
    if (workspaces?.length <= 0)
      return

    
    goToWorkspacePage(workspaces[0])
  }, [workspaces])
  


  return (
    <div className='main'>
      <MainNavbar
        workspaces={workspaces || []}
        onSelectWorkspace={(workspaceId: any) => goToWorkspacePage(workspaceId)}
      />
      <div className='main__content'>
        {workspaces?.length > 0 &&  <MainSidebar />}
        <div className='main__content__body'>

          {workspaces?.length > 0 ? <Outlet /> : <EmptyWorkspacesPlaceholser />}
        </div>
      </div>
    </div>
  )
}

export default Main 
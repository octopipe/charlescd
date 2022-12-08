import React, { useCallback, useEffect, useState } from 'react';
import { generatePath, matchRoutes, Outlet, useLocation, useNavigate, useParams } from 'react-router-dom';
import useFetch from 'use-http'
import Placeholder from '../../core/components/Placeholder';
import { ReactComponent as EmptyWorkspaces } from '../../core/assets/svg/empty-workspaces.svg'
import MainSidebar from './Sidebar';
import './style.scss'
import { Container } from 'react-bootstrap';
import { useAppDispatch } from '../../core/hooks/redux';
import { setDeployStrategy } from './mainSlice';
import Navbar from '../../core/components/Navbar';

const EmptyWorkspacesPlaceholser = () => (
  <Container>
    <Placeholder text="You don't have any workspace...">
      <EmptyWorkspaces />
    </Placeholder>
  </Container>
)

const Main = () => {
  const { workspaceId } = useParams()
  const [workspaces, setWorkspaces] = useState<any[]>([])
  const [selectedWorkspaceId, setSelectedWorkspaceId] = useState(workspaceId)
  const { response, get } = useFetch()

  const loadWorkspaces = useCallback(async () => {
    const workspaces = await get('/workspaces')
    if (response.ok) setWorkspaces(workspaces)
  }, [get, response])

  useEffect(() => {
    loadWorkspaces()
  }, [])


  return (
    <div className='main'>
      <Navbar
        workspaces={workspaces || []}
        selectedWorkspaceId={selectedWorkspaceId || ''}
        onSelectWorkspace={(workspaceId: any) => setSelectedWorkspaceId(workspaceId)}
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
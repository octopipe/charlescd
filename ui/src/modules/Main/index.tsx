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
import { ROUTES } from '../../core/constants/routes';

const EmptyWorkspacesPlaceholser = () => (
  <Container>
    <Placeholder text="You don't have any workspace...">
      <EmptyWorkspaces />
    </Placeholder>
  </Container>
)

const repalacePathParams = (path: string, params: object, prefix = ':') => {
  let newPath = path

  Object.entries(params).forEach(([key, value]) => {
    newPath = newPath.replace(prefix + key, value)
  })
  return newPath
}


const Main = () => {
  const { workspaceId } = useParams()
  const location = useLocation()
  const navigate = useNavigate()
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

  useEffect(() => {
    const routeMatches = matchRoutes(Object.values(ROUTES).map(i => ({ path: i })), location)
    if (routeMatches && routeMatches?.length > 0) {
      console.log(location)
      navigate(`${repalacePathParams(routeMatches[0].route?.path, { workspaceId: selectedWorkspaceId })}${location.search}`)
    }
    
  }, [selectedWorkspaceId])


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
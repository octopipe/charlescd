import React, { useCallback, useEffect, useState } from 'react';
import { generatePath, matchRoutes, Outlet, useLocation, useNavigate, useParams } from 'react-router-dom';
import Placeholder from '../../core/components/Placeholder';
import { ReactComponent as EmptyWorkspaces } from '../../core/assets/svg/empty-workspaces.svg'
import MainSidebar from './Sidebar';
import { Container } from 'react-bootstrap';
import { useAppDispatch, useAppSelector } from '../../core/hooks/redux';
import Navbar from '../../core/components/Navbar';
import { WorkspaceModel } from '../Workspaces/types';
import { setWorkspace as setWorkspaceDispatch } from '../Main/mainSlice'
import './style.scss'
import DynamicContainer from '../../core/components/DynamicContainer';
import useFetch from '../../core/hooks/fetch';

const Main = () => {
  const { workspaceId } = useParams()
  const [workspace, setWorkspace] = useState<WorkspaceModel>()
  const { loading, fetch } = useFetch()
  const dispatch = useAppDispatch()
  const { breadcrumbItems } = useAppSelector(state => state.main)

  useEffect(() => {
    fetch(`/workspaces/${workspaceId}`, "GET").then(res => setWorkspace(res))
  }, [])

  useEffect(() => {
    dispatch(setWorkspaceDispatch(workspace))
  }, [workspace])


  return (
    <div className='main'>
      <Navbar workspace={workspace} breadcrumbItems={breadcrumbItems} />
      <DynamicContainer loading={loading} className='main__content'>
        <MainSidebar />
        <div className='main__content__body'>
          <Outlet />
        </div>
      </DynamicContainer>
    </div>
  )
}

export default Main 
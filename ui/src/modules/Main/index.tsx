import React, { useCallback, useEffect, useState } from 'react';
import { Outlet } from 'react-router-dom';
import useFetch from 'use-http'
import { useAppDispatch } from '../../core/hooks/redux';
import { setCurrentWorkspace } from './mainSlice';
import MainNavbar from './Navbar';
import MainSidebar from './Sidebar';
import './style.scss'


const Main = () => {
  const dispatch = useAppDispatch()
  const [workspaces, setWorkspaces] = useState<any[]>([])
  const { response, get } = useFetch()

  const loadWorkspaces = useCallback(async () => {
    const workspaces = await get('/workspaces')
    if (response.ok) setWorkspaces(workspaces)
  }, [get, response])

  useEffect(() => {
    loadWorkspaces()
  }, [])

  useEffect(() => {
    if (workspaces?.length <= 0)
      return

    dispatch(setCurrentWorkspace(workspaces[0].id))
  }, [workspaces])

  return (
    <div className='main'>
      <MainNavbar
        workspaces={workspaces || []}
        onSelectWorkspace={(workspaceId: any) => dispatch(setCurrentWorkspace(workspaceId))}
      />
      <div className='main__content'>
        <MainSidebar />
        <Outlet />
      </div>
    </div>
  )
}

export default Main 
import React, { useCallback, useEffect, useState } from 'react';
import { Outlet } from 'react-router-dom';
import useFetch from 'use-http'
import MainNavbar from './Navbar';
import MainSidebar from './Sidebar';
import './style.scss'


const Main = () => {
  const [workspaces, setWorkspaces] = useState<any[]>()
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
      <MainSidebar />
      <div className='main__content'>
        <MainNavbar workspaces={workspaces} />
        <Outlet />
      </div>
    </div>
  )
}

export default Main 
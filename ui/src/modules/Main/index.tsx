import React from 'react'
import { Outlet } from 'react-router-dom'
import Sidebar from './Sidebar'
import './style.css'

const Main = () => {
  return (
    <div className='main'>
      <Sidebar />
      <div className='col-md-11'>
        <Outlet />
      </div>
    </div>
  )
}

export default Main 
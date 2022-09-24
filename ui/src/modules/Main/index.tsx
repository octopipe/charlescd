import React from 'react'
import { Outlet } from 'react-router-dom'
import Sidebar from './Sidebar'
import './style.css'

const Main = () => {
  return (
    <div className='main'>
      <Sidebar />
      <div className='col-md-9 ms-sm-auto col-lg-10 px-md-4'>
        <Outlet />
      </div>
    </div>
  )
}

export default Main 
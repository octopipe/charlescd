import React from 'react'
import { Outlet } from 'react-router-dom'
import Navbar from '../../core/components/Navbar'
import './style.scss'

const Root = () => {
  return (
    <>
      <Navbar />
      <Outlet />
    </>
  )
}

export default Root
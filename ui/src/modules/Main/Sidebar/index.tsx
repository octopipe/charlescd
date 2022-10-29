import { IconProp } from '@fortawesome/fontawesome-svg-core';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import React, { useState } from 'react';
import { Link, NavLink, useMatch } from 'react-router-dom';
import './style.scss'

const items = [
  {
    name: 'Home',
    to: '/',
    icon: 'home'
  },
  {
    name: 'Circles',
    to: '/circles',
    icon: ["far", "circle"]
  },
  {
    name: 'Modules',
    to: '/modules',
    icon: 'folder'
  },
  {
    name: 'Routes',
    to: '/routes',
    icon: 'signs-post'
  }
]

const MainSidebar = () => {
  return (
    <div className='main__sidebar'>
      <div className='main__sidebar__list'>
        {items.map(item => (
            <NavLink 
              className={({ isActive }) => isActive ? 'main__sidebar__list__item--active' : 'main__sidebar__list__item'}
              to={item.to}
              end
            >
              <FontAwesomeIcon icon={item.icon as IconProp} />
            </NavLink>
        ))}
      </div>
    </div>
  )
}

export default MainSidebar
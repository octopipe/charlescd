import React from 'react';
import { IconProp } from '@fortawesome/fontawesome-svg-core';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { NavLink, useParams } from 'react-router-dom';
import { ReactComponent as LogoWhite } from '../../../core/assets/svg/logo-white.svg'
import './style.scss'

const items = [
  {
    name: 'Home',
    to: '',
    icon: 'home'
  },
  {
    name: 'Circles',
    to: 'circles',
    icon: ["far", "circle"]
  },
  {
    name: 'Modules',
    to: 'modules',
    icon: 'folder'
  },
  {
    name: 'Routes',
    to: 'routes',
    icon: 'signs-post'
  }
]

const MainSidebar = () => {
  const { workspaceId } = useParams()

  return (
    <div className='main__sidebar'>
      <div className='main__sidebar__logo'>
        <NavLink to='/'>
          <LogoWhite />
        </NavLink>
      </div>
      <div className='main__sidebar__list'>
        {items.map(item => (
          <NavLink 
            key={item.name}
            className={({ isActive }) => isActive ? 'main__sidebar__list__item--active' : 'main__sidebar__list__item'}
            to={`/workspaces/${workspaceId}/${item.to}`}
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
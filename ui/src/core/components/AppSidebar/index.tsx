import { IconProp } from '@fortawesome/fontawesome-svg-core'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import React from 'react'
import { Nav } from 'react-bootstrap'
import DynamicContainer from '../DynamicContainer'
import './style.scss'

interface AppSidebarBaseProps {
  children: React.ReactNode
}

interface HeaderItemProps extends AppSidebarBaseProps {
  onClick: () => void
}

const Header = ({children}: AppSidebarBaseProps) => (
  <div className='app-sidebar__header'>
    <Nav>
      {children}
    </Nav>
  </div>
) 

const HeaderItem = ({children, onClick}: HeaderItemProps) => (
  <Nav.Item onClick={onClick}>
    {children}
  </Nav.Item>
) 

interface ListProps extends AppSidebarBaseProps {
  loading: boolean
}

const List = ({children, loading}: ListProps) => (
  <DynamicContainer loading={loading} className='app-sidebar__list'>
    {children}
  </DynamicContainer>
)

interface ListItemProps {
  isActive: boolean
  icon: IconProp
  activeIcon: IconProp
  text: string
  onClick: () => void
}

const ListItem = ({isActive, icon, activeIcon, text, onClick}: ListItemProps) => (
  <div className={isActive ? 'app-sidebar__list__item--active' : 'app-sidebar__list__item'} onClick={onClick}>
    <FontAwesomeIcon icon={isActive ? activeIcon : icon} className="me-2" />
    {text}
  </div>
)

const AppSidebar = ({children}: AppSidebarBaseProps) => (
  <div className='app-sidebar'>
    {children}
  </div>
)

AppSidebar.Header = Header
AppSidebar.HeaderItem = HeaderItem
AppSidebar.List = List
AppSidebar.ListItem = ListItem

export default AppSidebar
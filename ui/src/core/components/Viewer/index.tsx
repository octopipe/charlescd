import React from 'react'
import { IconProp } from '@fortawesome/fontawesome-svg-core'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import './style.scss'
import { Dropdown } from 'react-bootstrap'

interface ViewerPropsBase {
  children: React.ReactNode
}

interface ViewerTabsOption {
  name: string
  onAction: () => void
}

interface ViewerTabsProps {
  children: React.ReactNode
  hasOptions?: boolean
  options?: ViewerTabsOption[]
}

const ViewerTabs = ({ children, hasOptions, options }: ViewerTabsProps) => {
  const CustomToggle = React.forwardRef<any, any>(({ children, onClick }, ref) => (
    <a
      ref={ref}
      onClick={(e) => {
        e.preventDefault();
        onClick(e);
      }}
      className="circle-modules__item__menu"
    >
      {children}
    </a>
  ));

  return (
    <div className='viewer__tabs'>
      <div className='d-flex'>
        {children}
      </div>
      {hasOptions && (
        <div className='viewer__tabs__options'>
          <Dropdown className='mx-2'>
            <Dropdown.Toggle as={CustomToggle}>
              <FontAwesomeIcon icon="ellipsis-vertical" />
            </Dropdown.Toggle>
            <Dropdown.Menu>
              { options?.map(option => (
                <Dropdown.Item onClick={option.onAction}>
                  {option.name}
                </Dropdown.Item>
              )) }
            </Dropdown.Menu>
          </Dropdown>
        </div>
      )}
    </div>
  )
}

interface ViewerTabsItemProps {
  id: string
  icon: IconProp
  text: string
  isActive: boolean
  hasClose?: boolean
  onClick: () => void
  onClose?: () => void
}

const ViewerTabsItem = ({ id, isActive, icon, text, hasClose, onClick, onClose }: ViewerTabsItemProps) => {

  const handleClose = (e: any) => {
    e.stopPropagation()
    if (onClose)
      onClose()
  }

  const handleClick = (e: any) => {
    e.preventDefault()
    onClick()
  }

  return (
    <div className={isActive ? `viewer__tabs__item--active` : `viewer__tabs__item`} onClick={handleClick}>
      <div>
        <FontAwesomeIcon icon={icon} className="me-2" /> {text}
      </div>
      {hasClose && onClose && (
        <div onClick={handleClose} className="viewer__tabs__item__close">
          <FontAwesomeIcon icon="close" />
        </div>
      )}
    </div>
  )
}

const ViewerContent = ({ children }: ViewerPropsBase) => {
  return (
    <div className='viewer__content'>
      {children}
    </div>
  )
}

const Viewer = ({ children }: ViewerPropsBase) => {
  return (
    <div className='viewer'>
      {children}
    </div>
  )
}

Viewer.Tabs = ViewerTabs
Viewer.TabsItem = ViewerTabsItem
Viewer.Content = ViewerContent


export default Viewer
import React, { useState } from 'react'
import { icon, IconProp } from '@fortawesome/fontawesome-svg-core'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import './style.scss'
import { Dropdown } from 'react-bootstrap'

interface ViewerPropsBase {
  children: React.ReactNode
}

export interface ViewerTabsOption {
  icon?: IconProp
  name: string
  onAction: () => void
}

interface ViewerTabsProps {
  children: React.ReactNode
  hasOptions?: boolean
  options?: ViewerTabsOption[]
}

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

const CustomMenu = React.forwardRef<any, any>(
  ({ children, style, className, 'aria-labelledby': labeledBy }, ref) => {
    const [value, setValue] = useState('');

    return (
      <div
        ref={ref}
        style={style}
        className={className}
        aria-labelledby={labeledBy}
      >
        <ul className="list-unstyled">
          {React.Children.toArray(children).filter(
            (child: any) =>
              !value || child.props.children.toLowerCase().startsWith(value),
          )}
        </ul>
      </div>
    );
  },
);

const ViewerTabs = ({ children, hasOptions, options }: ViewerTabsProps) => {
  return (
    <div className='viewer__tabs'>
      <div className='d-flex'>
        {children}
      </div>
      {hasOptions && (
        <div className='viewer__tabs__options'>
          <Dropdown className='mx-3'>
            <Dropdown.Toggle as={CustomToggle}>
              <FontAwesomeIcon icon="ellipsis-vertical" />
            </Dropdown.Toggle>
            <Dropdown.Menu as={CustomMenu} className="viewer__tabs__options__menu">
              { options?.map(option => (
                <Dropdown.Item onClick={option.onAction}>
                  {option?.icon && <FontAwesomeIcon icon={option.icon} className='me-2' />}
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
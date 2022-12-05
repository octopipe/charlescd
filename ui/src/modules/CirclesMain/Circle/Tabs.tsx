import React from 'react'
import { Dropdown } from 'react-bootstrap'
import './style.scss'
import { CircleModel } from './types'

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

export enum TABS {
  CONTENT = 'content',
  TREE = 'tree'
}

interface Props {
  activeTab: TABS
  circleId: string
  circle?: CircleModel
  circleOp: string
  onChange: (tab: TABS) => void
  onClose: (circleId: string) => void
  onDelete: (circleId: string) => void
}



const CircleTabs = ({ activeTab, circleId, circle, circleOp, onChange, onClose, onDelete }: Props) => {
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
    <div className='circle__tabs'>
      <div className='circle__tabs'>
        <div
          onClick={() => onChange(TABS.CONTENT)}
          className={activeTab === TABS.CONTENT ? `circle__tabs__item--active` : `circle__tabs__item`}
        >
          <div>
            <FontAwesomeIcon icon={["far", "circle"]} className="me-2" /> {circleOp === 'C' ? 'Untitled' : circle?.name}
          </div>
          <div onClick={() => onClose(circleId)} className="circle__tabs__item__close">
            <FontAwesomeIcon icon="close" />
          </div>
        </div>
      </div>
      {circleOp !== "C" && (
        <div className='circle__tabs__options'>
          <div
            onClick={() => onChange(TABS.TREE)}
            className={activeTab === TABS.TREE ? `circle__tabs__options__item--active` : `circle__tabs__options__item`}
          >
            <div>
              <FontAwesomeIcon icon="diagram-project" className="me-2" /> Tree
            </div>
          </div>
          <Dropdown className='mx-2'>
            <Dropdown.Toggle as={CustomToggle}>
              <FontAwesomeIcon icon="ellipsis-vertical" />
            </Dropdown.Toggle>
            <Dropdown.Menu>
              <Dropdown.Item onClick={() => onDelete(circleId)}>
                Remove
              </Dropdown.Item>
            </Dropdown.Menu>
          </Dropdown>
        </div>
      )}
      
    </div>
  )
}

export default CircleTabs
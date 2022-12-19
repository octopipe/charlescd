import React from 'react'
import { Dropdown } from 'react-bootstrap'
import './style.scss'
import { ModuleModel } from '../types'

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

export enum TABS {
  CONTENT = 'content',
  TREE = 'tree'
}

interface Props {
  activeTab: TABS
  moduleId: string
  module?: ModuleModel
  moduleOp: string
  onChange: (tab: TABS) => void
  onClose: (moduleId: string) => void
  onDelete: (moduleId: string) => void
}



const ModuleTabs = ({ activeTab, moduleId, module, moduleOp, onChange, onClose, onDelete }: Props) => {
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
    <div className='module__tabs'>
      <div className='module__tabs'>
        <div
          onClick={() => onChange(TABS.CONTENT)}
          className={activeTab === TABS.CONTENT ? `module__tabs__item--active` : `module__tabs__item`}
        >
          <div>
            <FontAwesomeIcon icon="folder" className="me-2" /> {moduleOp === 'C' ? 'Untitled' : module?.name}
          </div>
          <div onClick={() => onClose(moduleId)} className="module__tabs__item__close">
            <FontAwesomeIcon icon="close" />
          </div>
        </div>
      </div>
      {moduleOp !== "C" && (
        <div className='module__tabs__options'>
          <div
            onClick={() => onChange(TABS.TREE)}
            className={activeTab === TABS.TREE ? `module__tabs__options__item--active` : `module__tabs__options__item`}
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
              <Dropdown.Item onClick={() => onDelete(moduleId)}>
                Remove
              </Dropdown.Item>
            </Dropdown.Menu>
          </Dropdown>
        </div>
      )}
      
    </div>
  )
}

export default ModuleTabs
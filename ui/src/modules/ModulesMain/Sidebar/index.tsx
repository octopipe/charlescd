import { IconProp } from '@fortawesome/fontawesome-svg-core';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import React, { useState } from 'react';
import { Link, NavLink, useMatch, useParams, useSearchParams } from 'react-router-dom';
import { ModulePagination } from '../types';
import { ReactComponent as EmptyModulesSVG } from '../../../core/assets/svg/empty-modules.svg'
import './style.scss'
import Placeholder from '../../../core/components/Placeholder';
import { Nav } from 'react-bootstrap';

interface Props {
  modules: ModulePagination
  onModuleClick: (id: string) => void
  onModuleCreateClick: () => void
}

const ModulesSidebar = ({modules, onModuleClick, onModuleCreateClick}: Props) => {
  const [searchParams] = useSearchParams();

  return (
    <div className='modules__sidebar'>
      <div className='modules__sidebar__header'>
        <Nav>
          <Nav.Item onClick={onModuleCreateClick}>
            <FontAwesomeIcon icon="plus-circle" className="me-1" /> Create module
          </Nav.Item>
        </Nav>
      </div>
      <div className='modules__sidebar__list'>
        {modules.items.length > 0 ? (
          <>
            {modules?.items.map(item => (
            <div key={item.id} className='modules__sidebar__list__item' onClick={() => onModuleClick(item.id)}>
              <FontAwesomeIcon icon={searchParams.has(item.id) ? "folder" : ["far", "folder"]} className="me-2" />
              {item.name}
            </div>
          ))}
          </>
        ) : (
          <div className='text-muted'>
            There are no modules here
          </div>
        )}
      </div>
    </div>
  )
}

export default ModulesSidebar
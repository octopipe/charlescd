import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import React, { useEffect, useState } from 'react';
import { Breadcrumb, Form } from 'react-bootstrap';
import { Link, useNavigate } from 'react-router-dom';
import { ReactComponent as LogoLight } from '../../../core/assets/svg/logo-light.svg'
import { WorkspaceModel } from '../../../modules/Workspaces/types';
import './style.scss'

export interface BreadcrumbItem {
  name: string
  to?: string
} 

export interface Props {
  workspace?: WorkspaceModel
  breadcrumbItems?: BreadcrumbItem[]
}

const Navbar = ({ workspace, breadcrumbItems }: Props) => {
  const navigate = useNavigate()

  return (
    <div className='navbar'>
      <div className='d-flex align-items-center'>
        <div className='navbar__logo' onClick={() => navigate('/')}>
          <LogoLight />
        </div>
        <div>
          <Breadcrumb className='ms-3 mt-3'>
            {breadcrumbItems && breadcrumbItems.map(breadcrumbItem => (
              <Breadcrumb.Item 
                linkAs={Link}
                linkProps={{ to: breadcrumbItem.to }}
                active={!breadcrumbItem?.to}>
                {breadcrumbItem.name}
              </Breadcrumb.Item>
            ))}
          </Breadcrumb>
        </div>
      </div>
      <div className='d-flex'>
        <FontAwesomeIcon icon="circle-user" size='xl' />
      </div>

    </div>
  )
}

export default Navbar
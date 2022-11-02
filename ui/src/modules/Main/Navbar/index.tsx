import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import React, { useState } from 'react';
import { Form } from 'react-bootstrap';
import { ReactComponent as LogoLight } from '../../../core/assets/svg/logo-light.svg'
import './style.scss'

export interface Props {
  workspaces: any[]
  onSelectWorkspace(workspaceId: string): void
}

const MainNavbar = ({ workspaces, onSelectWorkspace }: Props) => {
  return (
    <div className='main__navbar'>
      <div className='d-flex'>
        <div className='main__navbar__logo'>
          <LogoLight />
        </div>
        <div>
          <Form.Select onChange={(e) => onSelectWorkspace(e.target.value)}>
            {workspaces?.map((workspace: any) => (
              <option value={workspace?.id}>{workspace?.name}</option>
            ))}
          </Form.Select>
        </div>
      </div>
      <div className='d-flex'>
        <FontAwesomeIcon icon="circle-user" size='xl' />
      </div>

    </div>
  )
}

export default MainNavbar
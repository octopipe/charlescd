import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import React, { useEffect, useState } from 'react';
import { Form } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';
import { ReactComponent as LogoLight } from '../../../core/assets/svg/logo-light.svg'
import './style.scss'

export interface Props {
  workspaces: any[]
  selectedWorkspaceId: string
  onSelectWorkspace(workspaceId: string): void
}

const Navbar = ({ workspaces, selectedWorkspaceId, onSelectWorkspace }: Props) => {
  const navigate = useNavigate()

  return (
    <div className='navbar'>
      <div className='d-flex align-items-center'>
        <div className='navbar__logo' onClick={() => navigate('/')}>
          <LogoLight />
        </div>
        <div>
          {workspaces?.length > 0 && (
            <Form.Select defaultValue={selectedWorkspaceId} onChange={(e) => onSelectWorkspace(e.target.value)}>
              <option value="default" disabled>Select a workspace</option>
              {workspaces?.map((workspace: any) => (
                <option key={workspace?.id} value={workspace?.id}>{workspace?.name}</option>
              ))}
            </Form.Select>
          )}
        </div>
      </div>
      <div className='d-flex'>
        <FontAwesomeIcon icon="circle-user" size='xl' />
      </div>

    </div>
  )
}

export default Navbar
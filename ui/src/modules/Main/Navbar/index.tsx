import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import React, { useEffect, useState } from 'react';
import { Form } from 'react-bootstrap';
import { ReactComponent as LogoLight } from '../../../core/assets/svg/logo-light.svg'
import './style.scss'

export interface Props {
  workspaces: any[]
  onSelectWorkspace(workspaceId: string): void
}

const MainNavbar = ({ workspaces, onSelectWorkspace }: Props) => {
  const [defaultWorkspace, setDefaultWorkspace] = useState('')

  useEffect(() => {
    if (workspaces.length > 0)
      setDefaultWorkspace(workspaces[0].id)
  }, [workspaces])



  return (
    <div className='main__navbar'>
      <div className='d-flex'>
        <div className='main__navbar__logo'>
          <LogoLight />
        </div>
        <div>
          {workspaces?.length > 0 && (
            <Form.Select defaultValue={defaultWorkspace} onChange={(e) => onSelectWorkspace(e.target.value)}>
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

export default MainNavbar
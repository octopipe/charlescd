import React, { useState } from 'react';
import { Form } from 'react-bootstrap';
import './style.scss'

const MainNavbar = ({ workspaces }: any) => {
  return (
    <div className='main__navbar'>
      <div>
        <Form.Select>
          {workspaces?.map((workspace: any) => (
            <option value={workspace?.id}>{workspace?.name}</option>
          ))}
        </Form.Select>
      </div>
     
    </div>
  )
}

export default MainNavbar
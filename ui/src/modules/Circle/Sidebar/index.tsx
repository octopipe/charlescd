import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import React, { useCallback, useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import AceEditor from "react-ace";
import { ReactComponent as LogoLight } from '../../../core/assets/svg/logo-light.svg'
import CircleModules from '../../CircleModules'
import { CircleItemModule } from '../../CircleModules/types'
import { Circle } from '../types'
import './style.scss'

import "ace-builds/src-noconflict/mode-json";
import "ace-builds/src-noconflict/theme-monokai";
import { Form } from 'react-bootstrap';

interface Props {
  circle: Circle
}

const CircleSidebar = ({ circle }: Props) => {
  const { workspaceId } = useParams()

  const getModuleItems = (circle: Circle) => {
    let moduleItems: CircleItemModule[] = []
    
    for(let module of circle?.modules) {
      if (Object.keys(circle?.status).length <= 0) {
        return moduleItems
      }      

      const { modules } = circle.status      

      moduleItems.push({
        name: module.name,
        status: modules[module.name]?.status || '',
        error: modules[module.name]?.error || ''
      })
    }

    return moduleItems
  }

  return (
    <div className='circle-sidebar'>
      <div className='circle-sidebar__logo'>
        <Link to={`/workspaces/${workspaceId}/circles`} className='circle-sidebar__logo__back'>
          <FontAwesomeIcon icon="arrow-left" />
        </Link>
        <LogoLight />
      </div>
      <div className='circle-sidebar__content'>
        <div className='circle-sidebar__content__section'>
          <Form.Group className="mb-3" controlId="formBasicEmail">
            <Form.Label>Name</Form.Label>
            <Form.Control type="text" />
          </Form.Group>
          <Form.Group className="mb-3" controlId="formBasicEmail">
            <Form.Label>Description</Form.Label>
            <Form.Control as="textarea" rows={3} />
          </Form.Group>
        </div>
        <CircleModules modules={getModuleItems(circle)} />
        <div className='circle-sidebar__content__section'>
          <div className='circle-sidebar__content__title'>
            Environments
          </div>
          <AceEditor
            value={JSON.stringify(circle?.environments, null, 2)}
            width="100%"
            height='200px'
            mode="json"
            theme="monokai"
            showGutter={false}
          />
        </div>
      </div>
    </div>
  )
}

export default CircleSidebar
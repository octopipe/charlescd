import React, { useCallback, useEffect, useState } from 'react'
import { Button, Card, Col, Container, Form, FormControl, Row } from 'react-bootstrap'
import { useSelector } from 'react-redux'
import useFetch from 'use-http'
import './style.scss'
import { Circle as CircleType } from './types'
import { Link, Navigate, useNavigate, useParams } from 'react-router-dom'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import CircleModules from '../../CircleModules'
import AceEditor from "react-ace";

import "ace-builds/src-noconflict/mode-json";
import "ace-builds/src-noconflict/theme-monokai";

interface Props {
  circleId: string
}

const Circle = ({ circleId }: Props) => {
  const navigate = useNavigate()
  const { workspaceId } = useParams()
  const [circle, setCircle] = useState<CircleType>()
  const { response, get } = useFetch()

  const loadCircle = async () => {
    
    const circle = await get(`/workspaces/${workspaceId}/circles/${circleId}`)
    if (response.ok) setCircle(circle || [])
  }

  useEffect(() => {
    loadCircle()
  }, [workspaceId])

  return (
    <div className='circle'>
      <div className='circle__tabs'>

      </div>
      <div className='circle__content'>
        <div className='circle__content__title'>
          <FontAwesomeIcon icon={["far", "circle"]} className="me-2" /> {circle?.name}
        </div>
        <div className='circle__content__section'>
          <div className='circle__content__section__title'>
            <FontAwesomeIcon icon="folder" className="me-2" /> Modules
          </div>
          <CircleModules modules={circle?.status.modules || {}} />
        </div>
        <div className='circle__content__section'>
          <div className='circle__content__section__title'>
            <FontAwesomeIcon icon="folder" className="me-2" /> Environments
          </div>
          <AceEditor
            width='100%'
            height='200px'
            fontSize={16}
            mode="json"
            theme="monokai"
            value={JSON.stringify(circle?.environments, null, 2)}
          />
        </div>
      </div>
    </div>
  )
}

export default Circle
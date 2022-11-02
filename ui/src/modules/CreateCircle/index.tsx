import React, { useCallback, useEffect, useState } from 'react'
import { Button, Card, Col, Container, Form, Row } from 'react-bootstrap'
import useFetch from 'use-http'
import './style.scss'
import { CircleItem } from './types'
import AceEditor from "react-ace";
import { useNavigate, useParams } from 'react-router-dom'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Template, TEMPLATES } from '../../core/constants/templates'
import { IconProp } from '@fortawesome/fontawesome-svg-core'
import { useAppSelector } from '../../core/hooks/redux'

import "ace-builds/src-noconflict/mode-json";
import "ace-builds/src-noconflict/theme-monokai";

enum MODULE_STEPS {
  LIST = 'list',
  TEMPLATES = 'templates',
  EDIT = 'edit'
}

const initialCustomMatch = {
  headers: {
    "x-tenant-id": "1111",
  }
}

const initialEnviroments = {
  API_URL: ''
}


const CreateCircle = () => {
  const { workspaceId } = useParams()
  const { deployStrategy } = useAppSelector(store => store.main)
  const navigate = useNavigate()
  const [moduleStep, setModuleStep] = useState(MODULE_STEPS.LIST)
  const [currentTemplate, setCurrentTemplate] = useState<Template>()
  const [routing, setRouting] = useState('customMatch')
  const [customMatch, setCustomMatch] = useState(initialCustomMatch)
  const [environments, setEnvironments] = useState(initialEnviroments)
  const { response, get } = useFetch()

  useEffect(() => {

  }, [workspaceId])

  const handleTemplateUse = (template: Template) => {
    setCurrentTemplate(template)
    setModuleStep(MODULE_STEPS.EDIT)
  }

  return (
    <div className='create-circle'>
      <div className='create-circle__back' onClick={() => navigate(`/workspaces/${workspaceId}/circles`)}>
        <FontAwesomeIcon icon="arrow-left" />
      </div>
      <Row>
        <Col xs={6}>
          <Form.Group className="mb-3" controlId="formBasicEmail">
            <Form.Label>Name</Form.Label>
            <Form.Control type="text" />
          </Form.Group>
          <Form.Group className="mb-3" controlId="formBasicEmail">
            <Form.Label>Description</Form.Label>
            <Form.Control as="textarea" rows={3} />
          </Form.Group>
          <div className='create-circle__routing'>
            {deployStrategy == 'default' && (
              <>
                <Form.Group>
                  <Form.Label>Routing</Form.Label>
                  <Form.Check
                    type="radio"
                    label="Custom match"
                    value={routing}
                    checked={routing === "customMatch"}
                  />
                  <Form.Check
                    type="radio"
                    label="Segment (WIP)"
                    value={routing}
                    checked={routing === "segment"}
                  />
                </Form.Group>
                <div className='my-3'>
                  <AceEditor
                    value={JSON.stringify(customMatch, null, 2)}
                    width="100%"
                    height='150px'
                    fontSize="16px"
                    mode="json"
                    theme="monokai"
                    showGutter={false}
                  />
                </div>
              </>
            )}
            {deployStrategy == 'canary' && (
              <></>
            )}
          </div>
          <div className='create-circle__environments'>
            <div className='my-3'>
              <Form.Label>Environments</Form.Label>
              <AceEditor
                value={JSON.stringify(environments, null, 2)}
                width="100%"
                height='150px'
                fontSize="16px"
                mode="json"
                theme="monokai"
                showGutter={false}
              />
            </div>
          </div>
          

        </Col>
        <Col xs={6}>
          <div className='create-circle__modules'>
            <Form.Label>Modules</Form.Label>
            {moduleStep === MODULE_STEPS.LIST && (<div className='create-circle__modules__list'>
              <div className="d-grid gap-2">
                <Button variant="secondary" size="sm" className="circle-modules__btn-add" onClick={() => setModuleStep(MODULE_STEPS.TEMPLATES)}>
                  <FontAwesomeIcon icon="plus" />
                </Button>
              </div>
            </div>)}
            {moduleStep === MODULE_STEPS.TEMPLATES && (<div className='create-circle__modules__templates'>
              <FontAwesomeIcon icon="arrow-left" onClick={() => setModuleStep(MODULE_STEPS.LIST)}/>
              <Row>
                { TEMPLATES?.map(template => (
                  <Col xs={3}>
                    <Card className="text-center my-2">
                      <Card.Header>{template.name}</Card.Header>
                      <Card.Body>
                        <Card.Text>
                        <FontAwesomeIcon size='4x' icon={template.icon as IconProp} />
                        </Card.Text>
                        <div className="d-grid gap-2">
                          <Button onClick={() => handleTemplateUse(template)}>Use</Button>
                        </div>
                      </Card.Body>
                    </Card>
                  </Col>
                )) }
              </Row>
            </div>)}
            {moduleStep === MODULE_STEPS.EDIT && (<div className='create-circle__modules__edit'>
              <FontAwesomeIcon icon="arrow-left" onClick={() => setModuleStep(MODULE_STEPS.TEMPLATES)}/>
                <div className='my-3'>
                  {currentTemplate && currentTemplate.form?.map(input => (
                    <Form.Group>
                      <Form.Label>{input.label}</Form.Label>
                      <Form.Control type="text" />
                    </Form.Group>
                  ))}
                </div>
              <Button><FontAwesomeIcon icon="plus" /> Add module</Button>
            </div>)}
          </div>
        </Col>
      </Row>
    </div>
  )
}

export default CreateCircle
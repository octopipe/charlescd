import React, { useCallback, useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Container, Row, Col, Badge } from 'react-bootstrap';
import { Workspace, WorkspaceModel } from './types';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import WorkspaceForm from './ModalForm';
import './style.scss'
import useFetch from '../../core/hooks/fetch';
import DynamicContainer from '../../core/components/DynamicContainer';


const Workspaces = () => {
  const navigate = useNavigate()
  const [show, toggleModal] = useState(false)
  const [workspaces, setWorkspaces] = useState<WorkspaceModel[]>([])
  const { fetch, data, loading } = useFetch<WorkspaceModel[]>()


  const saveWorkspace = useCallback(async (workspace: Workspace) => {
    await fetch('/workspaces', {method: 'POST', data: workspace})
    await fetch('/workspaces')
  }, [fetch])

  useEffect(() => {
    fetch('/workspaces').then(res => setWorkspaces(res))
  }, [])  

  return (
    <DynamicContainer loading={loading} className='workspaces'>
      <div className='workspaces__background'></div>
      <Container className='workspaces__content'>
        <h2 className='mb-4'>Workspaces</h2>
        <Row>
          <Col xs={3}>
            <div className='workspaces__content__create' onClick={() => toggleModal(true)}>
              <FontAwesomeIcon icon="plus" size='2x'/>
            </div>
          </Col>
          {workspaces?.map((workspace, idx) => (
            <Col xs={3} key={idx}>
              <div className='workspaces__content__item' onClick={() => navigate(`/workspaces/${workspace.id}/circles`)}>
                <div>
                  <div>{workspace.name}</div>
                  <div className='text-muted'>{workspace.description}</div>
                  <Badge className='mt-3'>{workspace.routingStrategy}</Badge>
                </div>
                <span>
                  <span>{workspace.circles} <FontAwesomeIcon className='ms-1' icon={["far", "circle"]} /></span>
                  <span className='ms-3'>{workspace.modules} <FontAwesomeIcon className='ms-1' icon="folder" /></span>
                </span>
              </div>
            </Col>
          ))}
        </Row>
      </Container>
      <WorkspaceForm show={show} onHide={() => toggleModal(false)} onSave={saveWorkspace}/>
    </DynamicContainer>
  )
}

export default Workspaces
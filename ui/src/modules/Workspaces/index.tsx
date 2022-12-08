import React, { useCallback, useEffect, useState } from 'react';
import { generatePath, matchRoutes, Outlet, useLocation, useNavigate, useParams } from 'react-router-dom';
import useFetch from 'use-http'
import Placeholder from '../../core/components/Placeholder';
import { ReactComponent as EmptyWorkspaces } from '../../core/assets/svg/empty-workspaces.svg'
import './style.scss'
import { Container, Row, Col, Badge } from 'react-bootstrap';
import { useAppDispatch } from '../../core/hooks/redux';
import Navbar from '../../core/components/Navbar';
import { Workspace, WorkspaceModel } from './types';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import WorkspaceForm from './ModalForm';


const Workspaces = () => {
  const { workspaceId } = useParams()
  const navigate = useNavigate()
  const [show, toggleModal] = useState(false)
  const [workspaces, setWorkspaces] = useState<WorkspaceModel[]>([])
  const [selectedWorkspaceId, setSelectedWorkspaceId] = useState(workspaceId)
  const { response, get, post } = useFetch()

  const loadWorkspaces = useCallback(async () => {
    const workspaces = await get('/workspaces')
    if (response.ok) setWorkspaces(workspaces)
  }, [get, response])

  const saveWorkspace = useCallback(async (workspace: Workspace) => {
    await post('/workspaces', workspace)
    const workspaces = await get('/workspaces')
    
    if (response.ok) setWorkspaces(workspaces)

  }, [get, post, response])

  useEffect(() => {
    loadWorkspaces()
  }, [])

  return (
    <div className='workspaces'>
      <Container className='workspaces__content'>
        <h2 className='mb-4'>Workspaces</h2>
        <Row>
          <Col xs={3}>
            <div className='workspaces__content__create' onClick={() => toggleModal(true)}>
              <FontAwesomeIcon icon="plus" size='2x'/>
            </div>
          </Col>
          {workspaces?.map(workspace => (
            <Col xs={3}>
              <div className='workspaces__content__item' onClick={() => navigate(`/workspaces/${workspace.id}`)}>
                <div>
                  <div>{workspace.name}</div>
                  <div className='text-muted'>{workspace.description}</div>
                </div>
                <span>
                  <Badge>{workspace.routingStrategy}</Badge>
                </span>
              </div>
            </Col>
          ))}
        </Row>
      </Container>
      <WorkspaceForm show={show} onHide={() => toggleModal(false)} onSave={saveWorkspace}/>
    </div>
  )
}

export default Workspaces
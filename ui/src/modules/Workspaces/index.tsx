import React, { useCallback, useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Container, Row, Col, Badge } from 'react-bootstrap';
import { Workspace, WorkspaceModel } from './types';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import WorkspaceForm from './ModalForm';
import './style.scss'
import useFetch from '../../core/hooks/fetch';
import DynamicContainer from '../../core/components/DynamicContainer';
import { useAppDispatch, useAppSelector } from '../../core/hooks/redux';
import { createWorkspaceThunk, getWorkspacesThunk } from './workspacesSlice';
import { FETCH_STATUS } from '../../core/utils/fetch';


const Workspaces = () => {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()
  const workspacesState = useAppSelector(state => state.workspaces)
  const [show, toggleModal] = useState(false)

  const saveWorkspace = useCallback(async (workspace: Workspace) => {
    await dispatch(createWorkspaceThunk({ data: workspace }))
    dispatch(getWorkspacesThunk())
    toggleModal(false)
  }, [fetch])

  useEffect(() => {
    dispatch(getWorkspacesThunk())
  }, [])  

  return (
    <DynamicContainer loading={workspacesState.listStatus === FETCH_STATUS.LOADING} className='workspaces'>
      <div className='workspaces__background'></div>
      <Container className='workspaces__content'>
        <h2 className='mb-4'>Workspaces</h2>
        <Row>
          <Col xs={3}>
            <div className='workspaces__content__create' onClick={() => toggleModal(true)}>
              <FontAwesomeIcon icon="plus" size='2x'/>
            </div>
          </Col>
          {workspacesState.list?.map((workspace, idx) => (
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
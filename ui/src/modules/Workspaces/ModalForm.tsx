import React, { useCallback, useEffect, useState } from 'react';
import { Button, Form, Modal } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';
import ViewInput from '../../core/components/ViewInput';
import { Workspace } from './types';


interface Props {
  show: boolean
  onSave: (workspace: Workspace) => void
  onHide: () => void
}

const WorkspaceForm = ({ show, onHide, onSave }: Props) => {
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [routingStrategy, setRoutingStrategy] = useState('')
  
  const handleSave = () => {
    const newWorkspace = {
      name,
      description,
      routingStrategy,
    }
    console.log(newWorkspace)
    onSave(newWorkspace)
  }

  return (
    <Modal show={show} onHide={onHide}>
      <Modal.Header closeButton>
        <Modal.Title>Create workspace</Modal.Title>
      </Modal.Header>
      <Modal.Body className='p-4'>
        <ViewInput.Text
          icon="layer-group"
          label='Name'
          value={name}
          edit={true}
          canEdit={false}
          onChange={value => setName(value)}
          placeholder="Workspace name"
        />
        <ViewInput.Text
          icon="align-justify"
          label='Description'
          value={description}
          edit={true}
          canEdit={false}
          onChange={value => setDescription(value)}
          as="textarea"
          placeholder="Workspace description"
        />
        <ViewInput
          icon="signs-post"
          label='Routing strategy'
        >
          <Form.Select defaultValue="default" onChange={(e) => setRoutingStrategy(e.target.value)}>
            <option value="default" disabled>Select a routing strategy</option>
            <option value="MATCH">Circle Match</option>
            <option value="CANARY">Canary</option>
          </Form.Select>
        </ViewInput>
      </Modal.Body>
      <Modal.Footer>
        <Button onClick={handleSave}>Create workspace</Button>
        <Button variant='secondary' onClick={onHide}>Cancel</Button>
      </Modal.Footer>
    </Modal>
  )
}

export default WorkspaceForm
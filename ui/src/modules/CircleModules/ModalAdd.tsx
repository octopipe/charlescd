import React, { useEffect, useState } from 'react'
import { Button, Form, Modal, ModalProps } from 'react-bootstrap'
import { useParams } from 'react-router-dom'
import Editor from '../../core/components/Editor'
import ViewInput from '../../core/components/ViewInput'
import { useAppDispatch, useAppSelector } from '../../core/hooks/redux'
import { fetchModules } from '../Modules/modulesSlice'

const exampleOverridesValue = [
  {
    key: '$.spec.template.spec.containers[0].image',
    value: 'your-image:tag'
  }
]

const ModalAddModule = ({ show, onClose }: ModalProps) => {
  const {workspaceId} = useParams()
  const dispatch = useAppDispatch()
  const { list } = useAppSelector(state => state.modules)
  const [module, setModule] = useState<string>()

  useEffect(() => {
    dispatch(fetchModules(workspaceId))
  }, [])

  return (
    <Modal show={show} onHide={onClose}>
      <Modal.Header closeButton>
        <Modal.Title>Add module to circle</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <ViewInput
          icon="folder"
          label='Name'
        >
          <Form.Select defaultValue="default" onChange={(e) => setModule(e.target.value)}>
            <option value="default" disabled>Select a module</option>
            {list?.items?.map(item => <option value={item.id}>{item.name}</option>)}
          </Form.Select>
        </ViewInput>
        <ViewInput
          icon="folder"
          label='Override'
        >
          <Editor
            height='200px'
            value={JSON.stringify(exampleOverridesValue, null, 2)}
            onChange={() => {}}
          />
        </ViewInput>
     
      </Modal.Body>
      <Modal.Footer>
        <Button>Save</Button>
        <Button variant='secondary' onClick={onClose}>Cancel</Button>
      </Modal.Footer>
    </Modal>
  )
}

export default ModalAddModule
import React from 'react'
import { Button, Form, Modal, ModalProps } from 'react-bootstrap'
import Editor from '../../core/components/Editor'

const exampleOverridesValue = [
  {
    key: '$.spec.template.spec.containers[0].image',
    value: 'your-image:tag'
  }
]

const ModalForm = ({ show, onClose }: ModalProps) => {
  return (
    <Modal show={show} onHide={onClose}>
      <Modal.Header closeButton>
        <Modal.Title>Modal form</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <Form>
          <Form.Group className="mb-3" controlId="formBasicEmail">
            <Form.Label>Circle name</Form.Label>
            <Form.Control type="text" />
            <Form.Text className="text-muted">
              Using <a href="https://en.wiktionary.org/wiki/kebab_case">kebab-case</a> to write circle name
            </Form.Text>
          </Form.Group>
          <div>
            <Editor
              value={JSON.stringify(exampleOverridesValue, null, 2)}
              onChange={() => {}}
            />
          </div>
        </Form>
      </Modal.Body>
      <Modal.Footer>
        <Button>Save</Button>
        <Button variant='secondary' onClick={onClose}>Cancel</Button>
      </Modal.Footer>
    </Modal>
  )
}

export default ModalForm
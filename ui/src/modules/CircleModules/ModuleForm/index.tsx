import React, { useEffect, useState } from "react";
import { Button, Card, Dropdown, Modal, Form } from "react-bootstrap";
import ReactAce from "react-ace/lib/ace";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import './style.css'

import "ace-builds/src-noconflict/mode-yaml";
import "ace-builds/src-noconflict/theme-monokai";

const overrideExamples = [
  {
    key: '$.spec.template.spec.containers[0].image',
    value: 'your-image-repository/image:tag'
  }
]

const AddModule = ({ show, circle, isEdit, onSave, onClose }: any) => {
  const [modules, setModules] = useState<any[]>([])
  const [moduleRef, setModuleRef] = useState('')
  const [overrides, setOverrides] = useState(JSON.stringify(overrideExamples, null, '  '))

  useEffect(() => {
    fetch("http://localhost:8080/modules")
      .then(res => res.json())
      .then(res => setModules(res))
  }, [])

  useEffect(() => {
    if (isEdit !== '') {
      const currentModule = circle?.modules?.filter((module: any) => module?.moduleRef === isEdit)
      setModuleRef(currentModule[0].moduleRef)
      setOverrides(JSON.stringify(currentModule[0].overrides, null, '  '))
    }
  }, [isEdit])

  const AddModuleToCircle = () => {
    const newModule = {
      moduleRef,
      revision: '',
      overrides: JSON.parse(overrides)
    }

    if (circle?.modules) {
      circle.modules = [...circle?.modules, newModule]
    } else {
      circle = {
        ...circle,
        modules: [newModule]
      }
    }
  }

  const handleOnClick = () => {
    onClose()

    setModuleRef('')
    setOverrides(JSON.stringify(overrideExamples, null, '  '))
  }

  return (
    <Modal
      size="lg"
      aria-labelledby="contained-modal-title-vcenter"
      show={show}
      onHide={onClose}
      centered
    >
      <Modal.Header closeButton>
        <Modal.Title id="contained-modal-title-vcenter">
          Add module to circle
        </Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <h4>Centered Modal</h4>
        <p>
          <Form.Label>Module</Form.Label>
          <Form.Select aria-label="Default select example" value={moduleRef} onChange={value => setModuleRef(value.target.value)}>
            <option value='' disabled>Open this select menu</option>
            {modules.map((module: any) => (
              <option value={module.name}>{module.name}</option>
            ))}
          </Form.Select>
          <div className="mt-3">
            <Form.Label>Override</Form.Label>
            <ReactAce
              width="100%"
              height="200px"
              mode="yaml"
              theme="monokai"
              value={overrides}
              name="UNIQUE_ID_OF_DIV"
              onChange={value => setOverrides(value)}
              editorProps={{ $blockScrolling: true }}
            />
          </div>
        </p>
      </Modal.Body>
      <Modal.Footer>
        <Button onClick={handleOnClick}>Close</Button>
        <Button onClick={AddModuleToCircle}>Save</Button>
      </Modal.Footer>
    </Modal>
  )
}

export default AddModule
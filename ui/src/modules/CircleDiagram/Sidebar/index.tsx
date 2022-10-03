import React, { useEffect, useState } from "react";
import ReactAce from "react-ace/lib/ace";
import Form from 'react-bootstrap/Form'
import { useLocation, useMatch, useNavigate, useParams } from "react-router-dom";
import './style.css'
import { Alert, Badge, Button, Card, Dropdown } from "react-bootstrap";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import CircleModules from "../../CircleModules";

const colors = {
  "": "secondary",
  'Healthy': 'success',
  'Progressing': 'primary',
  'Degraded': 'danger'
} as any

const Sidebar = () => {
  const location = useLocation()
  const { circle: circleName } = useParams()
  const [ isEdit, toggleEdit ] = useState(false)
  const [circle, setCircle] = useState<any>()
  const navigate = useNavigate()

  useEffect(() => {
    fetch(`http://localhost:8080/circles/${circleName}`)
      .then(res => res.json())
      .then(res => setCircle(res))
  }, [location])

  useEffect(() => {
    console.log(circle)
  }, [circle])

  const handleBack = () => {
    navigate(`/circles`)
  }

  const handleToggleEdit = () => {
    toggleEdit(isEdit => !isEdit)
  }


  return (
    <div className="circle_sidebar">
      <div className="mb-3" style={{display: "flex", justifyContent: "space-between"}}>
        <FontAwesomeIcon icon="arrow-left" onClick={handleBack} />
        <FontAwesomeIcon icon={isEdit ? "eye" : "pen-to-square"} onClick={handleToggleEdit}/>
      </div>
      <div>
        <Form.Group className="mb-3" controlId="exampleForm.ControlInput1">
          <Form.Label>Name</Form.Label>
          <Form.Control
            className={isEdit ? 'circle_sidebar_input_edit' : 'circle_sidebar_input_readonly'}
            type="text"
            value={circle?.name}
            readOnly={!isEdit}
          />
        </Form.Group>
        <Form.Group className="mb-3" controlId="exampleForm.ControlTextarea1">
          <Form.Label>Description</Form.Label>
          <Form.Control
            className={isEdit ? 'circle_sidebar_input_edit' : 'circle_sidebar_input_readonly'}
            as="textarea"
            rows={3}
            value={circle?.description}
            readOnly={!isEdit}
          />
        </Form.Group>
        <CircleModules circle={circle} />
        <div className="mt-4">
          <Form.Label>Environments</Form.Label>
          <ReactAce
            width="100%"
            height="200px"
            mode="json"
            theme="monokai"
            value={JSON.stringify(circle?.environments, null, "  ")}
            name="UNIQUE_ID_OF_DIV"
            editorProps={{ $blockScrolling: true }}
          />
        </div>
        { isEdit && (
          <div className="d-grid gap-2 mt-4">
            <Button className='mt-2' variant='primary'>
              Save changes
            </Button>
          </div>
        )}
        
      </div>
    </div>
  )
}

export default Sidebar
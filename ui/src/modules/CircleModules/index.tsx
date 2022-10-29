import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import React, { useState } from "react";
import { Dropdown, Modal, Form, Button, ModalProps } from "react-bootstrap";
import ModalForm from './ModalForm'
import { CircleItemModule } from "./types";
import Alert from "../../core/components/Alert";
import './style.scss'

const ModalMoveTo = ({ show, onClose }: ModalProps) => {
  return (
    <Modal size="sm" show={show} onHide={onClose}>
      <Modal.Header closeButton>
        <Modal.Title>Move to</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <Form.Select>
          <option>Circle1</option>
        </Form.Select>
      </Modal.Body>
      <Modal.Footer>
        <Button>Move</Button>
        <Button variant="secondary" onClick={onClose}>Cancel</Button>
      </Modal.Footer>
    </Modal>
  )
}

export interface Props {
  modules: CircleItemModule[]
}

const CustomToggle = React.forwardRef<any, any>(({ children, onClick }, ref) => (
  <a
    ref={ref}
    onClick={(e) => {
      e.preventDefault();
      onClick(e);
    }}
    className="circle-modules__item__menu"
  >
    {children}
  </a>
));

const CircleModules = ({ modules }: Props) => {
  const [moveTo, toggleMoveTo] = useState(false)
  const [remove, toggleRemove] = useState(false)
  const [form, toggleForm] = useState(false)


  return (
    <>
      <div className="circle-modules">
        { modules.map(module => (
          <div className={`circle-modules__item--${module.status}`}>
            {module.name}
            <Dropdown>
              <Dropdown.Toggle as={CustomToggle}>
                <FontAwesomeIcon icon="ellipsis-vertical" />
              </Dropdown.Toggle>
              <Dropdown.Menu>
                <Dropdown.Item onClick={() => toggleForm(true)}>Edit</Dropdown.Item>
                <Dropdown.Item onClick={() => toggleMoveTo(true)}>Move to</Dropdown.Item>
                <Dropdown.Item onClick={() => toggleRemove(true)}>Remove</Dropdown.Item>
              </Dropdown.Menu>
            </Dropdown>
          </div>
        )) }
        <div className="d-grid gap-2">
          <Button variant="secondary" size="sm" className="circle-modules__btn-add" onClick={() => toggleForm(true)}>
            <FontAwesomeIcon icon="plus" />
          </Button>
        </div>
      </div>
      <ModalForm show={form} onClose={() => toggleForm(false)} />
      <ModalMoveTo show={moveTo} onClose={() => toggleMoveTo(false)}/>
      <Alert show={remove} onClose={() => toggleRemove(false)}/>
    </>
  )

}

export default CircleModules
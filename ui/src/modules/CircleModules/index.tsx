import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import React, { useEffect, useState } from "react";
import { Dropdown, Modal, Form, Button, ModalProps } from "react-bootstrap";
import ModalForm from './ModalForm'
import { CircleItemModule } from "./types";
import Alert from "../../core/components/Alert";
import './style.scss'
import { CircleItem } from "../CreateCircle/types";
import useFetch from "use-http";
import { useParams } from "react-router-dom";
import { CirclePagination } from "../CirclesMain/types";
import { CircleModule, ModuleStatus } from "../CirclesMain/Circle/types";


const ModalMoveTo = ({ show, onClose }: ModalProps) => {
  const [circles, setCircles] = useState<CirclePagination>({continue: '', items: []})
  const { response, get } = useFetch()
  const { workspaceId } = useParams()

  const loadCircles = async () => {
    const circle = await get(`/workspaces/${workspaceId}/circles`)
    if (response.ok) setCircles(circle)
  }

  useEffect(() => {
    loadCircles()
  }, [])

  return (
    <Modal size="sm" show={show} onHide={onClose}>
      <Modal.Header closeButton>
        <Modal.Title>Move to</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <Form.Select>
          {circles?.items?.map(circle => (
            <option value={circle.name}>{circle.name}</option>
          ))}
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
  modules: ModuleStatus
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
        <div className="circle-modules__title">
          Modules
        </div>
        { Object.keys(modules)?.map(moduleName => (
          <div className={modules[moduleName].status ? `circle-modules__item--${modules[moduleName].status}` : `circle-modules__item`} key={moduleName}>
            {moduleName}
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
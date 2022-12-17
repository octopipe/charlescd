import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import React, { useCallback, useEffect, useState } from "react";
import { Dropdown, Modal, Form, Button, ModalProps } from "react-bootstrap";
import ModalForm from './ModalForm'
import Alert from "../../core/components/Alert";
import './style.scss'
import { useParams } from "react-router-dom";
import { CirclePagination } from "../CirclesMain/types";
import { CircleModel } from "../CirclesMain/Circle/types";
import useFetch from "../../core/hooks/fetch";


const ModalMoveTo = ({ show, onClose }: ModalProps) => {
  const [circles, setCircles] = useState<CirclePagination>({continue: '', items: []})
  const { data, loading, fetch, error } = useFetch()
  const { workspaceId } = useParams()

  const loadCircles = useCallback(async () => {
    const circles = await fetch(`/workspaces/${workspaceId}/circles`)
    setCircles(circles)
  }, [setCircles, workspaceId])

  useEffect(() => {
    loadCircles()
  }, [loadCircles])

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
  circle: CircleModel
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

const CircleModules = ({ circle }: Props) => {
  const [moveTo, toggleMoveTo] = useState(false)
  const [remove, toggleRemove] = useState(false)
  const [form, toggleForm] = useState(false)

  return (
    <>
      <div className="circle-modules">
        <div className="circle-modules__title">
          Modules
        </div>
        { circle?.modules?.map(module => (
          <div className={circle.status.modules[module.name] ? `circle-modules__item--${circle.status.modules[module.name].status}` : `circle-modules__item`} key={module.name}>
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
      <Alert action={() => ({})} show={remove} onClose={() => toggleRemove(false)}/>
    </>
  )

}

export default CircleModules
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import React, { useCallback, useEffect, useState } from "react";
import { Dropdown, Modal, Form, Button, ModalProps } from "react-bootstrap";
import ModalAddModule from './ModalAdd'
import Alert from "../../core/components/Alert";
import './style.scss'
import { useParams } from "react-router-dom";
import { CirclePagination } from "../../core/types/circle";
import { CircleModel } from "../../core/types/circle";
import useFetch from "../../core/hooks/fetch";
import { ModuleResource } from "../../core/types/circle";


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

  const getStatus = (resources: ModuleResource[]) => {
    let status = 'Healthy'
    for (let i = 0; i < resources.length; i++) {
      if (resources[i]?.status === 'Progressing') {
        status = 'Progressing'
      }

      if (resources[i]?.status === 'Degraded') {
        status = 'Degraded'
        break
      }
    }

    console.log('STATUS', status)

    return status
  }

  return (
    <>
      <div className="circle-modules">
        <div className="circle-modules__title">
          Modules
        </div>
        { circle?.modules?.map(module => (
          <div className={`circle-modules__item--${getStatus(circle.status.modules[module.moduleId]?.resources)}`} key={module.name}>
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
      <ModalAddModule show={form} onClose={() => toggleForm(false)} />
      <ModalMoveTo show={moveTo} onClose={() => toggleMoveTo(false)}/>
      <Alert action={() => ({})} show={remove} onClose={() => toggleRemove(false)}/>
    </>
  )

}

export default CircleModules
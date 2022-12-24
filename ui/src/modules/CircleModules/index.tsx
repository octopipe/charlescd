import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import React, { useCallback, useEffect, useState } from "react";
import { Dropdown, Modal, Form, Button, ModalProps } from "react-bootstrap";
import ModalAddModule from './ModalAdd'
import Alert from "../../core/components/Alert";
import './style.scss'
import { useParams } from "react-router-dom";
import { CirclePagination, CircleStatusModel, CircleStatusModelModuleResource } from "../../core/types/circle";
import { CircleModel } from "../../core/types/circle";
import useFetch from "../../core/hooks/fetch";
import { ModuleResource } from "../../core/types/circle";
import usePolling from "../../core/hooks/polling";
import { circleApi } from "../../core/api/circle";
import Spinner from "../../core/components/Spinner";


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

interface Status {
  health: string
  message?: string
}

const CircleModules = ({ circle }: Props) => {
  const { workspaceId } = useParams()
  const [moveTo, toggleMoveTo] = useState(false)
  const [remove, toggleRemove] = useState(false)
  const [form, toggleForm] = useState(false)
  const { startPolling, stopPolling, data, loading } = usePolling<CircleStatusModel>({ timer: 3000, request: () => circleApi.getCircleStatus(workspaceId || '', circle.id) })


  useEffect(() => {
    startPolling()

    return () => stopPolling()
  }, [])

  const getStatus = (resources: CircleStatusModelModuleResource[]) => {
    let status: Status = {health: 'Healthy'}
    for (let i = 0; i < resources.length; i++) {
      if (resources[i]?.health === 'Progressing') {
        status =  { health: 'Progressing', message: resources[i]?.message}
      }

      if (resources[i]?.health === 'Degraded') {
        status =  { health: 'Degraded', message: resources[i]?.message}
        break
      }
    }

    return status
  }

  return (
    <>
      <div className="circle-modules">
        <div className="circle-modules__title">
          Modules
        </div>
        { loading ? <Spinner /> : Object.keys(data?.modules || {}).map(moduleName => (
          <div className={`circle-modules__item--${getStatus(data?.modules[moduleName]?.resources || []).health}`} key={moduleName}>
            <div className="circle-modules__item__header">
              <span>{moduleName}</span>
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
            {getStatus(data?.modules[moduleName]?.resources || [])?.message && (
              <div className="circle-modules__item__status">
                <hr />
                {getStatus(data?.modules[moduleName]?.resources || []).message}.message
              </div>
            )}
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
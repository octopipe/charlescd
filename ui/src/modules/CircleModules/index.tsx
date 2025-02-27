import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import React, { useCallback, useEffect, useState } from "react";
import { Dropdown, Modal, Form, Button, ModalProps } from "react-bootstrap";
import ModalForm from './ModalForm'
import Alert from "../../core/components/Alert";
import './style.scss'
import { useParams } from "react-router-dom";
import { CircleModule, CirclePagination, CircleStatusModel, CircleStatusModelModuleResource } from "../../core/types/circle";
import { CircleModel } from "../../core/types/circle";
import useFetch from "../../core/hooks/fetch";
import { ModuleResource } from "../../core/types/circle";
import usePolling from "../../core/hooks/polling";
import { circleApi } from "../../core/api/circle";
import Spinner from "../../core/components/Spinner";
import { CircleViewerState, fetchCircleStatus } from "../CircleViewer/circleViewerSlice";
import { useAppDispatch, useAppSelector } from "../../core/hooks/redux";
import { FETCH_STATUS } from "../../core/utils/fetch";


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
  circle?: CircleModel
  onChangeModules: (modules: CircleModule[]) => void
  onDelete: (circleModule: CircleModule) => void
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

const CircleModules = ({ circle, onChangeModules, onDelete }: Props) => {
  const { workspaceId } = useParams()
  const [moveTo, toggleMoveTo] = useState(false)
  const [remove, toggleRemove] = useState(false)
  const [form, toggleForm] = useState(false)
  const [selectedModule, setSelectedModule] = useState<CircleModule>()
  const [modules, setModules] = useState<CircleModule[]>(circle?.modules || [])
  const [circleStatus, setCircleStatus] = useState<CircleStatusModel>()
  const [draftModules, setDraftModules] = useState<CircleModule[]>([])

  const fetchCircleStatus = async (workspaceId: string, circleId: string) => {
    const res = await circleApi.getCircleStatus(workspaceId || '', circle?.id || '')
    setCircleStatus(res.data)
  }

  useEffect(() => {
    if (!circle) {
      return
    }

    fetchCircleStatus(workspaceId || '', circle?.id)
    const interval = setInterval(() => {
      fetchCircleStatus(workspaceId || '', circle?.id)
    }, 3000)

    return () => clearInterval(interval)
  }, [])

  const handleSelectModule = (module: CircleModule, cb: any) => {
    setSelectedModule(module)
    cb()
  }

  const handleAddModule = (module: CircleModule) => {
    setDraftModules(modules => ([...modules, module]))
  }

  const handleChangeModule = (module: CircleModule) => {
    let currentModules: CircleModule[] = []
    if (circle && circle?.modules) {
      currentModules = circle.modules
    }

    onChangeModules([module, ...currentModules?.filter(m => m.name !== module.name)])
  }

  useEffect(() => {
    let currentModules: CircleModule[] = []
    if (circle && circle?.modules) {
      currentModules = circle.modules
    }

    onChangeModules([...draftModules, ...currentModules])
  }, [draftModules])

  const getStatus = (resources: CircleStatusModelModuleResource[]) => {
    let status = { health: 'Healthy', message: '' }
    for (let i = 0; i < resources?.length; i++) {
      if (resources[i]?.health === 'Progressing') {
        status =  { health: 'Progressing', message: resources[i]?.message || ''}
      }

      if (resources[i]?.health === 'Degraded') {
        status =  { health: 'Degraded', message: resources[i]?.message || ''}
        break
      }
    }

    return status
  }

  const handleDelete = (module: CircleModule | undefined) => {
    if (!module)
      return
      
    onDelete(module)
    toggleRemove(false)
  }

  return (
    <>
      <div className="circle-modules">
        <div className="circle-modules__title">
          Modules
        </div>
        { draftModules?.map(module => (
          <div className='circle-modules__item' key={module?.name}>
            <div className="circle-modules__item__header">
              <span>{module.name}</span>
              <Dropdown>
                <Dropdown.Toggle as={CustomToggle}>
                  <FontAwesomeIcon icon="ellipsis-vertical" />
                </Dropdown.Toggle>
                <Dropdown.Menu>
                  <Dropdown.Item onClick={() => handleSelectModule(module, () => toggleForm(true))}>Edit</Dropdown.Item>
                </Dropdown.Menu>
              </Dropdown>
            </div>
          </div>
        )) }
        { circleStatus && circleStatus?.modules && modules?.map(module => (
          <div className={circleStatus && circleStatus?.modules ? `circle-modules__item--${getStatus(circleStatus?.modules[module?.name]?.resources)?.health}` : 'circle-modules__item'} key={module?.name}>
            <div className="circle-modules__item__header">
              <span>{module?.name}</span>
              <Dropdown>
                <Dropdown.Toggle as={CustomToggle}>
                  <FontAwesomeIcon icon="ellipsis-vertical" />
                </Dropdown.Toggle>
                <Dropdown.Menu>
                  <Dropdown.Item onClick={() => handleSelectModule(module, () => toggleForm(true))}>Edit</Dropdown.Item>
                  <Dropdown.Item onClick={() => toggleMoveTo(true)}>Move to</Dropdown.Item>
                  <Dropdown.Item onClick={() => handleSelectModule(module, () => toggleRemove(true))}>Remove</Dropdown.Item>
                </Dropdown.Menu>
              </Dropdown>
            </div>
            {circle && getStatus(circleStatus?.modules[module.name]?.resources || [])?.message && (
              <div className="circle-modules__item__status">
                <hr />
                {getStatus(circleStatus?.modules[module.name]?.resources || []).message}.message
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
      {form && <ModalForm module={selectedModule} show={true} onAdd={handleAddModule} onUpdate={handleChangeModule} onClose={() => toggleForm(false)} />}
      <ModalMoveTo show={moveTo} onClose={() => toggleMoveTo(false)}/>
      <Alert action={() => handleDelete(selectedModule)} show={remove} onClose={() => toggleRemove(false)}/>
    </>
  )

}

export default CircleModules
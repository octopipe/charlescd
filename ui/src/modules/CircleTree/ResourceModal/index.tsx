import React, { memo, useEffect, useMemo, useState } from "react";
import { Alert, Badge, ListGroup, Modal, ModalProps, Nav } from "react-bootstrap";
import { useParams } from "react-router-dom";
import Editor from "../../../core/components/Editor";
import useFetch from "../../../core/hooks/fetch";
import { Resource, ResourceMetadata } from "../types";
import './style.scss'


export interface ResourceModalProps extends ModalProps {
  circleId: string
  currentResource?: ResourceMetadata
}

enum EVENT_KEYS {
  OVERVIEW = 'overview',
  EVENTS = 'events',
  LOGS = 'logs'
}

const getAlertStatus = (resourceStatus: string) => {
  return {
    'Healthy': 'success',
    'Default': 'secondary',
  }[resourceStatus]
}


const ResourceModal = ({ circleId, show, onClose, selectedResource }: ResourceModalProps) => {
  const { workspaceId, circleName } = useParams()
  const { fetch } = useFetch()
  const [ activeKey, setActiveKey ] = useState<EVENT_KEYS>(EVENT_KEYS.OVERVIEW)
  const [resource, setResource] = useState<Resource>()
  const [events, setEvents] = useState<any>([])
  const [manifest, setManifest] = useState<any>({})
  const currentResource = useMemo(() => selectedResource, [])

  const getResource = async () => {
    const resource = await fetch(`/workspaces/${workspaceId}/circles/${circleId}/resources/${currentResource?.name}?group=${currentResource?.group || ''}&kind=${currentResource?.kind}`)
    setResource(resource || {})
  }

  const getManifest = async () => {
    const manifest = await fetch(`/workspaces/${workspaceId}/circles/${circleId}/resources/${currentResource?.name}/manifest?group=${currentResource?.group || ''}&kind=${currentResource?.kind}`)
    setManifest(manifest || {})
  }

  const getEvents = async () => {
    const events = await fetch(`/workspaces/${workspaceId}/circles/${circleId}/resources/${currentResource?.name}/events?kind=${currentResource?.kind}`)
    setEvents(events || [])
  }

  const handleSelect = (eventKey: string | null) => { setActiveKey(eventKey as EVENT_KEYS) }

  useEffect(() => {
    if (activeKey === EVENT_KEYS.OVERVIEW) {
      getResource()
      getManifest()
      return
    }

    if (activeKey === EVENT_KEYS.EVENTS) {
      getEvents()
      return
    }
  }, [activeKey])

  const Overview = () => (
    <>
      <div className="mb-3">
        <Badge><strong>Namespace: </strong>{ resource?.metadata?.namespace }</Badge>{' '}
        <Badge><strong>Kind: </strong>{ resource?.metadata?.kind }</Badge>{' '}
      </div>
      {resource?.metadata?.status && (
        <Alert variant={getAlertStatus(resource?.metadata?.status || 'Default')}>
          <strong>{resource?.metadata?.status}.</strong> {resource?.metadata?.message && <p>{resource?.metadata?.message}</p>}
        </Alert>
      )}
      
      <Editor
        value={JSON.stringify(manifest, null, 2)}
        onChange={() => {}}
        height="600px"
      />
    </>
  )

  const Events = () => (
    <>
      <ListGroup variant="flush">
        {events?.map((event: any) => (
          <ListGroup.Item
            as="li"
            className="d-flex justify-content-between align-items-start"
          >
            <div className="ms-2 me-auto">
              <div className="fw-bold">{ event?.reason }</div>
              { event?.message }
            </div>
            <Badge bg="primary" pill>
              { event?.count }
            </Badge>
          </ListGroup.Item>
        ))}
        {events.length <= 0 && (
          <p>Not found events</p>
        )}
      </ListGroup>
    </>
  )

  const Logs = () => (
    <>LOGS</>
  )

  return (
    <Modal show={show} onHide={onClose} size="xl" className="resource-modal">
      <Modal.Header closeButton>
        <Modal.Title>{currentResource?.name}</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <div className="resource-modal__header py-2">
          <Nav fill variant="pills" activeKey={activeKey} onSelect={handleSelect}>
            <Nav.Item>
              <Nav.Link eventKey={EVENT_KEYS.OVERVIEW}>Overview</Nav.Link>
            </Nav.Item>
            <Nav.Item>
              <Nav.Link eventKey={EVENT_KEYS.EVENTS}>Events</Nav.Link>
            </Nav.Item>
            <Nav.Item>
              <Nav.Link eventKey={EVENT_KEYS.LOGS}>Logs</Nav.Link>
            </Nav.Item>
          </Nav>
        </div>
        <div className="resource-modal__content">
          { activeKey === EVENT_KEYS.OVERVIEW && <Overview /> }
          { activeKey === EVENT_KEYS.EVENTS && <Events /> }
          { activeKey === EVENT_KEYS.LOGS && <Logs /> }
        </div>
      </Modal.Body>
      <Modal.Footer>
      </Modal.Footer>
    </Modal>
  )
  
}

export default React.memo(ResourceModal)
import React, { memo, useEffect, useState } from "react";
import { Alert, Badge, ListGroup, Modal, ModalProps, Nav } from "react-bootstrap";
import { Node } from "react-flow-renderer";
import { useParams } from "react-router-dom";
import useFetch from 'use-http'
import AceEditor from "react-ace";
import { Resource, ResourceMetadata } from "../types";
import './style.scss'

import "ace-builds/src-noconflict/mode-json";
import "ace-builds/src-noconflict/theme-monokai";

export interface ResourceModalProps extends ModalProps {
  node?: Node<ResourceMetadata>
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


const ResourceModal = ({ show, onClose, node }: ResourceModalProps) => {
  const { workspaceId, circleName } = useParams()
  const { response, get } = useFetch()
  const [ activeKey, setActiveKey ] = useState<EVENT_KEYS>(EVENT_KEYS.OVERVIEW)
  const [resource, setResource] = useState<Resource>()
  const [events, setEvents] = useState<any>([])

  const getResource = async () => {
    const resource = await get(`/workspaces/${workspaceId}/circles/${circleName}/resources/${node?.data.name}?group=${node?.data.group || ''}&kind=${node?.data.kind}`)
    if (response.ok) setResource(resource || {})
  }

  const getEvents = async () => {
    const events = await get(`/workspaces/${workspaceId}/circles/${circleName}/resources/${node?.data.name}/events?kind=${node?.data.kind}`)
    if (response.ok) setEvents(events || [])
  }

  const handleSelect = (eventKey: string | null) => { setActiveKey(eventKey as EVENT_KEYS) }

  useEffect(() => {
    if (activeKey === EVENT_KEYS.OVERVIEW) {
      getResource()
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
        <Badge><strong>Namespace: </strong>{ resource?.metadata.namespace }</Badge>{' '}
        <Badge><strong>Kind: </strong>{ resource?.metadata.kind }</Badge>{' '}
      </div>
      {resource?.metadata?.status && (
        <Alert variant={getAlertStatus(resource?.metadata?.status || 'Default')}>
          <strong>{resource?.metadata?.status}.</strong> {resource?.metadata?.error && <p>{resource?.metadata?.error}</p>}
        </Alert>
      )}
      
      <AceEditor
        value={JSON.stringify(resource?.manifest, null, 2)}
        width="100%"
        height='500px'
        mode="json"
        readOnly={true}
        theme="monokai"
        showGutter={false}
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
        <Modal.Title>{node?.data.name}</Modal.Title>
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

export default ResourceModal
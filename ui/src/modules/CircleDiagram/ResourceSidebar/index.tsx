import React, { useEffect, useState } from "react";
import ReactAce from "react-ace/lib/ace";
import { useLocation, useMatch, useNavigate, useParams } from "react-router-dom";
import './style.css'
import { Alert, Badge } from "react-bootstrap";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

import "ace-builds/src-noconflict/mode-json";
import "ace-builds/src-noconflict/theme-monokai";

const colors = {
  "": "secondary",
  'Healthy': 'success',
  'Progressing': 'primary',
  'Degraded': 'danger'
} as any

const Sidebar = () => {
  const location = useLocation()
  const { circle, namespace, ref, kind, resource: resourceName } = useParams()
  const [resource, setResource] = useState<any>()
  const navigate = useNavigate()

  useEffect(() => {
    const group = ref?.indexOf('/') === -1 ? "" : ref?.split('/')[0]
    const version = ref?.indexOf('/') === -1 ? ref : ref?.split('/')[1]
    fetch(`http://localhost:8080/circles/${circle}/resources/${resourceName}?group=${group}&kind=${kind}&version=${version}&namespace=${namespace}`)
      .then(res => res.json())
      .then(res => setResource(res))
  }, [location])

  const handleClose = () => {
    navigate(`/circles/${circle}`)
  }

  return (
    <div className="resource_sidebar">
      <FontAwesomeIcon
        icon="close"
        style={{position: "absolute", right: 20}}
        onClick={handleClose}
      />
      <strong>{resource?.name}</strong>
      <div>
       <Badge bg="primary">Namespace: {resource?.namespace}</Badge>{' '}
       <Badge bg="primary">Kind: {resource?.kind}</Badge>{' '}
       <Badge bg={colors[resource?.health || '']}>Health: {resource?.health}</Badge>{' '}
      </div>
      {resource?.error != "" && (
        <Alert className="mt-3" variant="danger">
          {resource?.error}
        </Alert>
      )}
      <hr />
      {resource?.resource !== "" && (
        <div className="mt-4">
          <ReactAce
            width="100%"
            mode="json"
            theme="monokai"
            value={JSON.stringify(resource?.resource, null, "  ")}
            name="UNIQUE_ID_OF_DIV"
            editorProps={{ $blockScrolling: true }}
          />
        </div>
      )}
      
    </div>
  )
}

export default Sidebar
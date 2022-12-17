import React, { useEffect, useState, memo } from "react";
import ReactFlow, { Background, ConnectionLineType, Edge, Handle, Node, Position, useEdgesState, useNodesState } from "react-flow-renderer";
import { Form, useParams } from "react-router-dom";
import dagre from 'dagre';
import './style.scss'
import { Resource, ResourceMetadata } from "./types";
import DefaultNode from "./DefaultNode";
import ResourceModal from "./ResourceModal";
import ProjectNode from "./ProjectNode";
import { Alert, Badge, ButtonGroup, FormCheck, ToggleButton } from "react-bootstrap";


interface Props {
  tree: ResourceMetadata[]
  onSelectResource: (resource: ResourceMetadata) => void
}

interface Items {
  [key: string]: ResourceMetadata
}

interface Alerts {
  [key: string]: string
}

interface Hierarchy {
  [key: string]: { chidren: string[] }
}

enum FILTERS {
  ALL = 'ALL',
  HEALTHY = 'HEALTHY',
  PROGRESSING = 'PROGRESSING',
  DEGRADED = 'DEGRAD'
}

const filters = [
  { name: 'All', value: FILTERS.ALL },
  { name: 'Healthy', value: FILTERS.HEALTHY },
  { name: 'Progressing', value: FILTERS.PROGRESSING },
  { name: 'Degraded', value: FILTERS.DEGRADED },
]

const alertStatus: Alerts = {
  "Degraded": "danger",
  "Progressing": "info",
  "Healthy": "success",
}

const TreeList = memo(({ tree, onSelectResource }: Props) => {
  const [filteredTree, setFilteredTree] = useState(tree)
  const [currentFilter, setCurrentFilter] = useState(FILTERS.ALL)

  const filterTree = (currentTree: ResourceMetadata[]) => {
    if (currentFilter === FILTERS.DEGRADED) {
      return tree.filter(item => item.status === 'Degraded') 
    }

    if (currentFilter === FILTERS.PROGRESSING) {
      return currentTree.filter(item => item.status === 'Progressing') 
    }

    if (currentFilter === FILTERS.HEALTHY) {
      return currentTree.filter(item => item.status === 'Healthy')
    }

    return tree
  }

  useEffect(() => {
    const f = filterTree(tree)
    setFilteredTree(f)
  }, [tree, currentFilter])

  

  return (
    <div className="circle-list">
      <ButtonGroup className="circle-list__filters">
        {filters.map((filter, idx) => (
          <ToggleButton
            id={`filter-${idx}`}
            key={idx}
            type="radio"
            name="circle-list-filters"
            variant="secondary"
            value={filter.value}
            checked={filter.value === currentFilter}
            onChange={(e) => setCurrentFilter(e.currentTarget.value as FILTERS)}
          >
            { filter.name }
          </ToggleButton>
        ))}
      </ButtonGroup>
      {filteredTree?.map(item => (
        <div className={item?.status ? `circle-list__item--${item?.status}` : 'circle-list__item'} onClick={() => onSelectResource(item)}>
          <div>{item.name}</div>
          <Badge className="me-2" bg="secondary">{item.kind}</Badge>
          <Badge className="me-2" bg={alertStatus[item?.status || 'Progressing']}>{item.status}</Badge>
          <Badge className="me-2" bg="secondary"><strong>Owner</strong> {`${item.owner?.kind} - ${item.owner?.name}`} </Badge>
          {item?.message && item.message !== "" && (
            <Alert variant={alertStatus[item?.status || 'Progressing']} className="mt-3">
              {item.message}
            </Alert>
          )}
        </div>
      ))}
    </div>
  )
})

export default TreeList
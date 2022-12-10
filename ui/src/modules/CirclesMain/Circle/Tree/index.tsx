import React, { useEffect, useState } from "react";
import { Edge, Node, Position} from "react-flow-renderer";
import { useParams } from "react-router-dom";
import useFetch, { CachePolicies } from 'use-http'
import dagre from 'dagre';
import './style.scss'
import { ResourceMetadata } from "./types";
import ResourceModal from "./ResourceModal";
import { ButtonGroup, ToggleButton } from "react-bootstrap";
import TreeDiagram from "./Diagram";
import TreeList from "./List";

const dagreGraph = new dagre.graphlib.Graph();
dagreGraph.setDefaultEdgeLabel(() => ({}));

const nodeWidth = 170;
const nodeHeight = 60;
const position = { x: 0, y: 0 }
const edgeType = 'smoothstep'

const getLayoutedElements = (nodes: Node[], edges: Edge[]) => {
  dagreGraph.setGraph({ rankdir: 'LR' });

  nodes.forEach((node) => {
    dagreGraph.setNode(node.id, { width: nodeWidth, height: nodeHeight });
  });

  edges.forEach((edge) => {
    dagreGraph.setEdge(edge.source, edge.target);
  });

  dagre.layout(dagreGraph);

  nodes.forEach((node) => {
    const nodeWithPosition = dagreGraph.node(node.id);
    node.targetPosition = Position.Left
    node.sourcePosition = Position.Right

    // We are shifting the dagre node position (anchor=center center) to the top left
    // so it matches the React Flow node anchor point (top left).
    node.position = {
      x: nodeWithPosition.x - nodeWidth / 2,
      y: nodeWithPosition.y - nodeHeight / 2,
    };

    return node;
  });

  return { nodes, edges };
};

interface Props {
  circleId: string
}

enum VIEWS {
  LIST = 'LIST',
  DIAGRAM = 'DIAGRAM'
}

const viewsOptions = [
  { name: 'List', icon: '', value: VIEWS.LIST },
  { name: 'Diagram', icon: '', value: VIEWS.DIAGRAM },
]

const CircleTree = ({ circleId }: Props) => {
  const { response, get } = useFetch({ cachePolicy: CachePolicies.NO_CACHE })
  const { workspaceId } = useParams()
  const [tree, setTree] = useState<ResourceMetadata[]>([])
  const [currentView, setCurrentView] = useState(VIEWS.LIST)
  const [selectedResource, setSelectedResource] = useState<ResourceMetadata | undefined>()

  const loadTree = async () => {
    const tree = await get(`/workspaces/${workspaceId}/circles/${circleId}/resources/tree`)
    if (response.ok) {
      setTree(tree || [])
    }
  }

  useEffect(() => {
    loadTree()
    const interval = setInterval(() => {
      loadTree()
    }, 3000)
    
    return () => clearInterval(interval)
  }, [])


  return (
    <div className="circle-tree">
      <div className="circle-tree__buttons">
        <ButtonGroup>
          {viewsOptions.map((option, idx) => (
            <ToggleButton
              key={idx}
              id={`radio-${idx}`}
              type="radio"
              variant='outline-primary'
              name="radio"
              value={option.value}
              checked={currentView === option.value}
              onChange={(e) => setCurrentView(e.currentTarget.value as VIEWS)}
            >
              {option.name}
            </ToggleButton>
          ))}
        </ButtonGroup>
      </div>
      <TreeList tree={tree} onSelectResource={setSelectedResource} />
      {currentView === VIEWS.DIAGRAM && <TreeDiagram show={true} tree={tree} onClose={() => setCurrentView(VIEWS.LIST)} onSelectResource={setSelectedResource} /> }
      {selectedResource && <ResourceModal show={!!selectedResource} circleId={circleId} selectedResource={selectedResource} onClose={() => setSelectedResource(undefined)}/>}
    </div>
  )
}

export default CircleTree
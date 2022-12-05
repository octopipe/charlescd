import React, { useEffect, useState, memo } from "react";
import ReactFlow, { Background, ConnectionLineType, Edge, Handle, Node, Position, useEdgesState, useNodesState } from "react-flow-renderer";
import { useParams } from "react-router-dom";
import useFetch, { CachePolicies } from 'use-http'
import dagre from 'dagre';
import './style.scss'
import { Resource, ResourceMetadata } from "./types";
import DefaultNode from "./DefaultNode";
import ResourceModal from "./ResourceModal";
import ProjectNode from "./ProjectNode";
import { Modal } from "react-bootstrap";

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

const getNodeAndEdgesByResources = (resources: ResourceMetadata[]) => {
  const nodes = resources
    .map(resource => ({
      id: `${resource.name}${resource.kind}`,
      type: resource.kind !== "Circle" && resource.kind !== "Module" ? 'default' : 'project',
      data: resource,
      position,
    }))

  const edges = resources
    .filter(resource => resource?.owner?.name !== '')
    .map(resource => ({
      id: `${resource.name}${resource.kind}-${resource.owner?.name}${resource.owner?.kind}`,
      source: `${resource.owner?.name}${resource.owner?.kind}`,
      target: `${resource.name}${resource.kind}`,
      type: edgeType,
    }))
 
 
  return {nodes, edges: [...edges]}
} 

const nodeTypes = {
  default: DefaultNode,
  project: ProjectNode,
};

interface Props {
  circleId: string
  show: boolean
  onClose: () => void
}

const CircleTree = ({ circleId, show, onClose }: Props) => {
  const { response, get } = useFetch({ cachePolicy: CachePolicies.NO_CACHE })
  const { workspaceId } = useParams()
  const [tree, setTree] = useState<ResourceMetadata[]>([])
  const [nodes, setNodes] = useNodesState([])
  const [edges, setedges] = useEdgesState([])
  const [selectedNode, setSelectedNode] = useState<Node<ResourceMetadata> | undefined>(undefined)

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

  useEffect(() => {
    if (tree.length <= 0) {
      return
    }

    const {nodes: treeNodes, edges: treeEdges} = getNodeAndEdgesByResources(tree)
    console.log('NODES', treeNodes, treeEdges)
    const {nodes: layoutedNodes, edges: layoutedEdges} = getLayoutedElements(treeNodes, treeEdges)
    
    setNodes([...layoutedNodes])
    setedges([...layoutedEdges])
  }, [tree])

  const handleNodeClick = (ev: any, node: Node<ResourceMetadata>) => {
    setSelectedNode(node)
  }

  return (
    <Modal show={show} fullscreen={true} onHide={() => onClose()} className="circle-tree">
      <Modal.Header closeButton style={{backgroundColor: "#1C1C1E"}}>
        <Modal.Title>Modal</Modal.Title>
      </Modal.Header>
      <Modal.Body style={{backgroundColor: "#1C1C1E"}}>
        <ReactFlow
          nodeTypes={nodeTypes}
          nodes={nodes}
          edges={edges}
          connectionLineType={ConnectionLineType.SmoothStep}
          nodesDraggable={false}
          nodesConnectable={false}
          onNodeClick={handleNodeClick}
          fitView
        >
          <Background/>
        </ReactFlow>
        {selectedNode && <ResourceModal show={!!selectedNode} circleId={circleId} node={selectedNode} onClose={() => setSelectedNode(undefined)}/>}
      </Modal.Body>
    </Modal>
  )
}

export default CircleTree
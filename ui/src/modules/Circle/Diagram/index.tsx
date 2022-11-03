import React, { useEffect, useState, memo } from "react";
import ReactFlow, { Background, ConnectionLineType, Edge, Handle, Node, Position, useEdgesState, useNodesState } from "react-flow-renderer";
import { useParams } from "react-router-dom";
import useFetch, { CachePolicies } from 'use-http'
import dagre from 'dagre';
import './style.scss'
import { Resource, ResourceMetadata } from "./types";
import DefaultNode from "./DefaultNode";
import ResourceModal from "./ResourceModal";

const dagreGraph = new dagre.graphlib.Graph();
dagreGraph.setDefaultEdgeLabel(() => ({}));

const nodeWidth = 200;
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
  const nodes = resources.map(resource => ({
    id: resource.name,
    type: 'default',
    data: resource,
    position,
  }))

  const edges = resources
    .filter(resource => resource?.owner)
    .map(resource => ({
      id: `${resource.name}-${resource.owner?.name}`,
      source: resource.owner?.name || '',
      target: resource.name,
      type: edgeType,
    }))

  return {nodes, edges}
} 

const nodeTypes = {
  default: DefaultNode,
};

const Diagram = () => {
  const { response, get } = useFetch({ cachePolicy: CachePolicies.NO_CACHE })
  const { workspaceId, circleName } = useParams()
  const [diagram, setDiagram] = useState<ResourceMetadata[]>([])
  const [nodes, setNodes] = useNodesState([])
  const [edges, setedges] = useEdgesState([])
  const [selectedNode, setSelectedNode] = useState<Node<ResourceMetadata> | undefined>(undefined)

  const loadDiagram = async () => {
    const circles = await get(`/workspaces/${workspaceId}/circles/${circleName}/diagram`)
    if (response.ok) setDiagram(circles || [])
  }

  useEffect(() => {
    const interval = setInterval(() => {
      loadDiagram()
      console.log('INTERBAL')
    }, 3000)
    
    return () => clearInterval(interval)
  }, [])

  useEffect(() => {
    const {nodes: diagramNodes, edges: diagramEdges} = getNodeAndEdgesByResources(diagram)
    const {nodes: layoutedNodes, edges: layoutedEdges} = getLayoutedElements(diagramNodes, diagramEdges)
    setNodes([...layoutedNodes])
    setedges([...layoutedEdges])
  }, [diagram])

  const handleNodeClick = (ev: any, node: Node<ResourceMetadata>) => {
    setSelectedNode(node)
  }

  return (
    <div className="circle-diagram">
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
        <Background />
      </ReactFlow>
      {selectedNode && <ResourceModal show={!!selectedNode} node={selectedNode} onClose={() => setSelectedNode(undefined)}/>}
    </div>
  )
}

export default Diagram
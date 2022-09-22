import React, { useCallback, useEffect, useState, memo } from 'react';
import ReactFlow, { useNodesState, useEdgesState, addEdge, Handle, Position, CoordinateExtent } from 'react-flow-renderer';
import dagre from 'dagre'
import "./style.css"
import Sidebar from './Sidebar';
import { Outlet, useNavigate } from 'react-router-dom';

const dagreGraph = new dagre.graphlib.Graph();
dagreGraph.setDefaultEdgeLabel(() => ({}));

const nodeWidth = 172;
const nodeHeight = 50;

const nodeExtent: CoordinateExtent = [
  [0, 0],
  [1000, 1000],
];

const CustomNode = memo(({ data, isConnectable }: any) => {
  const colors = {
    "": "gray",
    'Healthy': 'green',
    'Progressing': 'blue'
  } as any

  return (
    <>
      <Handle
        type="target"
        position={Position.Left}
        style={{ background: '#555' }}
        onConnect={(params) => console.log('handle onConnect', params)}
        isConnectable={isConnectable}
      />
      <div>
        <div style={{background: colors[data?.health || ''], padding: "5px"}}>{data?.kind}</div>
        <div style={{padding: "10px"}}>{data?.name}</div>
        
      </div>
      <Handle
        type="source"
        position={Position.Right}
        style={{ background: '#555' }}
        isConnectable={isConnectable}
      />
    </>
  );
});

const nodeTypes = {default: CustomNode}

const CircleDiagram = () => {
  const [resources, setResources] = useState([])
  const [nodes, setNodes, onNodesChange] = useNodesState([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);
  const navigate = useNavigate()
  const onConnect = useCallback((params: any) => setEdges((els) => addEdge(params, els)), [setEdges]);

  const onLayout = (direction: string, newNodes: any, newEdges: any) => {
    const isHorizontal = direction === 'LR';
    dagreGraph.setGraph({ rankdir: direction });

    newNodes.forEach((node: any) => {
      dagreGraph.setNode(node.id, { width: 150, height: 80 });
    });

    newEdges.forEach((edge: any) => {
      dagreGraph.setEdge(edge.source, edge.target);
    });

    dagre.layout(dagreGraph);

    const layoutedNodes = newNodes.map((node: any) => {
      const nodeWithPosition = dagreGraph.node(node.id);
      node.targetPosition = isHorizontal ? Position.Left : Position.Top;
      node.sourcePosition = isHorizontal ? Position.Right : Position.Bottom;
      // we need to pass a slightly different position in order to notify react flow about the change
      // @TODO how can we change the position handling so that we dont need this hack?
      node.position = { x: nodeWithPosition.x + Math.random() / 1000, y: nodeWithPosition.y };

      return node;
    });

    return { layoutedNodes, layoutedEdges: newEdges }
  };

  useEffect(() => {
    fetch("http://localhost:8080/circles/circle-sample")
      .then(res => res.json())
      .then(res => setResources(res))
    
    const interval = setInterval(() => {
      fetch("http://localhost:8080/circles/circle-sample")
        .then(res => res.json())
        .then(res => setResources(res))
    }, 3000)

    return () => clearInterval(interval)
  }, [])

  useEffect(() => {
    const newNodes = resources.map((res: any) => ({
      id: `${res.kind}-${res.name}`,
      sourcePosition: 'right',
      targetPosition: 'left',
      data: res,
      position: { x: 0, y: 0 },
      type: "default"
    }))

    const newEdges = resources
      .filter((res: any) => res.ownerRef != "")
      .map((res: any) => {
        const { ownerRef } = res
        return {
          id: `${ownerRef.kind}-${ownerRef.name}-${res.kind}-${res.name}`,
          source: `${ownerRef.kind}-${ownerRef.name}`,
          target: `${res.kind}-${res.name}`,
          type: 'smoothstep',
          animated: true 
        }
      })

    const { layoutedNodes, layoutedEdges } = onLayout(
      "LR",
      newNodes,
      newEdges,
    );

    setNodes(layoutedNodes)
    setEdges(layoutedEdges)
  }, [resources])

  const handleNodeClick = (ev: any, node: any) => {
    console.log(node.data)
    const { name } = node.data
    navigate(`./${name}`)
  }


  return (
    <>
      <div style={{position: "absolute", top: 0, bottom: 0, left: 0, right: 0}}>
        <ReactFlow
          nodes={nodes}
          edges={edges}
          nodeTypes={nodeTypes}
          nodeExtent={nodeExtent}
          onNodesChange={onNodesChange}
          onEdgesChange={onEdgesChange}
          onConnect={onConnect}
          onNodeClick={handleNodeClick}
          fitView
        ></ReactFlow>
      </div>
      <Outlet />
    </>
  );
};

export default CircleDiagram

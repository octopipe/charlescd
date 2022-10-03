import React, { useCallback, useEffect, useState, memo } from 'react';
import Card from 'react-bootstrap/Card';
import Button from 'react-bootstrap/Button';
import ReactFlow, { useNodesState, useEdgesState, addEdge, Handle, Position, CoordinateExtent, Background } from 'react-flow-renderer';
import dagre from 'dagre'
import "./style.css"
import { Outlet, useNavigate, useParams } from 'react-router-dom';
import Sidebar from './Sidebar';

const dagreGraph = new dagre.graphlib.Graph();
dagreGraph.setDefaultEdgeLabel(() => ({}));


const nodeExtent: CoordinateExtent = [
  [0, 0],
  [1000, 1000],
];

const CustomNode = memo(({ data, isConnectable }: any) => {
  const colors = {
    "": "secondary",
    'Healthy': 'success',
    'Progressing': 'primary',
    'Degraded': 'danger'
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
      <Card 
        border={colors[data.health]}
        text={colors[data.health] === 'light' ? 'dark' : 'white'}
        bg={colors[data.health]}
        className="text-center"
      >
        <Card.Header>{data.kind}</Card.Header>
        <Card.Body style={{background: '#fff', color: '#000'}}>
          <Card.Text>
            {data.name}
          </Card.Text>
        </Card.Body>
      </Card>
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
  const { circle } = useParams()
  const [resources, setResources] = useState([])
  const [nodes, setNodes, onNodesChange] = useNodesState([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);
  const navigate = useNavigate()
  const onConnect = useCallback((params: any) => setEdges((els) => addEdge(params, els)), [setEdges]);

  const onLayout = (direction: string, newNodes: any, newEdges: any) => {
    const isHorizontal = direction === 'LR';
    dagreGraph.setGraph({ rankdir: direction });

    newNodes.forEach((node: any) => {
      dagreGraph.setNode(node.id, { width: 220, height: 130 });
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

    fetch(`http://localhost:8080/circles/${circle}/diagram`)
      .then(res => res.json())
      .then(res => setResources(res))
    
    const interval = setInterval(() => {
      fetch(`http://localhost:8080/circles/${circle}/diagram`)
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
    const { name, ref, namespace, kind } = node.data
    navigate(`./namespaces/${namespace}/ref/${encodeURIComponent(ref)}/kind/${kind}/resource/${name}`)
  }

  


  return (
    <>
      <Sidebar />
      <div style={{position: "absolute", top: 0, bottom: 0, left: "380px", right: 0, background: "#1c1c1e"}}>
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
        >
          <Background color="#aaa" gap={16} />
        </ReactFlow>
      </div>
      <Outlet />
      
    </>
  );
};

export default CircleDiagram

import React, { memo } from "react";
import { Handle, Position, Node } from "react-flow-renderer";
import { Resource, ResourceMetadata } from "./types";

interface Props {
  data: ResourceMetadata
  isConnectable: boolean
}

const ProjectNode = memo(({ data, isConnectable }: Props) => (
  <>
    <Handle
      type="target"
      position={Position.Left}
      onConnect={(params) => console.log('handle onConnect', params)}
      style={{ background: '#555' }}
    />
    <div className={'circle-tree__item__project'}>
      <div className="circle-tree__item__header">
        { data.kind }
      </div>
      <div className="circle-tree__item__content">
        { data.name }
      </div>

    </div>
    <Handle
      type="source"
      position={Position.Right}
      onConnect={(params) => console.log('handle onConnect', params)}
      style={{ background: '#555' }}
    />
  </>
))

export default ProjectNode
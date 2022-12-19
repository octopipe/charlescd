import React, { memo } from "react";
import { Handle, Position, Node } from "react-flow-renderer";
import { Resource, ResourceMetadata } from "./types";

interface Props {
  data: ResourceMetadata
  isConnectable: boolean
}

const DefaultNode = memo(({ data, isConnectable }: Props) => (
  <>
    <Handle
      type="target"
      position={Position.Left}
      onConnect={(params) => console.log('handle onConnect', params)}
      style={{ background: '#555' }}
    />
    <div className={data?.status ? `circle-diagram__item--${data?.status}` : 'circle-diagram__item'}>
      <div className="circle-diagram__item__header">
        { data.kind }
      </div>
      <div className="circle-diagram__item__content">
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

export default DefaultNode
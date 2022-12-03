import React, { ElementType, useEffect, useState } from "react";
import Form from 'react-bootstrap/Form';
import './style.scss'

interface Props {
  value: string
  edit?: boolean
  only?: boolean
  canEdit?: boolean
  placeholder: string
  as?: string
  onChange: (value: string) => void
}

const ViewInput = (props: Props) => {
  const [isEdit, setEdit] = useState(props.edit)

  return (
    <div className="view-input">
      {isEdit || props.value === '' ? (
        <Form.Group>
          <Form.Control
            value={props.value}
            as={(props?.as) as any}
            onChange={e => props.onChange(e.target.value)}
            placeholder={props?.placeholder || ''}
          />
        </Form.Group>
      ) : (
        <div>
          <span onClick={() => props?.canEdit && setEdit(true)}>{props.value}</span>
        </div>
      )}
    </div>
  )
}


export default ViewInput
import { IconProp } from "@fortawesome/fontawesome-svg-core";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import React, { ElementType, useEffect, useState } from "react";
import Form from 'react-bootstrap/Form';
import './style.scss'

interface Props {
  value: string
  edit?: boolean
  only?: boolean
  canEdit?: boolean
  placeholder: string
  label: string
  as?: string
  icon: IconProp
  onChange: (value: string) => void
}

const ViewInput = (props: Props) => {
  const [isEdit, setEdit] = useState(props.edit)

  return (
    <div className="view-input mt-2 mb-4">
      <Form.Group>
        <Form.Label className="view-input__label">
          <FontAwesomeIcon icon={props.icon} />
          <span>{props.label}</span>
        </Form.Label>
        <Form.Control
          className="view-input__input px-2"
          value={props.value}
          as={(props?.as) as any}
          onChange={e => props.onChange(e.target.value)}
          placeholder={props?.placeholder || ''}
          plaintext={!isEdit || props.value !== ''}
          readOnly={!isEdit || props.value !== ''}
        />
      </Form.Group>
    
    </div>
  )
}


export default ViewInput
import { IconProp } from "@fortawesome/fontawesome-svg-core";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import React, { ElementType, useEffect, useState } from "react";
import Form from 'react-bootstrap/Form';
import './style.scss'

interface PropsViewInput {
  label: string
  icon: IconProp
  children: React.ReactNode
}

const ViewInput = (props: PropsViewInput) => {
  return (
    <div className="view-input mt-2 mb-4">
      <Form.Group>
        <Form.Label className="view-input__label">
          <FontAwesomeIcon icon={props.icon} />
          <span>{props.label}</span>
        </Form.Label>
        <div className="view-input__content">
          {props.children}
        </div>
      </Form.Group>
    </div>
  )
}

interface PropsViewInputText {
  value: string
  edit?: boolean
  canEdit?: boolean
  placeholder: string
  label: string
  as?: string
  icon: IconProp
  onChange: (value: string) => void
  
}

const ViewInputText = ({ label, icon, edit, value, as, placeholder, canEdit = true,  onChange }: PropsViewInputText) => {
  const [isEdit, setEdit] = useState(edit)

  return (
    <ViewInput
      label={label}
      icon={icon}
    >
      <Form.Control
        className="view-input__content__input px-2"
        value={value}
        as={as as any}
        onChange={e => onChange(e.target.value)}
        placeholder={placeholder || ''}
        plaintext={!isEdit}
        readOnly={!isEdit}
        onClick={() => canEdit && setEdit(true)}
      />
    </ViewInput>
  )
}

ViewInput.Text = ViewInputText

export default ViewInput
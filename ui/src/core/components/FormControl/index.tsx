import React from 'react'
import { Control, Controller, FieldErrors } from 'react-hook-form'
import './style.scss'

interface Props {
  name: string
  control: Control<any>
  rules?: object
  defaultValue?: any
  children: any
  errors?: any
}

const FormControl = ({ name, control, rules, defaultValue, errors, children }: Props) => {
  return (
    <Controller
      name={name}
      control={control}
      defaultValue={defaultValue}
      rules={rules}
      render={({ field }) => (
        <div className='formcontrol'>
          {React.cloneElement(children, field)}
          {errors && errors[name] && <div className='formcontrol__error'>{errors[name]?.message}</div>}
        </div>
      )}
    />
  )

}

export default FormControl
import React, { useEffect, useState } from 'react'
import { Button, Form, Modal, ModalProps } from 'react-bootstrap'
import { useForm } from 'react-hook-form'
import { useParams } from 'react-router-dom'
import Editor from '../../core/components/Editor'
import FormControl from '../../core/components/FormControl'
import ViewInput from '../../core/components/ViewInput'
import { formValidations } from '../../core/form/validate'
import { useAppDispatch, useAppSelector } from '../../core/hooks/redux'
import { CircleModule, Override } from '../../core/types/circle'
import { fetchModules } from '../Modules/modulesSlice'

const exampleOverridesValue = [
  {
    key: '$.spec.template.spec.containers[0].image',
    value: 'your-image:tag'
  }
]

interface ModalFormProps extends ModalProps {
  module?: CircleModule
  onSave: (module: CircleModule) => void
}

interface Form {
  name: string
  overrides: string
}

const ModalForm = ({ show, module, onClose, onSave }: ModalFormProps) => {
  const {workspaceId} = useParams()
  const dispatch = useAppDispatch()
  const { list } = useAppSelector(state => state.modules)
  const { control, formState: {errors}, handleSubmit } = useForm<Form>({
    defaultValues: {
      name: module?.name || '',
      overrides: JSON.stringify(module?.overrides || exampleOverridesValue, null, 2),
    }
  })

  useEffect(() => {
    dispatch(fetchModules(workspaceId))
  }, [])

  const handleSave = (form: Form) => {
    onSave({overrides: JSON.parse(form.overrides), revision: "", name: form.name})
    onClose()
  }

  return (
    <Modal show={show} onHide={onClose}>
      <Modal.Header closeButton>
        <Modal.Title>Module form</Modal.Title>
      </Modal.Header>
      <form onSubmit={handleSubmit(handleSave)}>
      <Modal.Body>
        <ViewInput
          icon="folder"
          label='Name'
        >
          <FormControl  name="name" rules={{required: 'Select a module'}} errors={errors} control={control}>
            <Form.Select>
              <option value="" disabled>Select a module</option>
              {list?.items?.map(item => <option key={item.name} value={item.name}>{item.name}</option>)}
            </Form.Select>
          </FormControl>
        </ViewInput>
        <ViewInput
          icon="folder"
          label='Override'
        >
          <FormControl name="overrides" rules={{validate: formValidations.jsonValidate}} control={control} errors={errors}>
            <Editor
              height='200px'
              readonly={false}
            />
          </FormControl>
        </ViewInput>
     
      </Modal.Body>
      <Modal.Footer>
        <Button type='submit'>Save</Button>
        <Button variant='secondary' onClick={onClose}>Cancel</Button>
      </Modal.Footer>
      </form>
    </Modal>
  )
}

export default ModalForm
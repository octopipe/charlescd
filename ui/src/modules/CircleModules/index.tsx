import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import React, { useState } from 'react'
import { Button, Form } from 'react-bootstrap'
import Module from './Module'
import AddModule from './ModuleForm'

const CircleModules = ({ circle, onChange }: any) => {
  const [showAddModule, setShowAddModule] = useState(false)
  const [isEditModule, setIsEditModule] = useState('')

  const handleRemoveModule = (name: string) => {
    const modules = circle.modules.filter((mod: any) => mod.moduleRef !== name)

    circle.modules = modules
    
    fetch("http://localhost:8080/circles/" + circle?.name, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(circle),
    })
      .then(res => res.json())
      .then(res => onChange(res))
  }

  const save = (newCircle: any) => {
    fetch("http://localhost:8080/circles/" + newCircle?.name, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(newCircle),
    })
      .then(res => res.json())
      .then(res => onChange(res))
  }

  const handleAddModuleOpen = (circle: any) => {
    setShowAddModule(true)
    setIsEditModule('')
  }

  const handleEditModuleOpen = (moduleName: any) => {
    setIsEditModule(moduleName)
    setShowAddModule(true)
  }

  const handleAddModuleClose = () => {
    setShowAddModule(false)
    setIsEditModule('')
  }

  return (
    <div>
      <Form.Label>Modules</Form.Label>
      {circle?.status && Object.keys(circle?.status?.modules || {}).map((name: any) => (
        <Module
          module={circle?.status?.modules[name]}
          name={name}
          circle={circle}
          onRemove={handleRemoveModule}
          onEdit={handleEditModuleOpen}
        />
      ))}
      <div className="d-grid gap-2">
        <Button className='mt-2' variant='secondary' style={{background: '#373739'}} onClick={handleAddModuleOpen}>
          <FontAwesomeIcon icon='plus' />{' '}Add module
        </Button>
      </div>
      <AddModule
        show={showAddModule}
        isEdit={isEditModule}
        circle={circle}
        onClose={handleAddModuleClose}
        onSave={handleAddModuleClose}
      />
    </div> 
  )
}

export default CircleModules
import React, { useEffect, useState } from "react";
import ReactAce from "react-ace/lib/ace";
import { Button, Dialog, DialogActions, DialogContent, DialogTitle, FormControl, InputLabel, MenuItem, Modal, Select } from "@mui/material";

import "ace-builds/src-noconflict/mode-yaml";
import "ace-builds/src-noconflict/theme-monokai";

const overrideExamples = [
  {
    key: '$.spec.template.spec.containers[0].image',
    value: 'your-image-repository/image:tag'
  }
]

const AddModule = ({ show, circle, isEdit, onSave, onClose }: any) => {
  const [modules, setModules] = useState<any[]>([])
  const [moduleRef, setModuleRef] = useState('')
  const [overrides, setOverrides] = useState(JSON.stringify(overrideExamples, null, '  '))

  useEffect(() => {
    fetch("http://localhost:8080/modules")
      .then(res => res.json())
      .then(res => setModules(res))
  }, [])

  useEffect(() => {
    if (isEdit !== '') {
      const currentModule = circle?.modules?.filter((module: any) => module?.moduleRef === isEdit)
      setModuleRef(currentModule[0].moduleRef)
      setOverrides(JSON.stringify(currentModule[0].overrides, null, '  '))
    }
  }, [isEdit])

  const AddModuleToCircle = () => {
    const newModule = {
      moduleRef,
      revision: '',
      overrides: JSON.parse(overrides)
    }

    if (circle?.modules) {
      circle.modules = [...circle?.modules, newModule]
    } else {
      circle = {
        ...circle,
        modules: [newModule]
      }
    }

    onSave(circle.modules)
  }

  const handleOnClick = () => {
    onClose()

    setModuleRef('')
    setOverrides(JSON.stringify(overrideExamples, null, '  '))
  }

  return (
    <Dialog
      aria-labelledby="contained-modal-title-vcenter"
      fullWidth={true}
      maxWidth="sm"
      open={show}
      onClose={onClose}
    >
      <DialogTitle id="contained-modal-title-vcenter">
        Add module to circle
      </DialogTitle>
      <DialogContent>
        <p>
          <FormControl sx={{ m: 1, minWidth: 300 }}>
            <InputLabel id="select-module">Select a module</InputLabel>
            <Select
              labelId="select-module"
              label="Select a module"
              value={moduleRef}
              onChange={value => setModuleRef(value.target.value)}
            >
              {modules.map((module: any) => (
                <MenuItem value={module.name}>{module.name}</MenuItem>
              ))}
            </Select>
          </FormControl>
          <div className="mt-3">
            <label>Override</label>
            <ReactAce
              width="100%"
              height="200px"
              mode="yaml"
              theme="monokai"
              value={overrides}
              name="UNIQUE_ID_OF_DIV"
              onChange={value => setOverrides(value)}
              editorProps={{ $blockScrolling: true }}
            />
          </div>
        </p>
      </DialogContent>
      <DialogActions>
        <Button onClick={handleOnClick}>Close</Button>
        <Button onClick={AddModuleToCircle}>Save</Button>
      </DialogActions>
    </Dialog>
  )
}

export default AddModule
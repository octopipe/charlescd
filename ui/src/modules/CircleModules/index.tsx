import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { Button, Card, CardActions, CardContent, Divider, IconButton, List, ListItem, ListItemButton, ListItemIcon, ListItemText, Menu, MenuItem, Paper, Typography } from '@mui/material'
import AddIcon from '@mui/icons-material/Add'
import FolderIcon from '@mui/icons-material/Folder'
import MoreVert from '@mui/icons-material/MoreVert'
import React, { useState } from 'react'
import Module from './Module'
import AddModule from './ModuleForm'
import { Box } from '@mui/system'

const ITEM_HEIGHT = 48;

const colors = {
  "": "secondary",
  'Healthy': '#43a047',
  'Progressing': '#039be5',
  'Degraded': '#e53935'
} as any

const CircleModules = ({ circle, onChange }: any) => {
  const [showAddModule, setShowAddModule] = useState(false)
  const [isEditModule, setIsEditModule] = useState('')
  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);
  const open = Boolean(anchorEl)

  const handleClick = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };
  const handleClose = () => {
    setAnchorEl(null);
  };

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
    handleClose()
    setIsEditModule(moduleName)
    setShowAddModule(true)
  }

  const handleAddModuleClose = (modules: any) => {
    setShowAddModule(false)
    setIsEditModule('')
    save({...circle, modules})
  }

  const handleRemove = (name: string) => {
    handleClose()
  }

  return (
    <>
      <Card variant='outlined'>
        <CardContent>
          <Typography variant="h5" component="div">Modules</Typography>
          <List>
            {!circle?.status?.modules && (
              <Typography variant="body2" color="text.secondary">
                This circle no have modules
              </Typography>
              )}
            {circle?.status && Object.keys(circle?.status?.modules || {}).map((name: any) => (
              <Paper sx={{ background: colors[circle?.status?.modules[name].status] }}>
                <ListItem sx={{mb: 2}} secondaryAction={
                  <>
                    <IconButton onClick={handleClick}>
                      <MoreVert />
                    </IconButton>
                    <Menu
                      open={open}
                      anchorEl={anchorEl}
                      onClose={handleClose}
                      PaperProps={{
                        style: {
                          maxHeight: ITEM_HEIGHT * 4.5,
                          width: '20ch',
                        },
                      }}
                    >
                      <MenuItem onClick={() => handleRemove(name)}>Remove</MenuItem>
                      <MenuItem>Move To</MenuItem>
                      <MenuItem onClick={() => handleEditModuleOpen(name)}>Edit</MenuItem>
                    </Menu>
                  </>
                }>
                  <ListItemText primary={name} />
                </ListItem>
              </Paper>
            ))}
          </List>
        </CardContent>
        <CardActions>
          <Button onClick={handleAddModuleOpen} startIcon={<AddIcon />}>
            Add module
          </Button>
        </CardActions>
      </Card> 
      <AddModule
        show={showAddModule}
        isEdit={isEditModule}
        circle={circle}
        onClose={handleAddModuleClose}
        onSave={handleAddModuleClose}
      />
    </>
  )
}

export default CircleModules
import React, { useEffect, useState } from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import './style.css'
import { Card, CardHeader, IconButton, Menu, MenuItem } from "@mui/material";

const colors = {
  "": "secondary",
  'Healthy': 'success',
  'Progressing': 'primary',
  'Degraded': 'danger'
} as any

const ITEM_HEIGHT = 48;

const Module = ({ module, name, onRemove, onEdit }: any) => {
  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);
  const open = Boolean(anchorEl)

  const handleClick = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };
  const handleClose = () => {
    setAnchorEl(null);
  };

  const handleRemove = (name: string) => {
    handleClose()
    onRemove(name)
  }

  const handleEdit = (name: string) => {
    handleClose()
    onEdit(name)
  }

  return (
    <Card
      variant='outlined'
      className="text-center"
    >
      <CardHeader
        subheader={name}
        action={
          <>
            <IconButton onClick={handleClick}>
              <FontAwesomeIcon icon="ellipsis-vertical" />
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
              <MenuItem onClick={() => handleEdit(name)}>Edit</MenuItem>
            </Menu>
          </>
        }
      />
    </Card>
  )
}

export default Module
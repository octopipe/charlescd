import React, { useEffect, useState } from "react";
import ReactAce from "react-ace/lib/ace";
import { useLocation, useMatch, useNavigate, useParams } from "react-router-dom";
import './style.css'
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import CircleModules from "../../CircleModules";
import { Box, Button, FormControl, TextField } from "@mui/material";

const colors = {
  "": "secondary",
  'Healthy': 'success',
  'Progressing': 'primary',
  'Degraded': 'danger'
} as any

const Sidebar = () => {
  const location = useLocation()
  const { circle: circleName } = useParams()
  const [ isEdit, toggleEdit ] = useState(false)
  const [circle, setCircle] = useState<any>()
  const navigate = useNavigate()

  useEffect(() => {
    fetch(`http://localhost:8080/circles/${circleName}`)
      .then(res => res.json())
      .then(res => setCircle(res))
  }, [location])

  useEffect(() => {
    console.log(circle)
  }, [circle])

  const handleBack = () => {
    navigate(`/circles`)
  }

  const handleToggleEdit = () => {
    toggleEdit(isEdit => !isEdit)
  }


  return (
    <Box sx={{mt: 7, position: "absolute", zIndex: 99, width: "350px", top: 0, bottom: 0, p: 2, background: "#121212"}}>
      <div className="mb-3" style={{display: "flex", justifyContent: "space-between"}}>
        <FontAwesomeIcon icon="arrow-left" onClick={handleBack} />
        <FontAwesomeIcon icon={isEdit ? "eye" : "pen-to-square"} onClick={handleToggleEdit}/>
      </div>
      <div>
        <Box sx={{mt: 2}}>
          <TextField
            id="outlined-basic"
            label="Name"
            variant="outlined"
            value={circle?.name}
            InputProps={{
              readOnly: !isEdit
            }}
            sx={{width: "100%"}}
          />
        </Box>
        <Box sx={{mt: 2}}>
          <TextField
            id="outlined-basic"
            label="Description"
            multiline
            rows={4}
            variant="outlined"
            value={circle?.description}
            InputProps={{
              readOnly: !isEdit
            }}
            sx={{width: "100%"}}
          />
        </Box>
        <Box sx={{mt: 2}}>
          <CircleModules circle={circle} />
        </Box>
        <Box sx={{mt: 2}}>
          <label>Environments</label>
          <ReactAce
            width="100%"
            height="200px"
            mode="json"
            theme="monokai"
            value={JSON.stringify(circle?.environments, null, "  ")}
            name="UNIQUE_ID_OF_DIV"
            editorProps={{ $blockScrolling: true }}
          />
        </Box>
        { isEdit && (
          <div className="d-grid gap-2 mt-4">
            <Button className='mt-2' variant='outlined'>
              Save changes
            </Button>
          </div>
        )}
        
      </div>
    </Box>
  )
}

export default Sidebar
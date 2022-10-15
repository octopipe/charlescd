import React, { useEffect, useState } from "react";
import ReactAce from "react-ace/lib/ace";
import { useLocation, useMatch, useNavigate, useParams } from "react-router-dom";
import './style.css'
import CloseIcon from '@mui/icons-material/Close'
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

import "ace-builds/src-noconflict/mode-json";
import "ace-builds/src-noconflict/theme-monokai";
import { Alert, AppBar, Button, Card, CardContent, Chip, Container, Dialog, IconButton, List, ListItemText, Stack, Tab, Tabs, Toolbar, Typography } from "@mui/material";
import { Box } from "@mui/system";

const colors = {
  "": "secondary",
  'Healthy': 'success',
  'Progressing': 'primary',
  'Degraded': 'error'
} as any

const eventColors = {
  "": "warning",
  'Success': 'success',
  'Pulling': 'info',
  'Failed': 'error',
  'BackOff': 'error'
} as any

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`simple-tabpanel-${index}`}
      aria-labelledby={`simple-tab-${index}`}
      {...other}
    >
      {value === index && (
        <Box sx={{ p: 3 }}>
          <Typography>{children}</Typography>
        </Box>
      )}
    </div>
  );
}

function a11yProps(index: number) {
  return {
    id: `simple-tab-${index}`,
    'aria-controls': `simple-tabpanel-${index}`,
  };
}

const Sidebar = ({ open }: any) => {
  const location = useLocation()
  const { circle, namespace, ref, kind, resource: resourceName } = useParams()
  const [resource, setResource] = useState<any>()
  const [events, setEvents] = useState([])
  const [logs, setLogs] = useState('')
  const [tabValue, setTabValue] = useState(0)
  const navigate = useNavigate()

  useEffect(() => {
    const group = ref?.indexOf('/') === -1 ? "" : ref?.split('/')[0]
    const version = ref?.indexOf('/') === -1 ? ref : ref?.split('/')[1]
    fetch(`http://localhost:8080/circles/${circle}/resources/${resourceName}?group=${group}&kind=${kind}&version=${version}&namespace=${namespace}`)
      .then(res => res.json())
      .then(res => setResource(res))

    fetch(`http://localhost:8080/circles/${circle}/events/${resourceName}?kind=${kind}`)
      .then(res => res.json())
      .then(res => setEvents(res))

    fetch(`http://localhost:8080/circles/${circle}/logs/${resourceName}`)
      .then(res => res.json())
      .then(res => setLogs(res?.logs))
  }, [location])

  const handleClose = () => {
    navigate(`/circles/${circle}`)
  }

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
  };


  return (
    <Dialog
      fullScreen
      open={true}
    >
      <AppBar sx={{ position: 'relative' }}>
        <Toolbar>
          <IconButton
            edge="start"
            color="inherit"
            onClick={handleClose}
            aria-label="close"
          >
            <CloseIcon />
          </IconButton>
          <Typography sx={{ ml: 2, flex: 1 }} variant="h6" component="div">
            {resource?.name}
          </Typography>
        </Toolbar>
      </AppBar>
      <Container sx={{m: 2}}>
        <Box sx={{mb: 2}}>
          <Chip label={`Namespace: ${resource?.namespace}`} color="primary" />{' '}
          <Chip label={`Kind: ${resource?.kind}`} color="primary" />{' '}
          <Chip label={`Health: ${resource?.health}`} color={colors[resource?.health || '']} />{' '}
        </Box>
        {resource?.error != "" && (
          <Alert severity="error">
            {resource?.error}
          </Alert>
        )}
        <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
          <Tabs value={tabValue} onChange={handleTabChange} aria-label="basic tabs example">
            <Tab label="Description" {...a11yProps(0)} />
            <Tab label="Events" {...a11yProps(1)} />
            <Tab label="Logs" {...a11yProps(2)} />
          </Tabs>
        </Box>
        <TabPanel value={tabValue} index={0}>
          {resource?.resource !== "" && (
            <div className="mt-4">
              <ReactAce
                width="100%"
                mode="json"
                theme="monokai"
                value={JSON.stringify(resource?.resource, null, "  ")}
                name="UNIQUE_ID_OF_DIV"
                editorProps={{ $blockScrolling: true }}
              />
            </div>
          )}
        </TabPanel>
        <TabPanel value={tabValue} index={1}>
          <Stack spacing={2}>
            {events?.map((event: any) => (
              <Alert severity={eventColors[event.reason]}>
                {event?.message}
              </Alert>
            ))} 
          </Stack>
        </TabPanel>
        <TabPanel value={tabValue} index={2}>
          <pre style={{background: "#000"}}>
            {logs}
          </pre>
        </TabPanel>
      </Container>
    </Dialog>
  )
}

export default Sidebar
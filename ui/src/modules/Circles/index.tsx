import React, { useEffect, useState } from 'react'
import './style.css'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Link, useNavigate } from 'react-router-dom';
import CircleModules from '../CircleModules';
import { Accordion, AccordionDetails, AccordionSummary, Box, Button, Card, CardActions, CardContent, CardHeader, Chip, CircularProgress, Container, Divider, Grid, IconButton, Typography } from '@mui/material';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import AccountTreeIcon from '@mui/icons-material/AccountTree';

const colors = {
  "": "secondary",
  'Healthy': '#43a047',
  'Progressing': '#039be5',
  'Degraded': '#e53935'
} as any

const Circles = () => {
  const [expanded, setExpanded] = useState<string | false>(false);
  const [circles, setCircles] = useState<any>([])
  const [currentCircle, setCurrentCircle] = useState<any>({})
  const [showAddModule, setShowAddModule] = useState(false)
  const navigate = useNavigate()

  const getCircleStatusByModules = (modules: any) => {
    const dangerModules = Object.keys(modules)
      .filter(module => modules[module].health !== "Healthy")
    
    return dangerModules.length <= 0 ? "Healthy" : modules[dangerModules[0]]["health"]
  }

  const handleChange =
    (panel: string) => (event: React.SyntheticEvent, isExpanded: boolean) => {
      setExpanded(isExpanded ? panel : false);
    };

  useEffect(() => {
    fetch("http://localhost:8080/circles")
      .then(res => res.json())
      .then(res => setCircles(res))

  }, [])

  const CircleStatus = ({ modules }: any) => {
    const progressing = Object.keys(modules).filter((name: any) => modules[name].status === 'Progressing')
    if (progressing.length > 0) {
      return <CircularProgress size={20} sx={{marginRight: "10px"}} />
    }

    const notHealthy = Object.keys(modules).filter((name: any) => modules[name].status !== 'Healthy')
    const currentColor = notHealthy.length > 0 ? colors['Degraded'] : colors['Healthy']

    return <FontAwesomeIcon icon="circle" size='sm' color={currentColor} style={{marginRight: "10px" }}/>

  }

  return (
    <Container>
      <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mt: 6, pb: 1 }}>
        <h1 className='text-white'>Circles</h1>
        <Button variant='contained'>
          <FontAwesomeIcon icon='plus' /> New circle
        </Button>
      </Box>
      <Divider />
      <Grid container mt={2} spacing={2}>
        {circles.map((circle: any) => (
          <Grid item xs={4}>
            <Card>
              <CardContent>
                <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                  <Typography variant="h5" component="div">
                    <CircleStatus modules={circle?.status?.modules || []} /> {circle?.name}
                  </Typography>
                  <IconButton onClick={() => navigate('/circles/' + circle.name)}><AccountTreeIcon /></IconButton>
                </Box>
                <Typography sx={{ mb: 1.5 }} color="text.secondary">
                  {circle?.description}
                </Typography>
                <CircleModules circle={circle} />
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>
    </Container>
  )
}


export default Circles
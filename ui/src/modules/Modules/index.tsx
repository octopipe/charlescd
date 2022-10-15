import React, { useEffect, useState } from 'react'
import './style.css'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Link } from 'react-router-dom';
import { Box, Button, Card, CardContent, CardHeader, Container, Divider, Stack } from '@mui/material';

const Modules = () => {
  const [modules, setModules] = useState<any>([])

  useEffect(() => {
    fetch("http://localhost:8080/modules")
      .then(res => res.json())
      .then(res => setModules(res))

  }, [])

  return (
    <Container>
      <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mt: 6, pb: 1 }}>
        <h1 className='text-white'>Modules</h1>
        <Button variant='text'>
          <FontAwesomeIcon icon='plus' /> New module
        </Button>
      </Box>
      <Divider />
      <Stack spacing={2} mt={2}>
      {modules.map((module: any) => (
        <>
          <Card>
            <CardHeader
              title={module?.name}
            />
            <CardContent>
              <div className='font-weight-light'>
                <div><strong>Repository: </strong> {module.repositoryPath}</div>
                <div><strong>Deployment path: </strong> {module.deploymentPath}</div>
                <div><strong>Secret: </strong> {module.secretRef}</div>
                <div><strong>Template: </strong> {module.templateType}</div>
              </div>
            </CardContent>
          </Card>
        </>
      ))}
      </Stack>
    </Container>
  )
}

export default Modules
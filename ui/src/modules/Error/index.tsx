import React from 'react'
import { Button, Container } from 'react-bootstrap'
import Placeholder from '../../core/components/Placeholder'
import { ReactComponent as ErrorSVG } from '../../core/assets/svg/error.svg'
import { useNavigate } from 'react-router-dom'

const Error = () => {
  const navigate = useNavigate()

  return (
    <Container>
      <div className='d-flex flex-column align-items-center'>
        <Placeholder text='Ops! An unexpected error occurred'>
          <ErrorSVG /> 
        </Placeholder>
        <Button className='my-4' style={{width: '200px'}} onClick={() => navigate(`/`)}>
          Go to home
        </Button>
      </div>
    </Container>
  )
}

export default Error
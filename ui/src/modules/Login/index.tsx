import React from 'react'
import { Button, Form } from 'react-bootstrap'
import { useNavigate } from 'react-router-dom'
import { ReactComponent as LogoLight } from '../../core/assets/svg/logo-light.svg'
import './style.scss'

const Login = () => {
  const navigate = useNavigate()

  return (
    <div className='login'>
      <div className='login__circle-1'></div>
      <div className='login__circle-2'></div>
      <div className='login__circle-3'></div>
      <div className='login__circle-4'></div>
      <div className='login__content'>
        <div className='mb-4'>
          <LogoLight />
        </div>
        <h5>Sign in with your Charles Account</h5>
        <Form.Group className='mb-3 w-100'>
          <Form.Control placeholder='Username'  />
        </Form.Group>
        <Form.Group className='mb-3 w-100'>
          <Form.Control type="password" placeholder='Passoword' />
        </Form.Group>
        <div className="d-flex justify-content-end">
          <Button variant="primary" className='mb-2 rounded-pill' onClick={() => navigate('/')}>
            Login
          </Button>
        </div>
        <hr />
        <div className="d-grid gap-2">
          <Button variant="secondary" className='rounded-pill'>
            Enter with LDAP
          </Button>
        </div>
      </div>
    </div>
  )
}

export default Login
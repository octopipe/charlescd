import React from "react";
import { Button, Modal, ModalProps } from "react-bootstrap";
import './style.scss'

interface AlertProps extends ModalProps {
  action: () => void
}

const Alert = ({ show, action, onClose }: AlertProps) => {
  return (
    <Modal className="alert" size="sm" show={show} onHide={onClose}>
      <Modal.Body>
        Are you sure?
      </Modal.Body>
      <Modal.Footer>
        <Button onClick={action}>Yes</Button>
        <Button variant="secondary" onClick={onClose}>No</Button>
      </Modal.Footer>
    </Modal>
  )
}

export default Alert
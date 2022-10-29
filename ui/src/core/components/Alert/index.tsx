import React from "react";
import { Button, Modal, ModalProps } from "react-bootstrap";
import './style.scss'

interface AlertProps extends ModalProps {

}

const Alert = ({ show, onClose }: AlertProps) => {
  return (
    <Modal className="alert" size="sm" show={show} onHide={onClose}>
      <Modal.Body>
        Are you sure?
      </Modal.Body>
      <Modal.Footer>
        <Button>Yes</Button>
        <Button variant="secondary" onClick={onClose}>No</Button>
      </Modal.Footer>
    </Modal>
  )
}

export default Alert
import React from 'react'
import AceEditor from "react-ace";

import 'ace-builds/src-noconflict/ace';
import "ace-builds/src-noconflict/mode-json";
import "ace-builds/src-noconflict/theme-monokai";

interface Props {
  value: string
  readonly?: boolean
  width?: string
  height?: string
  onChange: (value: string) => void
}

const Editor = ({ value, readonly = false, height = "100%", width = "100%", onChange }: Props) => {
  return  (
    <AceEditor
      value={value}
      width={width}
      height={height}
      mode="json"
      readOnly={readonly}
      theme="monokai"
      showGutter={false}
      onChange={onChange}
      setOptions={{
        useWorker: false
      }}
    />
  )
}

export default Editor
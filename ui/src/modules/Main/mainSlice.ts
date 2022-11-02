import { createSlice } from '@reduxjs/toolkit'

export const mainSlice = createSlice({
  name: 'main',
  initialState: {
    currentWorkspace: '',
    deployStrategy: '',
  },
  reducers: {
    setCurrentWorkspace: (state, action) => {
      state.currentWorkspace = action.payload
    },
    setDeployStrategy: (state, action) => {
      state.deployStrategy = action.payload
    }
  }
})

export const { setCurrentWorkspace, setDeployStrategy } = mainSlice.actions

export default mainSlice.reducer
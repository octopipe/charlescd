import { createSlice } from '@reduxjs/toolkit'

export const mainSlice = createSlice({
  name: 'main',
  initialState: {
    currentWorkspace: '',
    routingStrategy: '',
  },
  reducers: {
    setCurrentWorkspace: (state, action) => {
      state.currentWorkspace = action.payload
    },
    setDeployStrategy: (state, action) => {
      state.routingStrategy = action.payload
    }
  }
})

export const { setCurrentWorkspace, setDeployStrategy } = mainSlice.actions

export default mainSlice.reducer
import { createSlice } from '@reduxjs/toolkit'

export const mainSlice = createSlice({
  name: 'main',
  initialState: {
    currentWorkspace: '',
  },
  reducers: {
    setCurrentWorkspace: (state, action) => {
      state.currentWorkspace = action.payload
    }
  }
})

export const { setCurrentWorkspace } = mainSlice.actions

export default mainSlice.reducer
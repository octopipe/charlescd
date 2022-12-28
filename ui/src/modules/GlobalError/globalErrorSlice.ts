import { createAsyncThunk, createSlice, SliceCaseReducers } from '@reduxjs/toolkit'
import { FETCH_STATUS } from '../../core/utils/fetch'
import { moduleApi } from '../../core/api/module'
import { ModulePagination } from '../../core/types/module'

export const fetchModules = createAsyncThunk('modules/fetchModules', async (workspaceId?: string) => {
  const res = await moduleApi.getModules(workspaceId || '')
  return res.data
})

interface InitialState {
  code: string
  message: string
}

export const errorSlice = createSlice<InitialState, SliceCaseReducers<InitialState>, string>({
  name: 'modules',
  initialState: {
    code: '',
    message: '',
  },
  reducers: {
    setGlobalError: (state, action) => {
      state = action.payload
    },
  },
})

export const { setGlobalError } = errorSlice.actions

export default errorSlice.reducer
import { createAsyncThunk, createSlice, SliceCaseReducers } from '@reduxjs/toolkit'
import { FETCH_STATUS } from '../../core/utils/fetch'
import { moduleApi } from '../../core/api/module'
import { ModulePagination } from '../../core/types/module'

export const fetchModules = createAsyncThunk('modules/fetchModules', async (workspaceId?: string) => {
  const res = await moduleApi.getModules(workspaceId || '')
  return res.data
})

interface InitialState {
  list: ModulePagination
  status: FETCH_STATUS
  error: string | undefined
}

export const modulesSlice = createSlice<InitialState, SliceCaseReducers<InitialState>, string>({
  name: 'modules',
  initialState: {
    list: {} as ModulePagination,
    status: FETCH_STATUS.LOADING,
    error: ''
  },
  reducers: {
    setCircles: (state, action) => {
      state.list = action.payload
    },
  },
  extraReducers(builder) {
    builder
      .addCase(fetchModules.pending, (state, action) => {
        state.status = FETCH_STATUS.LOADING
      })
      .addCase(fetchModules.fulfilled, (state, action) => {
        state.status = FETCH_STATUS.SUCCEEDED
        state.list = action.payload
      })
      .addCase(fetchModules.rejected, (state, action) => {
        state.status = FETCH_STATUS.FAILED
        state.error = action.error.message
      })
  }
})

export const { setCircles } = modulesSlice.actions

export default modulesSlice.reducer
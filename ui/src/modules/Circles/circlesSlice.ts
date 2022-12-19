import { createAsyncThunk, createSlice, SliceCaseReducers } from '@reduxjs/toolkit'
import { circleApi } from '../../core/api/circle'
import { FETCH_STATUS } from '../../core/utils/fetch'
import { CircleItem, CirclePagination } from '../../core/types/circle'

export const fetchCircles = createAsyncThunk('circles/fetchCircles', async (workspaceId?: string) => {
  const res = await circleApi.getCircles(workspaceId || '')
  return res.data
})

interface InitialState {
  list: CirclePagination
  status: FETCH_STATUS
  error: string | undefined
}

export const circlesSlice = createSlice<InitialState, SliceCaseReducers<InitialState>, string>({
  name: 'circles',
  initialState: {
    list: {} as CirclePagination,
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
      .addCase(fetchCircles.pending, (state, action) => {
        state.status = FETCH_STATUS.LOADING
      })
      .addCase(fetchCircles.fulfilled, (state, action) => {
        state.status = FETCH_STATUS.SUCCEEDED
        state.list = action.payload
      })
      .addCase(fetchCircles.rejected, (state, action) => {
        state.status = FETCH_STATUS.FAILED
        state.error = action.error.message
      })
  }
})

export const { setCircles } = circlesSlice.actions

export default circlesSlice.reducer
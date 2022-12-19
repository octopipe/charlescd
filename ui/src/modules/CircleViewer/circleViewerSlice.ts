import { createAsyncThunk, createSlice, SliceCaseReducers } from '@reduxjs/toolkit'
import { circleApi } from '../../core/api/circle'
import { FETCH_STATUS } from '../../core/utils/fetch'
import { CircleModel } from '../../core/types/circle'

interface ThunkArg {
  workspaceId?: string
  circleId?: string
}

export const fetchCircle = createAsyncThunk('circleViewer/fetchCircle', async (arg: ThunkArg) => {
  const res = await circleApi.getCircle(arg.workspaceId || '', arg.circleId || '')
  return res.data
})

interface ViewerState {
  item: CircleModel
  status: FETCH_STATUS
  error: string | undefined
}

interface InitialState {
  [key: string]: ViewerState
}

export const circleViewerSlice = createSlice<InitialState, SliceCaseReducers<InitialState>, string>({
  name: 'circleViewer',
  initialState: {} as InitialState,
  reducers: {
    removeCircleViewer: (state, action) => {
      const idToRemove = action.payload.circleId
      delete state[idToRemove]
    },

  },
  extraReducers(builder) {
    builder
      .addCase(fetchCircle.pending, (state, action) => {
        const { meta } = action
        state[meta.arg?.circleId || ''] = {...state[meta.arg?.circleId || ''], status: FETCH_STATUS.LOADING}
      })
      .addCase(fetchCircle.fulfilled, (state, action) => {
        const { meta } = action
        state[meta.arg?.circleId || ''] = {...state[meta.arg?.circleId || ''], status: FETCH_STATUS.SUCCEEDED, item: action.payload}
      })
      .addCase(fetchCircle.rejected, (state, action) => {
      })
  }
})

export const { removeCircleViewer } = circleViewerSlice.actions

export default circleViewerSlice.reducer
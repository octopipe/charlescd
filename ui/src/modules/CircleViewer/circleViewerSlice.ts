import { createAsyncThunk, createSlice, SliceCaseReducers } from '@reduxjs/toolkit'
import { circleApi } from '../../core/api/circle'
import { FETCH_STATUS } from '../../core/utils/fetch'
import { Circle, CircleModel, CircleStatusModel } from '../../core/types/circle'
import { toast } from 'react-toastify'

interface ThunkArg {
  workspaceId?: string
  circleId?: string
  data?: Circle
}

export const fetchCircle = createAsyncThunk('circleViewer/fetchCircle', async (arg: ThunkArg) => {
  const res = await circleApi.getCircle(arg.workspaceId || '', arg.circleId || '')
  return res.data
})

export const fetchCircleStatus = createAsyncThunk('circleViewer/fetchCircleStatus', async (arg: any) => {
  const res = await circleApi.getCircleStatus(arg.workspaceId, arg.circleId)
  return res.data
})

export const fetchCircleSync = createAsyncThunk('circleViewer/fetchCircleSync', async (arg: any) => {
  const res = await circleApi.sync(arg.workspaceId, arg.circleId)
  return res.data
})

export const fetchCircleCreate = createAsyncThunk('circleViewer/fetchCircleCreate', async (arg: any) => {
  const res = await circleApi.createCircle(arg.workspaceId, arg.data )
  return res.data
})

export const fetchCircleUpdate = createAsyncThunk('circleViewer/fetchCircleUpdate', async (arg: any) => {
  const res = await circleApi.updateCircle(arg.workspaceId, arg.circleId, arg.data )
  return res.data
})



export interface CircleViewerState {
  item: { data: CircleModel, status: FETCH_STATUS }
  status: { data: CircleStatusModel, status: FETCH_STATUS, syncedAt: string }
}

export interface ViewerState {
  [key: string]: CircleViewerState
}

const getCurrentTime = () => {
  const date = new Date()
  const hour = '0' + date.getHours()
  const minutes = '0' + date.getMinutes()
  const seconds = date.getSeconds()

  return `${hour.slice(-2)}:${minutes.slice(-2)}:${seconds}`
}

export const circleViewerSlice = createSlice<ViewerState, SliceCaseReducers<ViewerState>, string>({
  name: 'circleViewer',
  initialState: {} as ViewerState,
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
        state[meta.arg?.circleId || ''] = {...state[meta.arg?.circleId || ''], item: {...state[meta.arg?.circleId || '']?.item, status: FETCH_STATUS.LOADING}}
      })
      .addCase(fetchCircle.fulfilled, (state, action) => {
        const { meta } = action
        state[meta.arg?.circleId || ''] = {...state[meta.arg?.circleId || ''], item: {data: action.payload, status: FETCH_STATUS.SUCCEEDED}}
      })
      .addCase(fetchCircle.rejected, (state, action) => {
      })
      .addCase(fetchCircleStatus.pending, (state, action) => {
        const { meta } = action
        state[meta.arg?.circleId || ''] = {...state[meta.arg?.circleId || ''], status: {...state[meta.arg?.circleId || ''].status, status: FETCH_STATUS.LOADING}}
      })
      .addCase(fetchCircleStatus.fulfilled, (state, action) => {
        const { meta } = action
        state[meta.arg?.circleId || ''] = {...state[meta.arg?.circleId || ''], status: {data: action.payload, status: FETCH_STATUS.SUCCEEDED, syncedAt: getCurrentTime()}}
      })
      .addCase(fetchCircleStatus.rejected, (state, action) => {
      })
      .addCase(fetchCircleCreate.fulfilled, (state, action) => {
        toast.success(`Circle ${action?.payload?.name} create with success!`)
      })
      .addCase(fetchCircleCreate.rejected, (state, action) => {
        toast.error(`Failed to create circle: ${action.error.message}`)
      })
      .addCase(fetchCircleUpdate.fulfilled, (state, action) => {
        toast.success(`Circle ${action?.payload?.name} updated with success!`)
      })
      .addCase(fetchCircleUpdate.rejected, (state, action) => {
        toast.error(`Failed to update circle: ${action.error.message}`)
      })
  }
})

export const { removeCircleViewer } = circleViewerSlice.actions

export default circleViewerSlice.reducer
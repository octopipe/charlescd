import { createAsyncThunk, createSlice, SliceCaseReducers } from '@reduxjs/toolkit'
import { circleApi } from '../../core/api/circle'
import { FETCH_STATUS } from '../../core/utils/fetch'
import { Circle, CircleModel, CircleStatusModel } from '../../core/types/circle'
import { WorkspaceModel, WorkspacePagination } from '../../core/types/workspace'
import { workspaceApi } from '../../core/api/workspace'

interface ThunkArg {
  workspaceId?: string
  circleId?: string
  data?: Circle
}

export const getWorkspacesThunk = createAsyncThunk('workspaces/getWorkspacesThunk', async () => {
  const res = await workspaceApi.getWorkspaces()
  return res.data
})

export const createWorkspaceThunk = createAsyncThunk('workspaces/createWorkspaceThunk', async (arg: any) => {
  const res = await workspaceApi.createWorkspace(arg.data)
  return res.data
})

interface InitialState {
  list: WorkspaceModel[]
  listStatus: FETCH_STATUS
  createStatus: FETCH_STATUS
  error: string | undefined
}

export const circleViewerSlice = createSlice<InitialState, SliceCaseReducers<InitialState>, string>({
  name: 'circleViewer',
  initialState: {
    list: [] as WorkspaceModel[],
    listStatus: FETCH_STATUS.LOADING,
    createStatus: FETCH_STATUS.SUCCEEDED,
    error: ''
  },
  reducers: {
  },
  extraReducers(builder) {
    builder
      .addCase(getWorkspacesThunk.pending, (state, action) => {
        state.listStatus = FETCH_STATUS.LOADING
      })
      .addCase(getWorkspacesThunk.fulfilled, (state, action) => {
        state.listStatus = FETCH_STATUS.SUCCEEDED
        state.list = action.payload
      })
      .addCase(getWorkspacesThunk.rejected, (state, action) => {
      })
      .addCase(createWorkspaceThunk.pending, (state, action) => {
        state.createStatus = FETCH_STATUS.LOADING
      })
      .addCase(createWorkspaceThunk.fulfilled, (state, action) => {
        state.createStatus = FETCH_STATUS.LOADING
      })
      .addCase(createWorkspaceThunk.rejected, (state, action) => {
      })
  }
})

export const { removeCircleViewer } = circleViewerSlice.actions

export default circleViewerSlice.reducer
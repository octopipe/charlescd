import { createSlice } from '@reduxjs/toolkit'
import { BreadcrumbItem } from '../../core/components/Navbar'
import { WorkspaceModel } from '../Workspaces/types'

export const mainSlice = createSlice({
  name: 'main',
  initialState: {
    workspace: {} as WorkspaceModel,
    breadcrumbItems: [] as BreadcrumbItem[],
  },
  reducers: {
    setWorkspace: (state, action) => {
      state.workspace = action.payload
    },
    setBreadcrumbItems: (state, action) => {
      let breadcrumbItems = [
        { name: 'Workspaces', to: '/' },
        { name: state.workspace?.name || '' },
      ]

      breadcrumbItems = [...breadcrumbItems, ...action.payload]
      state.breadcrumbItems = breadcrumbItems
    }
  }
})

export const { setWorkspace, setBreadcrumbItems } = mainSlice.actions

export default mainSlice.reducer
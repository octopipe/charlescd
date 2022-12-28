import { configureStore } from '@reduxjs/toolkit'
import mainReducer from './modules/Main/mainSlice'
import workspacesReducer from './modules/Workspaces/workspacesSlice'
import circlesReducer from './modules/Circles/circlesSlice'
import circleViewerReducer from './modules/CircleViewer/circleViewerSlice'
import modulesReducer from './modules/Modules/modulesSlice'
import globalErrorReducer from './modules/GlobalError/globalErrorSlice'
import { handleErrorsMiddleware } from './core/api/middlewares'

const store = configureStore({
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware().concat(handleErrorsMiddleware),
  reducer: {
    main: mainReducer,
    workspaces: workspacesReducer,
    circles: circlesReducer,
    circleViewer: circleViewerReducer,
    modules: modulesReducer,
    globalError: globalErrorReducer,
  },
})


export type RootState = ReturnType<typeof store.getState>

export type AppDispatch = typeof store.dispatch

export default store
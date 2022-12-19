import { configureStore } from '@reduxjs/toolkit'
import mainReducer from './modules/Main/mainSlice'
import circlesReducer from './modules/Circles/circlesSlice'
import circleViewerReducer from './modules/CircleViewer/circleViewerSlice'
import modulesReducer from './modules/Modules/modulesSlice'

const store = configureStore({
  reducer: {
    main: mainReducer,
    circles: circlesReducer,
    circleViewer: circleViewerReducer,
    modules: modulesReducer,
  },
})


export type RootState = ReturnType<typeof store.getState>

export type AppDispatch = typeof store.dispatch

export default store
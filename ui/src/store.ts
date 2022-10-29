import { configureStore } from '@reduxjs/toolkit'
import mainReducer from './modules/Main/mainSlice'

const store = configureStore({
  reducer: {
    main: mainReducer,
  },
})


export type RootState = ReturnType<typeof store.getState>

export type AppDispatch = typeof store.dispatch

export default store
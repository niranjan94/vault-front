import { createSlice, PayloadAction } from '@reduxjs/toolkit'


export type AuthState = {
  readonly isAuthenticated: boolean;
  readonly token?: string | null;
};

export const initialState: AuthState = {
  isAuthenticated: true,
  token: null
};

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    login(state, action: PayloadAction<string>) {
      Object.assign( state, {
        isAuthenticated: true,
        token: action.payload
      })
    },
    logout(state) {
      Object.assign(state, {
        isAuthenticated: false,
        token: null
      })
    },
  }
})

export const { login, logout } = authSlice.actions
export default authSlice.reducer

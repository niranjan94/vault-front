import { RootState } from '../reducers';

export const isAuthenticatedSelector = (state: RootState) => state.auth.isAuthenticated;

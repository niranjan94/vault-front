import { combineReducers } from 'redux';
import auth from './auth/slice';

const rootReducer = combineReducers({
  auth,
});

export type RootState = ReturnType<typeof rootReducer>
export default rootReducer;

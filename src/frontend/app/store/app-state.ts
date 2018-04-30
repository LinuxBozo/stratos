import { ActionHistoryState } from './reducers/action-history-reducer';
import { AuthState } from './reducers/auth.reducer';
import { DashboardState } from './reducers/dashboard-reducer';
import { ListsState } from './reducers/list.reducer';
import { CreateNewApplicationState } from './types/create-application.types';
import { DeployApplicationState } from './types/deploy-application.types';
import { EndpointState } from './types/endpoint.types';
import { IRequestDataState, IRequestState } from './types/entity.types';
import { PaginationState } from './types/pagination.types';
import { RoutingHistory } from './types/routing.type';
import { UAASetupState } from './types/uaa-setup.types';
import { UserProfileInfo } from './types/user-profile.types';
import { InternalEventsState } from './types/internal-events.types';

export interface IRequestTypeState {
  [entityKey: string]: IRequestEntityTypeState<any>;
}
export interface IRequestEntityTypeState<T> {
  [guid: string]: T;
}
export interface AppState {
  actionHistory: ActionHistoryState;
  auth: AuthState;
  uaaSetup: UAASetupState;
  endpoints: EndpointState;
  pagination: PaginationState;
  request: IRequestState;
  requestData: IRequestDataState;
  dashboard: DashboardState;
  createApplication: CreateNewApplicationState;
  deployApplication: DeployApplicationState;
  lists: ListsState;
  routing: RoutingHistory;
  internalEvents: InternalEventsState;
}

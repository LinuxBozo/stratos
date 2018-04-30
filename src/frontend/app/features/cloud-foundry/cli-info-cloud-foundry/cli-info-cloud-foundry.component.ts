import { Component, OnInit, ViewEncapsulation, Optional } from '@angular/core';
import { AppState } from '../../../store/app-state';
import { Store } from '@ngrx/store';
import { getPreviousRoutingState, RoutingEvent } from '../../../store/types/routing.type';
import { first, map, filter } from 'rxjs/operators';
import { RouterNav } from '../../../store/actions/router.actions';
import { ActivatedRoute } from '@angular/router';
import { EntityServiceFactory } from '../../../core/entity-service-factory.service';
import { APIResource, EntityInfo } from '../../../store/types/api.types';
import { endpointSchemaKey, entityFactory } from '../../../store/helpers/entity-factory';
import { GetAllEndpoints } from '../../../store/actions/endpoint.actions';
import { combineLatest } from 'rxjs/observable/combineLatest';
import { getFullEndpointApiUrl } from '../../endpoints/endpoint-helpers';
import { Observable } from 'rxjs/Observable';
import { IApp, IOrganization, ISpace } from '../../../core/cf-api.types';
import { IHeaderBreadcrumb } from '../../../shared/components/page-header/page-header.types';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { CFAppCLIInfoContext } from '../../../shared/components/cli-info/cli-info.component';
import { ActiveRouteCfOrgSpace } from '../cf-page.types';
import { CloudFoundryEndpointService } from '../services/cloud-foundry-endpoint.service';
import { CloudFoundryOrganizationService } from '../services/cloud-foundry-organization.service';
import { EndpointModel } from '../../../store/types/endpoint.types';
import { CloudFoundrySpaceService } from '../services/cloud-foundry-space.service';
import { getActiveRouteCfOrgSpaceProvider } from '../cf.helpers';
import { CfUserService } from '../../../shared/data-services/cf-user.service';

@Component({
  selector: 'app-cli-info-cloud-foundry',
  templateUrl: './cli-info-cloud-foundry.component.html',
  styleUrls: ['./cli-info-cloud-foundry.component.scss'],
  providers: [
    getActiveRouteCfOrgSpaceProvider,
    CloudFoundryEndpointService,
    CloudFoundryOrganizationService,
    CloudFoundrySpaceService
  ]
})
export class CliInfoCloudFoundryComponent implements OnInit {

  cfEndpointEntityService: any;
  public previousUrl: string;
  public previousQueryParams: {
    [key: string]: string;
  };

  public context$: Observable<CFAppCLIInfoContext>;
  public breadcrumbs$: Observable<IHeaderBreadcrumb[]>;
  public route$: Observable<{ url: string, queryParams: any }>;

  public endpointOrgSpace$: Observable<[
    EntityInfo<EndpointModel>,
    EntityInfo<APIResource<IOrganization>>,
    EntityInfo<APIResource<ISpace>>
  ]>;

  constructor(
    private store: Store<AppState>,
    private activeRouteCfOrgSpace: ActiveRouteCfOrgSpace,
    private cfEndpointService: CloudFoundryEndpointService,
    @Optional() private cfOrgService: CloudFoundryOrganizationService,
    @Optional() private cfSpaceService: CloudFoundrySpaceService
  ) {
    this.breadcrumbs$ = new BehaviorSubject<IHeaderBreadcrumb[]>([]);
  }

  ngOnInit() {
    const { cfGuid, orgGuid, spaceGuid } = this.activeRouteCfOrgSpace;
    this.setupRouteObservable(this.getDefaultBackLink());
    // Will auto unsubscribe as we are using 'first'
    this.route$.pipe(first()).subscribe(route => {
      this.previousUrl = route.url;
      this.previousQueryParams = route.queryParams;
    });

    this.setupObservables();
    this.setupBreadcrumbs();
  }

  private getDefaultBackLink(): string {
    const { cfGuid, orgGuid, spaceGuid } = this.activeRouteCfOrgSpace;
    let backLink = `/cloud-foundry/${cfGuid}`;
    if (orgGuid) {
      backLink = `${backLink}/${orgGuid}`;
      if (spaceGuid) {
        backLink = `${backLink}/${spaceGuid}`;
      }
    }
    return backLink;
  }

  private setupRouteObservable(defaultBackLink: string) {
    this.route$ = this.store.select(getPreviousRoutingState).pipe(
      map(route => {
        return {
          url: route && route.state ? route.state.url : defaultBackLink,
          queryParams: route && route.state.queryParams ? route.state.queryParams : {}
        };
      })
    );
  }

  private setupObservables() {
    const { cfGuid, orgGuid, spaceGuid } = this.activeRouteCfOrgSpace;
    const org$ = orgGuid ? this.cfOrgService.org$ : Observable.of(null);
    const space$ = spaceGuid ? this.cfSpaceService.space$ : Observable.of(null);
    this.endpointOrgSpace$ = combineLatest(
      this.cfEndpointService.endpoint$,
      org$,
      space$
    );

    this.context$ = this.endpointOrgSpace$.pipe(
      map(([cf, org, space]) => {
        return {
          orgName: org ? org.entity.entity.name : null,
          spaceName: space ? space.entity.entity.name : null,
          apiEndpoint: getFullEndpointApiUrl(cf.entity),
          username: cf.entity.user ? cf.entity.user.name : ''
        };
      }),
      first()
    );
  }

  private setupBreadcrumbs() {
    const { cfGuid, orgGuid, spaceGuid } = this.activeRouteCfOrgSpace;
    this.breadcrumbs$ = this.endpointOrgSpace$.pipe(
      map(([cf, org, space]) => {
        const breadcrumbs = [{
          value: cf.entity.name,
          routerLink: `/cloud-foundry/${cf.entity.guid}`
        }];
        if (org) {
          breadcrumbs.push({
            value: org.entity.entity.name,
            routerLink: `/cloud-foundry/${cf.entity.guid}/organizations/${org.entity.metadata.guid}`
          });
          if (space) {
            breadcrumbs.push({
              value: space.entity.entity.name,
              routerLink: `/cloud-foundry/${cf.entity.guid}/organizations/${org.entity.metadata.guid}/spaces/${space.entity.metadata.guid}`
            });
          }
        }
        return [{ breadcrumbs }];
      }),
      first()
    );
  }

  back() {
    this.store.dispatch(new RouterNav({
      path: this.previousUrl,
      query: this.previousQueryParams
    }
    ));
  }
}

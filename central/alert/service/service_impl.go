package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/stackrox/rox/central/alert/datastore"
	notifierProcessor "github.com/stackrox/rox/central/notifier/processor"
	"github.com/stackrox/rox/central/processwhitelist"
	whitelistDatastore "github.com/stackrox/rox/central/processwhitelist/datastore"
	"github.com/stackrox/rox/central/role/resources"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/auth/permissions"
	"github.com/stackrox/rox/pkg/grpc/authz"
	"github.com/stackrox/rox/pkg/grpc/authz/perrpc"
	"github.com/stackrox/rox/pkg/grpc/authz/user"
	"github.com/stackrox/rox/pkg/protoconv"
	"github.com/stackrox/rox/pkg/search"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	badSnoozeErrorMsg = "'snooze_till' timestamp must be at a future time"

	maxListAlertsReturned = 1000
)

var (
	authorizer = perrpc.FromMap(map[authz.Authorizer][]string{
		user.With(permissions.View(resources.Alert)): {
			"/v1.AlertService/GetAlert",
			"/v1.AlertService/ListAlerts",
			"/v1.AlertService/CountAlerts",
			"/v1.AlertService/GetAlertsGroup",
			"/v1.AlertService/GetAlertsCounts",
			"/v1.AlertService/GetAlertTimeseries",
		},
		user.With(permissions.Modify(resources.Alert)): {
			"/v1.AlertService/ResolveAlert",
			"/v1.AlertService/SnoozeAlert",
			"/v1.AlertService/ResolveAlerts",
			"/v1.AlertService/DeleteAlerts",
		},
	})

	// groupByFunctions provides a map of functions that group slices of ListAlet objects by category or by cluser.
	groupByFunctions = map[v1.GetAlertsCountsRequest_RequestGroup]func(*storage.ListAlert) []string{
		v1.GetAlertsCountsRequest_UNSET: func(*storage.ListAlert) []string { return []string{""} },
		v1.GetAlertsCountsRequest_CATEGORY: func(a *storage.ListAlert) (output []string) {
			output = append(output, a.GetPolicy().GetCategories()...)
			return
		},
		v1.GetAlertsCountsRequest_CLUSTER: func(a *storage.ListAlert) []string { return []string{a.GetDeployment().GetClusterName()} },
	}
)

// serviceImpl is a thin facade over a domain layer that handles CRUD use cases on Alert objects from API clients.
type serviceImpl struct {
	dataStore  datastore.DataStore
	notifier   notifierProcessor.Processor
	whitelists whitelistDatastore.DataStore
}

// RegisterServiceServer registers this service with the given gRPC Server.
func (s *serviceImpl) RegisterServiceServer(grpcServer *grpc.Server) {
	v1.RegisterAlertServiceServer(grpcServer, s)
}

// RegisterServiceHandler registers this service with the given gRPC Gateway endpoint.
func (s *serviceImpl) RegisterServiceHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return v1.RegisterAlertServiceHandler(ctx, mux, conn)
}

// AuthFuncOverride specifies the auth criteria for this API.
func (s *serviceImpl) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	return ctx, authorizer.Authorized(ctx, fullMethodName)
}

// GetAlert returns the alert with given id.
func (s *serviceImpl) GetAlert(ctx context.Context, request *v1.ResourceByID) (*storage.Alert, error) {
	alert, exists, err := s.dataStore.GetAlert(ctx, request.GetId())
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !exists {
		return nil, status.Errorf(codes.NotFound, "alert with id '%s' does not exist", request.GetId())
	}

	return alert, nil
}

// ListAlerts returns ListAlerts according to the request.
func (s *serviceImpl) ListAlerts(ctx context.Context, request *v1.ListAlertsRequest) (*v1.ListAlertsResponse, error) {
	if request.GetPagination() == nil {
		request.Pagination = &v1.Pagination{
			Limit: maxListAlertsReturned,
		}
	}
	alerts, err := s.dataStore.ListAlerts(ctx, request)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &v1.ListAlertsResponse{Alerts: alerts}, nil
}

// CountAlerts counts the number of alerts that match the input query.
func (s *serviceImpl) CountAlerts(ctx context.Context, request *v1.RawQuery) (*v1.CountAlertsResponse, error) {
	// Fill in Query.
	parsedQuery, err := search.ParseRawQueryOrEmpty(request.GetQuery())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	alerts, err := s.dataStore.Search(ctx, parsedQuery)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &v1.CountAlertsResponse{Count: int32(len(alerts))}, nil
}

// GetAlertsGroup returns alerts according to the request, grouped by category and policy.
func (s *serviceImpl) GetAlertsGroup(ctx context.Context, request *v1.ListAlertsRequest) (*v1.GetAlertsGroupResponse, error) {
	alerts, err := s.dataStore.ListAlerts(ctx, request)
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := alertsGroupResponseFrom(alerts)
	return response, nil
}

// GetAlertsCounts returns alert counts by severity according to the request.
// Counts can be grouped by policy category or cluster.
func (s *serviceImpl) GetAlertsCounts(ctx context.Context, request *v1.GetAlertsCountsRequest) (*v1.GetAlertsCountsResponse, error) {
	alerts, err := s.dataStore.ListAlerts(ctx, request.GetRequest())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if response, ok := alertsCountsResponseFrom(alerts, request.GetGroupBy()); ok {
		return response, nil
	}

	return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("unknown group by: %v", request.GetGroupBy()))
}

// GetAlertTimeseries returns the timeseries format of the events based on the request parameters
func (s *serviceImpl) GetAlertTimeseries(ctx context.Context, req *v1.ListAlertsRequest) (*v1.GetAlertTimeseriesResponse, error) {
	alerts, err := s.dataStore.ListAlerts(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := alertTimeseriesResponseFrom(alerts)
	return response, nil
}

func (s *serviceImpl) ResolveAlert(ctx context.Context, req *v1.ResolveAlertRequest) (*v1.Empty, error) {
	alert, exists, err := s.dataStore.GetAlert(ctx, req.GetId())
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !exists {
		return nil, status.Errorf(codes.NotFound, "alert with id '%s' does not exist", req.GetId())
	}

	if req.GetWhitelist() {
		// This isn't great as it assumes a single whitelist key
		itemMap := make(map[string][]*storage.WhitelistItem)
		for _, process := range alert.GetProcessViolation().GetProcesses() {
			itemMap[process.GetContainerName()] = append(itemMap[process.GetContainerName()], &storage.WhitelistItem{
				Item: &storage.WhitelistItem_ProcessName{
					ProcessName: processwhitelist.WhitelistItemFromProcess(process),
				},
			})
		}
		for containerName, items := range itemMap {
			key := &storage.ProcessWhitelistKey{
				DeploymentId:  alert.GetDeployment().GetId(),
				ContainerName: containerName,
				ClusterId:     alert.GetDeployment().GetClusterId(),
				Namespace:     alert.GetDeployment().GetNamespace(),
			}
			if _, err := s.whitelists.UpdateProcessWhitelistElements(ctx, key, items, nil, false); err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
		}
	}

	if alert.LifecycleStage == storage.LifecycleStage_RUNTIME {
		err = s.changeAlertState(ctx, alert, storage.ViolationState_RESOLVED)
	}
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &v1.Empty{}, nil
}

func (s *serviceImpl) ResolveAlerts(ctx context.Context, req *v1.ResolveAlertsRequest) (*v1.Empty, error) {
	query, err := search.ParseRawQuery(req.GetQuery())
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	runtimeQuery := search.NewQueryBuilder().AddStrings(search.LifecycleStage, storage.LifecycleStage_RUNTIME.String()).ProtoQuery()
	cq := search.NewConjunctionQuery(query, runtimeQuery)
	alerts, err := s.dataStore.SearchRawAlerts(ctx, cq)
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	for _, alert := range alerts {
		err := s.changeAlertState(ctx, alert, storage.ViolationState_RESOLVED)
		if err != nil {
			log.Error(err)
		}
	}
	return &v1.Empty{}, nil
}

func (s *serviceImpl) changeAlertState(ctx context.Context, alert *storage.Alert, state storage.ViolationState) error {
	if state != storage.ViolationState_SNOOZED {
		alert.SnoozeTill = nil
	}
	alert.State = state
	err := s.dataStore.UpdateAlert(ctx, alert)
	if err != nil {
		log.Error(err)
		return status.Error(codes.Internal, err.Error())
	}
	s.notifier.ProcessAlert(alert)
	return nil
}

func (s *serviceImpl) SnoozeAlert(ctx context.Context, req *v1.SnoozeAlertRequest) (*v1.Empty, error) {
	if req.GetSnoozeTill() == nil {
		return nil, status.Error(codes.InvalidArgument, "'snooze_till' cannot be nil")
	}
	if protoconv.ConvertTimestampToTimeOrNow(req.GetSnoozeTill()).Before(time.Now()) {
		return nil, status.Error(codes.InvalidArgument, badSnoozeErrorMsg)
	}
	alert, exists, err := s.dataStore.GetAlert(ctx, req.GetId())
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !exists {
		return nil, status.Errorf(codes.NotFound, "alert with id '%s' does not exist", req.GetId())
	}
	alert.SnoozeTill = req.GetSnoozeTill()
	err = s.changeAlertState(ctx, alert, storage.ViolationState_SNOOZED)
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &v1.Empty{}, nil
}

// DeleteAlerts is a maintenance function that deletes alerts from the store
func (s *serviceImpl) DeleteAlerts(ctx context.Context, request *v1.DeleteAlertsRequest) (*v1.DeleteAlertsResponse, error) {
	if request.GetQuery() == nil {
		return nil, fmt.Errorf("a scoping query is required")
	}

	query, err := search.ParseRawQueryOrEmpty(request.GetQuery().GetQuery())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "error parsing query: %v", err)
	}
	query.Pagination = request.GetQuery().GetPagination()

	specified := false
	search.ApplyFnToAllBaseQueries(query, func(bq *v1.BaseQuery) {
		matchFieldQuery, ok := bq.GetQuery().(*v1.BaseQuery_MatchFieldQuery)
		if !ok {
			return
		}
		if matchFieldQuery.MatchFieldQuery.GetField() == search.ViolationState.String() {
			if matchFieldQuery.MatchFieldQuery.Value != storage.ViolationState_RESOLVED.String() {
				err = status.Errorf(codes.InvalidArgument, "invalid value for violation state: %q. Only resolved alerts can be deleted", matchFieldQuery.MatchFieldQuery.Value)
				return
			}
			specified = true
		}
	})
	if err != nil {
		return nil, err
	}
	if !specified {
		return nil, status.Errorf(codes.InvalidArgument, "please specify Violation State:%s in the query to confirm deletion", storage.ViolationState_RESOLVED.String())
	}

	results, err := s.dataStore.Search(ctx, query)
	if err != nil {
		return nil, err
	}

	response := &v1.DeleteAlertsResponse{
		NumDeleted: uint32(len(results)),
		DryRun:     !request.GetConfirm(),
	}

	if !request.GetConfirm() {
		return response, nil
	}

	idSlice := search.ResultsToIDs(results)
	if err := s.dataStore.DeleteAlerts(ctx, idSlice...); err != nil {
		return nil, err
	}
	return response, nil
}

// alertsGroupResponseFrom returns a slice of storage.ListAlert objects translated into a v1.GetAlertsGroupResponse object.
func alertsGroupResponseFrom(alerts []*storage.ListAlert) (output *v1.GetAlertsGroupResponse) {
	policiesMap := make(map[string]*storage.ListAlertPolicy)
	alertCountsByPolicy := make(map[string]int)

	for _, a := range alerts {
		if _, ok := policiesMap[a.GetPolicy().GetId()]; !ok {
			policiesMap[a.GetPolicy().GetId()] = a.GetPolicy()
		}
		alertCountsByPolicy[a.GetPolicy().GetId()]++
	}

	output = new(v1.GetAlertsGroupResponse)
	output.AlertsByPolicies = make([]*v1.GetAlertsGroupResponse_PolicyGroup, 0, len(policiesMap))

	for id, p := range policiesMap {
		output.AlertsByPolicies = append(output.AlertsByPolicies, &v1.GetAlertsGroupResponse_PolicyGroup{
			Policy:    p,
			NumAlerts: int64(alertCountsByPolicy[id]),
		})
	}

	sort.Slice(output.AlertsByPolicies, func(i, j int) bool {
		return output.AlertsByPolicies[i].GetPolicy().GetName() < output.AlertsByPolicies[j].GetPolicy().GetName()
	})

	return
}

// alertsCountsResponseFrom returns a slice of storage.ListAlert objects translated into a v1.GetAlertsCountsResponse
// object. True is returned if the translation was successful; otherwise false when the requested group is unknown.
func alertsCountsResponseFrom(alerts []*storage.ListAlert, groupBy v1.GetAlertsCountsRequest_RequestGroup) (*v1.GetAlertsCountsResponse, bool) {
	if groupByFunc, ok := groupByFunctions[groupBy]; ok {
		response := countAlerts(alerts, groupByFunc)
		return response, true
	}

	return nil, false
}

// alertTimeseriesResponseFrom returns a slice of storage.ListAlert objects translated into a v1.GetAlertTimeseriesResponse
// object.
func alertTimeseriesResponseFrom(alerts []*storage.ListAlert) *v1.GetAlertTimeseriesResponse {
	response := new(v1.GetAlertTimeseriesResponse)
	for cluster, severityMap := range getGroupToAlertEvents(alerts) {
		alertCluster := &v1.GetAlertTimeseriesResponse_ClusterAlerts{Cluster: cluster}
		for severity, alertEvents := range severityMap {
			// Sort the alert events so they are chronological
			sort.SliceStable(alertEvents, func(i, j int) bool { return alertEvents[i].GetTime() < alertEvents[j].GetTime() })
			alertCluster.Severities = append(alertCluster.Severities, &v1.GetAlertTimeseriesResponse_ClusterAlerts_AlertEvents{
				Severity: severity,
				Events:   alertEvents,
			})
		}
		sort.Slice(alertCluster.Severities, func(i, j int) bool { return alertCluster.Severities[i].Severity < alertCluster.Severities[j].Severity })
		response.Clusters = append(response.Clusters, alertCluster)
	}
	sort.SliceStable(response.Clusters, func(i, j int) bool { return response.Clusters[i].Cluster < response.Clusters[j].Cluster })
	return response
}

func countAlerts(alerts []*storage.ListAlert, groupByFunc func(*storage.ListAlert) []string) (output *v1.GetAlertsCountsResponse) {
	groups := getMapOfAlertCounts(alerts, groupByFunc)

	output = new(v1.GetAlertsCountsResponse)
	output.Groups = make([]*v1.GetAlertsCountsResponse_AlertGroup, 0, len(groups))

	for group, countsBySeverity := range groups {
		bySeverity := make([]*v1.GetAlertsCountsResponse_AlertGroup_AlertCounts, 0, len(countsBySeverity))

		for severity, count := range countsBySeverity {
			bySeverity = append(bySeverity, &v1.GetAlertsCountsResponse_AlertGroup_AlertCounts{
				Severity: severity,
				Count:    int64(count),
			})
		}

		sort.Slice(bySeverity, func(i, j int) bool {
			return bySeverity[i].Severity < bySeverity[j].Severity
		})

		output.Groups = append(output.Groups, &v1.GetAlertsCountsResponse_AlertGroup{
			Group:  group,
			Counts: bySeverity,
		})
	}

	sort.Slice(output.Groups, func(i, j int) bool {
		return output.Groups[i].Group < output.Groups[j].Group
	})

	return
}

func getMapOfAlertCounts(alerts []*storage.ListAlert, groupByFunc func(alert *storage.ListAlert) []string) (groups map[string]map[storage.Severity]int) {
	groups = make(map[string]map[storage.Severity]int)

	for _, a := range alerts {
		for _, g := range groupByFunc(a) {
			if groups[g] == nil {
				groups[g] = make(map[storage.Severity]int)
			}

			groups[g][a.GetPolicy().GetSeverity()]++
		}
	}

	return
}

func getGroupToAlertEvents(alerts []*storage.ListAlert) (clusters map[string]map[storage.Severity][]*v1.AlertEvent) {
	clusters = make(map[string]map[storage.Severity][]*v1.AlertEvent)
	for _, a := range alerts {
		alertCluster := a.GetDeployment().GetClusterName()
		if clusters[alertCluster] == nil {
			clusters[alertCluster] = make(map[storage.Severity][]*v1.AlertEvent)
		}
		eventList := clusters[alertCluster][a.GetPolicy().GetSeverity()]
		eventList = append(eventList, &v1.AlertEvent{Time: a.GetTime().GetSeconds() * 1000, Id: a.GetId(), Type: v1.Type_CREATED})
		if a.GetState() == storage.ViolationState_RESOLVED {
			eventList = append(eventList, &v1.AlertEvent{Time: a.GetTime().GetSeconds() * 1000, Id: a.GetId(), Type: v1.Type_REMOVED})
		}
		clusters[alertCluster][a.GetPolicy().GetSeverity()] = eventList
	}

	for _, v1 := range clusters {
		for k2, v2 := range v1 {
			sort.SliceStable(v2, func(i, j int) bool { return v2[i].GetTime() < v2[j].GetTime() })
			v1[k2] = v2
		}
	}
	return
}

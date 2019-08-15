package resolvers

import (
	"context"
	"time"

	"github.com/graph-gophers/graphql-go"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/central/metrics"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	pkgMetrics "github.com/stackrox/rox/pkg/metrics"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/set"
	"github.com/stackrox/rox/pkg/utils"
)

func init() {
	schema := getBuilder()
	utils.Must(
		schema.AddQuery("secret(id:ID!): Secret"),
		schema.AddQuery("secrets(query: String): [Secret!]!"),
		schema.AddExtraResolver("Secret", "deployments(query: String): [Deployment!]!"),
		schema.AddExtraResolver("Secret", "deploymentCount: Int!"),
	)
}

// Secret gets a single secret by ID
func (resolver *Resolver) Secret(ctx context.Context, arg struct{ graphql.ID }) (*secretResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Root, "Secret")
	if err := readSecrets(ctx); err != nil {
		return nil, err
	}

	secret := resolver.getSecret(ctx, string(arg.ID))
	if secret == nil {
		return resolver.wrapSecret(nil, false, errors.Errorf("error locating secret with id: %s", arg.ID))
	}
	return resolver.wrapSecret(secret, true, nil)
}

// Secrets gets a list of all secrets
func (resolver *Resolver) Secrets(ctx context.Context, args rawQuery) ([]*secretResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Root, "Secrets")
	if err := readSecrets(ctx); err != nil {
		return nil, err
	}
	q, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return nil, err
	}
	secrets, err := resolver.SecretsDataStore.SearchRawSecrets(ctx, q)
	if err != nil {
		return nil, err
	}

	for _, secret := range secrets {
		resolver.getDeploymentRelationships(ctx, secret)
	}
	return resolver.wrapSecrets(secrets, nil)
}

func (resolver *secretResolver) Deployments(ctx context.Context, args rawQuery) ([]*deploymentResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Secrets, "Deployments")
	if err := readDeployments(ctx); err != nil {
		return nil, err
	}

	q, err := resolver.getDeploymentQuery(args)
	if err != nil {
		return nil, err
	}
	return resolver.root.wrapDeployments(
		resolver.root.DeploymentDataStore.SearchRawDeployments(ctx, q))
}

func (resolver *secretResolver) DeploymentCount(ctx context.Context) (int32, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Secrets, "DeploymentCount")
	if err := readDeployments(ctx); err != nil {
		return 0, err
	}

	q, err := resolver.getDeploymentQuery(rawQuery{})
	if err != nil {
		return 0, err
	}
	results, err := resolver.root.DeploymentDataStore.Search(ctx, q)
	if err != nil {
		return 0, err
	}
	return int32(len(results)), err
}

func (resolver *secretResolver) getDeploymentQuery(args rawQuery) (*v1.Query, error) {
	deploymentFilterQuery, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return nil, err
	}

	secret := resolver.data
	deploymentIDs := set.NewStringSet()

	for _, dr := range secret.Relationship.GetDeploymentRelationships() {
		deploymentIDs.Add(dr.GetId())
	}
	deploymentIDQuery := search.NewQueryBuilder().AddDocIDs(deploymentIDs.AsSlice()...).ProtoQuery()

	return search.NewConjunctionQuery(deploymentIDQuery, deploymentFilterQuery), nil
}

func (resolver *Resolver) getSecret(ctx context.Context, id string) *storage.Secret {
	secret, ok, err := resolver.SecretsDataStore.GetSecret(ctx, id)
	if err != nil || !ok {
		return nil
	}

	resolver.getDeploymentRelationships(ctx, secret)
	return secret
}

func (resolver *Resolver) getDeploymentRelationships(ctx context.Context, secret *storage.Secret) {
	psr := search.NewQueryBuilder().
		AddExactMatches(search.ClusterID, secret.GetClusterId()).
		AddExactMatches(search.Namespace, secret.GetNamespace()).
		AddExactMatches(search.SecretName, secret.GetName()).
		ProtoQuery()

	deploymentResults, err := resolver.DeploymentDataStore.SearchListDeployments(ctx, psr)
	if err != nil {
		return
	}

	var deployments []*storage.SecretDeploymentRelationship
	for _, r := range deploymentResults {
		deployments = append(deployments, &storage.SecretDeploymentRelationship{
			Id:   r.Id,
			Name: r.Name,
		})
	}
	secret.Relationship = &storage.SecretRelationship{
		DeploymentRelationships: deployments,
	}
}

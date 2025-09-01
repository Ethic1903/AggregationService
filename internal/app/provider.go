package app

import (
	"AggregationService/internal/adapters/http/handlers"
	"AggregationService/internal/adapters/repository/postgres"
	"AggregationService/internal/converters"
	"AggregationService/internal/domain/ports/repository"
	"AggregationService/internal/domain/usecase/subscription_usecase"
	"AggregationService/internal/infrastructure/database/go_postgres"
	"AggregationService/internal/infrastructure/server"
	"AggregationService/internal/pkg/validation"
	"context"
)

type Provider struct {
	pgConfig     *go_postgres.IPGConfig
	pgClient     *go_postgres.PostgresClient
	serverConfig server.IServerConfig
	converter    *converters.SubscriptionConverter
	repo         repository.ISubscriptionRepository
	handler      *handlers.SubscriptionHandler
	usecase      subscription_usecase.ISubscriptionUseCase
	validator    *validation.Validator
}

func NewAppProvider() *Provider {
	return &Provider{}
}

func (p *Provider) PGConfig() *go_postgres.IPGConfig {
	if p.pgConfig == nil {
		config, err := go_postgres.NewPGConfig()
		if err != nil {
			panic(err)
		}
		p.pgConfig = config
	}
	return p.pgConfig
}

func (p *Provider) PGClient(ctx context.Context) *go_postgres.PostgresClient {
	if p.pgClient == nil {
		client, err := go_postgres.NewPGClient(ctx, *p.PGConfig())
		if err != nil {
			panic(err)
		}
		p.pgClient = client
	}
	return p.pgClient
}

func (p *Provider) SubscriptionRepo(ctx context.Context) repository.ISubscriptionRepository {
	if p.repo == nil {
		p.repo = postgres.NewSubscriptionsRepository(p.PGClient(ctx))
	}
	return p.repo
}

func (p *Provider) UseCase(ctx context.Context) subscription_usecase.ISubscriptionUseCase {
	if p.usecase == nil {
		p.usecase = subscription_usecase.New(
			p.SubscriptionRepo(ctx),
			p.Validator(),
			p.Converter(),
		)
	}
	return p.usecase
}

func (p *Provider) ServerConfig() server.IServerConfig {
	if p.serverConfig == nil {
		config, err := server.NewServerConfig()
		if err != nil {
			panic(err)
		}
		p.serverConfig = config
	}
	return p.serverConfig
}

func (p *Provider) Handler(ctx context.Context) *handlers.SubscriptionHandler {
	if p.handler == nil {
		p.handler = handlers.NewSubscriptionHandler(p.UseCase(ctx))
	}
	return p.handler
}

func (p *Provider) Converter() *converters.SubscriptionConverter {
	if p.converter == nil {
		p.converter = converters.New()
	}
	return p.converter
}

func (p *Provider) Validator() *validation.Validator {
	if p.validator == nil {
		p.validator, _ = validation.New()
	}
	return p.validator
}

package src

import (
	"github.com/MetaDandy/cuent-ai-core/config"
	"github.com/MetaDandy/cuent-ai-core/src/core/subscription"
	"github.com/MetaDandy/cuent-ai-core/src/core/user"
	"github.com/MetaDandy/cuent-ai-core/src/modules/asset"
	generatejob "github.com/MetaDandy/cuent-ai-core/src/modules/generate_job"
	"github.com/MetaDandy/cuent-ai-core/src/modules/project"
	"github.com/MetaDandy/cuent-ai-core/src/modules/script"
)

type Container struct {
	// User
	UserRepo *user.Repository
	UserSvc  *user.Service
	UserHdl  *user.Handler

	// Project
	ProjectRepo *project.Repository
	ProjectSvc  *project.Service
	ProjectHdl  *project.Handler

	// Script
	ScriptRepo *script.Repository
	ScriptSvc  *script.Service
	ScriptHdl  *script.Handler

	// Asset
	AssetRepo *asset.Repository
	AssetSvc  *asset.Service
	AssetHdl  *asset.Handler

	// Generated Job
	GeneratedJobRepo *generatejob.Repository

	// Subscription
	SubsRepo *subscription.Repository
	SubsSvc  *subscription.Service
	SubsHdl  *subscription.Handler
}

func SetupContainer() *Container {
	// User
	userRepo := user.NewRepository(config.DB)
	userSvc := user.NewService(userRepo)
	userHdl := user.NewHandler(userSvc)

	// Project
	projectRepo := project.NewRepository(config.DB)
	projectSvc := project.NewService(projectRepo, userRepo)
	projectHdl := project.NewHandler(projectSvc)

	// Generated Job
	generatedJobRepo := generatejob.NewRepository(config.DB)

	// Asset
	assetRepo := asset.NewRepository(config.DB)
	assetSvc := asset.NewService(assetRepo, generatedJobRepo, userRepo)
	assetHdl := asset.NewHandler(assetSvc)

	// Script
	scriptRepo := script.NewRepository(config.DB)
	scriptSvc := script.NewService(scriptRepo, projectRepo, assetRepo, userRepo)
	scriptHdl := script.NewHandler(scriptSvc)

	// Subscription
	subsRepo := subscription.NewRepository(config.DB)
	subsSvc := subscription.NewService(subsRepo)
	subsHdl := subscription.NewHandler(subsSvc)

	return &Container{
		// User
		UserRepo: userRepo,
		UserSvc:  userSvc,
		UserHdl:  userHdl,

		// Project
		ProjectRepo: projectRepo,
		ProjectSvc:  projectSvc,
		ProjectHdl:  projectHdl,

		// Asset
		AssetRepo: assetRepo,
		AssetSvc:  assetSvc,
		AssetHdl:  assetHdl,

		//Script
		ScriptRepo: scriptRepo,
		ScriptSvc:  scriptSvc,
		ScriptHdl:  scriptHdl,

		//Generated Job
		GeneratedJobRepo: generatedJobRepo,

		// Subscription
		SubsRepo: subsRepo,
		SubsSvc:  subsSvc,
		SubsHdl:  subsHdl,
	}
}

package src

import (
	"github.com/MetaDandy/cuent-ai-core/config"
	"github.com/MetaDandy/cuent-ai-core/src/core/subscription"
	"github.com/MetaDandy/cuent-ai-core/src/core/user"
	tts "github.com/MetaDandy/cuent-ai-core/src/modules/Tts"
	"github.com/MetaDandy/cuent-ai-core/src/modules/asset"
	"github.com/MetaDandy/cuent-ai-core/src/modules/cuentai"
	generatejob "github.com/MetaDandy/cuent-ai-core/src/modules/generate_job"
	"github.com/MetaDandy/cuent-ai-core/src/modules/project"
	"github.com/MetaDandy/cuent-ai-core/src/modules/script"
	"github.com/MetaDandy/cuent-ai-core/src/modules/supabase"
)

type Container struct {
	// TTS
	TtsSvc     *tts.Service
	TtsHandler *tts.Handler

	// CuentAI
	CuentSvc     *cuentai.Service
	CuentHandler *cuentai.Handler

	// Supabase
	SupaSvc     *supabase.Service
	SupaHandler *supabase.Handler

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
	// TTS
	ttsSvc := tts.NewService()
	ttsHandler := tts.NewHandler(ttsSvc)

	// Supabase
	supaSvc := supabase.NewService()
	supaHandler := supabase.NewHandler(supaSvc)

	// CuentAI
	cuentSvc := cuentai.NewService(ttsSvc, supaSvc)
	cuentHandler := cuentai.NewHandler(cuentSvc)

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
		// TTS
		TtsSvc:     ttsSvc,
		TtsHandler: ttsHandler,

		// CuentAI
		CuentSvc:     cuentSvc,
		CuentHandler: cuentHandler,

		// Supabase
		SupaSvc:     supaSvc,
		SupaHandler: supaHandler,

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

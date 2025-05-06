package src

import (
	"github.com/MetaDandy/cuent-ai-core/config"
	"github.com/MetaDandy/cuent-ai-core/src/core/user"
	tts "github.com/MetaDandy/cuent-ai-core/src/modules/Tts"
	"github.com/MetaDandy/cuent-ai-core/src/modules/cuentai"
	"github.com/MetaDandy/cuent-ai-core/src/modules/project"
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

	// Project
	ProjectRepo *project.Repository
	ProjectSvc  *project.Service
	ProjectHdl  *project.Handler
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

	// Project
	projectRepo := project.NewRepository(config.DB)
	projectSvc := project.NewService(projectRepo, userRepo)
	projectHdl := project.NewHandler(projectSvc)

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

		// Project
		ProjectRepo: projectRepo,
		ProjectSvc:  projectSvc,
		ProjectHdl:  projectHdl,
	}
}

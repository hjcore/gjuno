package x

import (
	"fmt"

	"github.com/hjcore/gjuno/database"
	"github.com/hjcore/gjuno/x/apis"
	"github.com/hjcore/gjuno/x/authz"

	"github.com/forbole/juno/v4/modules"
	"github.com/forbole/juno/v4/modules/registrar"
	"github.com/forbole/juno/v4/modules/telemetry"
	"github.com/forbole/juno/v4/node/remote"
)

type RegistrarOptions struct {
	APIsRegistrar    apis.Registrar
	APIsConfigurator apis.Configurator
}

func (o RegistrarOptions) GetAPIsRegistrar() apis.Registrar {
	if o.APIsRegistrar != nil {
		return o.APIsRegistrar
	}
	return apis.DefaultRegistrar
}

func (o RegistrarOptions) GetAPIsConfigurator() apis.Configurator {
	return o.APIsConfigurator
}

// --------------------------------------------------------------------------------------------------------------------

// ModulesRegistrar represents the modules.Registrar that allows to register all custom GJuno modules
type ModulesRegistrar struct {
	options RegistrarOptions
}

// NewModulesRegistrar allows to build a new ModulesRegistrar instance
func NewModulesRegistrar() *ModulesRegistrar {
	return &ModulesRegistrar{}
}

// WithOptions sets the given option inside this registrar
func (r *ModulesRegistrar) WithOptions(options RegistrarOptions) *ModulesRegistrar {
	r.options = options
	return r
}

// BuildModules implements modules.Registrar
func (r *ModulesRegistrar) BuildModules(ctx registrar.Context) modules.Modules {
	cdc := ctx.EncodingConfig.Marshaler
	gjunoDb := database.Cast(ctx.Database)

	remoteCfg, ok := ctx.JunoConfig.Node.Details.(*remote.Details)
	if !ok {
		panic(fmt.Errorf("cannot run GJuno on local node"))
	}

	grpcConnection := remote.MustCreateGrpcConnection(remoteCfg.GRPC)

	// Juno modules
	telemetryModule := telemetry.NewModule(ctx.JunoConfig)

	// GJuno modules
	apisModule := apis.NewModule(apis.NewContext(ctx, grpcConnection))
	if apisModule != nil {
		apisModule = apisModule.WithRegistrar(r.options.GetAPIsRegistrar())
		apisModule = apisModule.WithConfigurator(r.options.GetAPIsConfigurator())
	}

	authzModule := authz.NewModule(ctx.Proxy, cdc, gjunoDb)

	return []modules.Module{
		apisModule,
		authzModule,
		telemetryModule,
	}
}

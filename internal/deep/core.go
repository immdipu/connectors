package deep

import (
	"github.com/amp-labs/connectors/common/interpreter"
	"github.com/amp-labs/connectors/common/paramsbuilder"
	"github.com/amp-labs/connectors/internal/deep/dprequests"
	"github.com/amp-labs/connectors/internal/deep/dpvars"
	"github.com/amp-labs/connectors/internal/deep/requirements"
	"github.com/amp-labs/connectors/providers"
	"go.uber.org/dig"
)

// Connector is the main method to build deep connectors.
// It accepts constructor method and connector components.
// This procedure builds dependency tree and provides them by constructor injections.
// Your connector constructor will therefore receive any requirements.ConnectorComponent passed by arguments.
//
// Ex:
//
// (connector constructor)
//
//	 	constructor := func(
//			clients *deep.Clients,
//			closer *deep.EmptyCloser,
//			reader *deep.Reader,
//			writer *deep.Writer,
//			staticMetadata *deep.StaticMetadata,
//			remover *deep.Remover,
//			// Any custom requirements.ConnectorComponent can be injected here
//		) *Connector {
//			return &Connector{
//				Clients:        *clients,
//				EmptyCloser:    *closer,
//				Reader:         *reader,
//				Writer:         *writer,
//				StaticMetadata: *staticMetadata,
//				Remover:        *remover,
//			}
//		}
//
// (putting everything together)
//
//	 deep.Connector[Connector, parameters](constructor, providers.Salesforce, opts,
//			errorHandler,
//			objectURLResolver,
//			firstPage,
//			nextPage,
//			...
//			// Custom requirements.ConnectorComponent are passed here
//			// to be available for all ConnectorComponent constructors
//		)
//
// List of connector components, will be injected resembling "Template Design Pattern" using "uber dig" package.
// Ex: To build Connector, it needs Reader, which in turn needs Clients, and the last depends on Parameters.
// One by one will be built ultimately resolving Connector constructor.
//
// You can replace default connector components by implementing requirements.ConnectorComponent interface.
// It must return the same requirements.Dependency ID to act as override.
func Connector[C any, P paramsbuilder.ParamAssurance](
	connectorConstructor any,
	provider providers.Provider,
	options []func(params *P),
	components ...requirements.ConnectorComponent,
) (*C, error) {
	return ExtendedConnector[C, P, *dpvars.EmptyMetadataVariables](
		connectorConstructor, provider, &dpvars.EmptyMetadataVariables{}, options, components...,
	)
}

// ExtendedConnector is the same connector builder as Connector method.
// The main difference is introduction of concrete dpvars.MetadataVariables type.
// Connector may have additional metadata stored in the struct.
// You can define arbitrary holder of this data using this builder: ExtendedConnector.
func ExtendedConnector[C any, P paramsbuilder.ParamAssurance, D dpvars.MetadataVariables]( //nolint:funlen
	connectorConstructor any,
	provider providers.Provider,
	metadataVariables D,
	options []func(params *P),
	components ...requirements.ConnectorComponent,
) (*C, error) {
	// This is a default list of dependencies available for a "connectorConstructor" to pick up.
	dependencies := requirements.NewDependencies([]requirements.Dependency{
		{
			// Connector must have Provider name.
			ID: requirements.Provider,
			Constructor: func() providers.Provider {
				return provider
			},
		},
		{
			// Connector is configured using options.
			ID: requirements.Options,
			Constructor: func() []func(params *P) {
				return options
			},
		},
		// Metadata Variables hold connector specific data fields.
		// They are inferred from parameters
		metadataVariables.Satisfies(),

		// Options are realized into parameters.
		// Some parameters can be catalog or metadata variables.
		dpvars.Parameters[P]{}.Satisfies(),
		dpvars.ConnectorData[P, D]{}.Satisfies(),
		dpvars.CatalogVariables[P, D]{}.Satisfies(),

		// Every connector makes requests. Clients holds authenticated HTTP clients
		// capable of doing JSON, XML calls. It needs parameters, catalog vars for proper setup.
		// ErrorHandler would parse error response depending on media type.
		// HeaderSupplements is used to attach headers when performing said requests.
		{
			ID:          requirements.Clients,
			Constructor: dprequests.NewClients[P, D],
		},
		interpreter.ErrorHandler{}.Satisfies(),

		// Most connectors do no-op on close.
		// *EmptyCloser is available as constructor argument.
		EmptyCloser{}.Satisfies(),

		{
			// This is the main constructor which will get all dependencies resolved.
			// It is possible that not all dependencies are needed, this list is exhaustive,
			// which describes all the building blocks that Deep connector may have.
			ID:          requirements.Connector,
			Constructor: connectorConstructor,
		},
	})

	for _, component := range components {
		dependencies.Add(component.Satisfies())
	}

	container := dig.New()
	if err := dependencies.Apply(container); err != nil {
		return nil, err
	}

	return resolveDependencies[C](container)
}

// Tries to invoke connector constructor. The dig.Container will have dependencies that will be injected.
// Produces connector of a specified T type.
func resolveDependencies[T any](container *dig.Container) (*T, error) {
	var result *T

	err := container.Invoke(func(builder *T) {
		result = builder
	})

	return result, err
}

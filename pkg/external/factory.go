package external

import (
	"fmt"

	"github.com/sanjasimiccc/movieManagementSystem/pkg/types"
)

type ProviderFactory struct {
	providers map[string]types.MovieProvider
}

func NewProviderFactory() *ProviderFactory {
	return &ProviderFactory{
		providers: map[string]types.MovieProvider{
			"OMDb": NewOMDbProvider("5257cb6d"),
			"TMDb": NewTMDbProvider("081ceea64560bb90d333e9fb727f1927"),
		},
	}
}

func (f *ProviderFactory) Get(source string) (types.MovieProvider, error) {
	provider, ok := f.providers[source]
	if !ok {
		return nil, fmt.Errorf("unknown provider: %s", source)
	}
	return provider, nil
}

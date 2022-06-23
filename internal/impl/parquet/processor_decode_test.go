package parquet

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/benthosdev/benthos/v4/public/service"
)

func TestParquetDecodeProcessorPokemonFans(t *testing.T) {
	reader, err := newParquetDecodeProcessor(nil)
	require.NoError(t, err)

	pBytes, err := os.ReadFile("./resources/pokemon_fans.parquet")
	require.NoError(t, err)

	readerResBatch, err := reader.Process(context.Background(), service.NewMessage(pBytes))
	require.NoError(t, err)

	var readerResStrs []string
	for _, m := range readerResBatch {
		mBytes, err := m.AsBytes()
		require.NoError(t, err)
		readerResStrs = append(readerResStrs, string(mBytes))
	}

	assert.Equal(t, []string{
		`{"nameIn":"fooer first","age":21,"Id":1,"weight":60.1,"favPokemon":null}`,
		`{"nameIn":"fooer second","age":22,"Id":2,"weight":60.2,"favPokemon":null}`,
		`{"nameIn":"fooer third","age":23,"Id":3,"weight":60.3,"favPokemon":[{"pokeName":"bulbasaur","coolness":99}]}`,
		`{"nameIn":"fooer fourth","age":24,"Id":4,"weight":60.4,"favPokemon":null}`,
		`{"nameIn":"fooer fifth","age":25,"Id":5,"weight":60.5,"favPokemon":null}`,
		`{"nameIn":"fooer sixth","age":26,"Id":6,"weight":60.6,"favPokemon":null}`,
	}, readerResStrs)
}

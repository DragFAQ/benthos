package parquet

import (
	"errors"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/segmentio/parquet-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParquetDecodePokemonFans(t *testing.T) {
	pFile, err := os.Open("./resources/pokemon_fans.parquet")
	require.NoError(t, err)

	var resRows []map[string]any

	// File schema is as follows:
	//
	// message root {
	// 	required binary name (STRING);
	// 	required int32 age;
	// 	required int64 id;
	// 	required float weight;
	// 	optional group favPokemon {
	// 		repeated group list {
	// 			required group element {
	// 				required binary name (STRING);
	// 				required float coolness;
	// 			}
	// 		}
	// 	}
	// }

	pRdr := parquet.NewGenericReader[any](pFile)
	schema := pRdr.Schema()
	rowBuf := make([]parquet.Row, 10)
	for {
		n, err := pRdr.ReadRows(rowBuf)
		if err != nil && !errors.Is(err, io.EOF) {
			t.Fatal(err)
		}
		if n == 0 {
			break
		}

		for i := 0; i < n; i++ {
			row := rowBuf[i]

			mappedData := map[string]any{}
			for j, f := range schema.Fields() {
				value := row[j]
				if value.IsNull() {
					mappedData[f.Name()] = nil
					continue
				}

				if !f.Leaf() {
					for _, innerField := range f.Fields() {
						fmt.Printf("Field %v is a child of %v but I don't know how to expand it\n", innerField.Name(), f.Name())
					}
					continue
				}
				mappedData[f.Name()] = value.String()
			}

			resRows = append(resRows, mappedData)
		}
	}

	assert.Equal(t, []map[string]any{
		{"name": "fooer first", "age": "21", "id": "1", "weight": "60.1", "favPokemon": nil},
		{"name": "fooer second", "age": "22", "id": "2", "weight": "60.2", "favPokemon": nil},
		{"name": "fooer third", "age": "23", "id": "3", "weight": "60.3", "favPokemon": []map[string]any{{"pokeName": "bulbasaur", "coolness": "99"}}},
		{"name": "fooer fourth", "age": "24", "id": "4", "weight": "60.4", "favPokemon": nil},
		{"name": "fooer fifth", "age": "25", "id": "5", "weight": "60.5", "favPokemon": nil},
		{"name": "fooer sixth", "age": "26", "id": "6", "weight": "60.6", "favPokemon": nil},
	}, resRows)
}

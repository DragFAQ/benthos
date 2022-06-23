package parquet

import (
	"bytes"
	"context"
	"errors"
	"io"

	"github.com/segmentio/parquet-go"

	"github.com/benthosdev/benthos/v4/public/service"
)

func parquetDecodeProcessorConfig() *service.ConfigSpec {
	return service.NewConfigSpec().
		// Stable(). TODO
		Categories("Parsing").
		Summary("Decodes [Parquet files](https://parquet.apache.org/docs/) into a batch of structured messages.").
		Description(``).
		Version("4.3.0")
}

func init() {
	err := service.RegisterProcessor(
		"parquet_decode", parquetDecodeProcessorConfig(),
		func(conf *service.ParsedConfig, mgr *service.Resources) (service.Processor, error) {
			return newParquetDecodeProcessorFromConfig(conf, mgr.Logger())
		})

	if err != nil {
		panic(err)
	}
}

//------------------------------------------------------------------------------

func newParquetDecodeProcessorFromConfig(conf *service.ParsedConfig, logger *service.Logger) (*parquetDecodeProcessor, error) {
	return newParquetDecodeProcessor(logger)
}

type parquetDecodeProcessor struct {
	logger *service.Logger
}

func newParquetDecodeProcessor(logger *service.Logger) (*parquetDecodeProcessor, error) {
	s := &parquetDecodeProcessor{
		logger: logger,
	}
	return s, nil
}

func (s *parquetDecodeProcessor) Process(ctx context.Context, msg *service.Message) (service.MessageBatch, error) {
	mBytes, err := msg.AsBytes()
	if err != nil {
		return nil, err
	}

	rdr := bytes.NewReader(mBytes)
	pRdr := parquet.NewReader(rdr)

	rowBuf := make([]parquet.Row, 10)
	var resBatch service.MessageBatch

	schema := pRdr.Schema()
	for {
		n, err := pRdr.ReadRows(rowBuf)
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, err
		}
		if n == 0 {
			break
		}

		for i := 0; i < n; i++ {
			row := rowBuf[i]

			// TODO: Do this properly
			data := map[string]interface{}{}
			for j, f := range schema.Fields() {
				value := row[j]
				data[f.Name()] = value.String()
			}

			newMsg := msg.Copy()
			newMsg.SetStructured(data)
			resBatch = append(resBatch, newMsg)
		}
	}

	return resBatch, nil
}

func (s *parquetDecodeProcessor) Close(ctx context.Context) error {
	return nil
}

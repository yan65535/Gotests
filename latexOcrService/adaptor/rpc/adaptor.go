package rpc

import (
	"latexOcrService/config"
	ocrRpc "latexOcrService/latexOcrProto"
)

type ILatex interface {
}
type Adaptor struct {
	conf        config.Conf
	latexClient ocrRpc.LatexServiceClient
}

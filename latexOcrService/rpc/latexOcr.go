package rpc

import (
	"latexOcrService/adaptor"
	"latexOcrService/config"
	ocrRpc "latexOcrService/latexOcrProto"
)

type OcrGrpcService struct {
	ocrRpc.UnimplementedLatexServiceServer
}

func NewOcrGrpcService(adaptors adaptor.Adaptors, conf config.Conf) *OcrGrpcService {
	return &OcrGrpcService{
		//sender: sender.NewService(adaptors, conf),
	}
}

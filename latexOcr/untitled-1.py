import grpc
from concurrent import futures
import latex_pb2
import latex_pb2_grpc
import io
from PIL import Image
import pytesseract
from pix2tex.cli import LatexOCR
class LatexService(latex_pb2_grpc.LatexServiceServicer):
    def RecognizeLatex(self, request, context):
        try:
            # 将接收到的字节数据转换为图片
            image = Image.open(io.BytesIO(request.image_data))
            # 确保图像是 RGB 模式
            image = image.convert('RGB')
            # 初始化 LatexOCR 模型
            model = LatexOCR()
            # 使用模型进行识别
            result = model(image)
            return latex_pb2.LatexResponse(result=result)
        except Exception as e:
            # 处理可能出现的异常
            print(f"识别出错: {e}")
            return latex_pb2.LatexResponse(result="")

# ... existing code ...
def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    latex_pb2_grpc.add_LatexServiceServicer_to_server(LatexService(), server)
    # 修改监听端口为 50052
    server.add_insecure_port('[::]:50052')
    server.start()
    print("Server started, listening on port 50052")
    server.wait_for_termination()
# ... existing code ...

if __name__ == '__main__':
    serve()

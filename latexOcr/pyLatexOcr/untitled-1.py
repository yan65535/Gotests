import grpc
from concurrent import futures
import latex_pb2
import latex_pb2_grpc
import io
from PIL import Image
from pix2tex.cli import LatexOCR
from kazoo.client import KazooClient
#protobuf==5.29.0
class LatexService(latex_pb2_grpc.LatexServiceServicer):
    def __init__(self):
        # 初始化时加载模型（避免每次请求都重新加载）
        self.model = LatexOCR()
        # ... existing code ...
    def RecognizeLatex(self, request, context):
        try:
            # 将接收到的字节数据转换为图片
            image = Image.open(io.BytesIO(request.image_data))
            # 确保图像是 RGB 模式
            image = image.convert('RGB')
            # 初始化 LatexOCR 模型

            # 使用模型进行识别
            result = self.model(image)
            return latex_pb2.LatexResponse(result=result)
        except Exception as e:
            # 处理可能出现的异常
            print(f"识别出错: {e}")
            return latex_pb2.LatexResponse(result="")

def serve():
    zk = KazooClient(hosts='127.0.0.1:2181')  # 请根据实际情况修改 ZooKeeper 地址
    zk.start()
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    latex_pb2_grpc.add_LatexServiceServicer_to_server(LatexService(), server)
    server.add_insecure_port('[::]:50053')
    server.start()
    print("Server started, listening on port 50053")

    # 注册服务到 ZooKeeper
    service_path = '/services/latex_ocr'
    port = '[::]:50053'
    if not zk.exists(service_path):
        zk.create(service_path, makepath=True)
    zk.create(f'{service_path}/node', value=port.encode(), ephemeral=True, sequence=True)

    try:
        server.wait_for_termination()
    except KeyboardInterrupt:
        server.stop(0)
        zk.stop()


    server.wait_for_termination()

if __name__ == '__main__':
    serve()
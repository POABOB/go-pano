from concurrent import futures
import grpc
from protos.predict_pb2 import Response
import protos.predict_pb2_grpc as predict_pb2_grpc

# AI辨識LIB
import argparse
import json
import sys
import os


from pathlib import Path

from utils.core import main


class RPCService(predict_pb2_grpc.PredictServicer):
    def predict(self, request, context):
        if request.Dir == "":
            print(request)
            return Response(isSuccessful = False, msg = "辨識失敗，請傳遞正確的資料夾名稱", predict = "")
        Dir = request.Dir

        try:
            data_dir = Path("/app/go/static/img/" + Dir)
            _D = list(data_dir.glob('*.jpg'))
            tooth_anomaly_dict = main(_D)

            text = """Tx:Check Panoramic radiography films for initial examination in this year
            C.F.:
            1.缺牙: %(missing)s
            2.殘根: %(R.R)s  埋伏齒: %(embedded)s
            3.固定補綴物: %(filling)s
            4.齲齒: %(caries)s
            5.曾經根管治療: %(endo)s"""

            predict = []
            anomaly_list = ['R.R', 'caries', 'crown', 'endo', 'post', 'filling', 'Imp', 'embedded', 'impacted', 'missing']
            for filename, teeth in tooth_anomaly_dict.items():
                filename = f'{filename}.jpg'
                temp = {}
                temp["filename"] = filename

                teeth_anomalies_dict = {anomaly: [] for anomaly in anomaly_list}
                for tooth_number, anomalies in teeth.items():
                    for anomaly in anomalies:
                        teeth_anomalies_dict[anomaly].append(tooth_number)

                text_anomalies_list = ['missing', 'R.R', 'embedded', 'filling', 'caries', 'endo']
                text_dict = {
                    anomaly: ' '.join(map(str, teeth_anomalies_dict[anomaly])) if teeth_anomalies_dict[
                        anomaly] else 'no finding' for
                    anomaly in text_anomalies_list}

                temp['text'] = text % text_dict
                temp['data'] = teeth_anomalies_dict
                predict.append(temp)

                print(predict)
                return Response(isSuccessful = True, msg = "成功", predict = json.dumps(predict))
        except Exception as err:
            print(f"Unexpected {err=}, {type(err)=}")
            return Response(isSuccessful = False, msg = str(err), predict = "")

# 將剛剛實作完的service架起來
def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=5))
    predict_pb2_grpc.add_PredictServicer_to_server(RPCService(), server)
    print('server start at 0.0.0.0:5001')
    server.add_insecure_port("0.0.0.0:5001")
    server.start()
    server.wait_for_termination()
 
 
if __name__ == "__main__":
    serve() # run gRPC server
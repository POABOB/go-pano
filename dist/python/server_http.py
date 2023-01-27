import argparse
import json
import sys
import os

from flask import Flask, request
from utils.core import main

app = Flask(__name__)

@app.route("/")
def Start():
    # 獲取queryString
    query = request.args

    if "Dir" not in query.keys():
        return json.dumps({
            "isSuccessful": False,
            "msg": "辨識失敗，請傳遞正確的資料夾名稱",
            "predict": ""
        }, ensure_ascii=False)

    Dir = query.get("Dir")
    
    # try:
    tooth_anomaly_dict = main("/app/go/static/img/" + Dir)

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


    # 假資料
    # predict = {"msg": "Good"}
    result = {
        "isSuccessful": True,
        "msg": "成功",
        # 把json轉成STR
        "predict": json.dumps(predict)
    }


    return json.dumps(result, ensure_ascii=False)
    # except Exception as e:
    #     result = {
    #         "isSuccessful": False,
    #         "msg": str(e),
    #         "predict": ""
    #     }
    #     print(result)
    #     return json.dumps(result, ensure_ascii=False)



if __name__ == '__main__':
    app.run(debug=True,host='0.0.0.0')


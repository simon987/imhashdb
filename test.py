import requests
from base64 import b64encode
import json

with open("/home/simon/Downloads/a.jpg", "rb") as f:
    data = f.read()

r = requests.post("http://localhost:8080/api/hash", data=json.dumps({
    "data": b64encode(data).decode()
}))

# print(r.content)

for i in range (0, 49):
    r2 = requests.post("http://localhost:8080/api/query", data=json.dumps({
        "hash": r.json()["ahash:12"],
        "type": "ahash:12",
        "distance": 30,
        "limit": 500 + i,
        "offset": 0
    }))
    print(r2.content.decode())



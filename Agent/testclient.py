import asyncio
import json
from websockets import connect

async def hello(uri):
    print("Connected")
    async with connect(uri) as websocket:
        # set ping timeout to 5 seconds
        # websocket.ping_interval = 1

        # create json object
        json_object = {
            "uuid": "myuuid",
            "name": "agent007",
            "os": "Windows 10.1",
            "arch": "amd64",
        }

        # turn json into string
        json_string = json.dumps(json_object)

        # send json
        await websocket.send(json_string)
        
        while (1):
            print("Waiting for work")
            # get work
            data = await websocket.recv()

            print("work received:", data)

            # send "file data"
            await websocket.send("test file contents")

            print("file sent")

            # get return value
            data = await websocket.recv()

            print("response:", data)

asyncio.run(hello("ws://localhost:8080/api/v1/agent/ws"))
import asyncio
from websockets import connect

async def hello(uri):
    print("Connected")
    async with connect(uri) as websocket:
        await websocket.send("\{uuid\:myuuid,uuid\:myuuid,name\:agent007,os\:windows,arch\:amd64\}")
        data = await websocket.recv()

        print(data)

        await websocket.send("test file contents")

        data = await websocket.recv()

        print(data)

asyncio.run(hello("ws://localhost/agent:4201"))
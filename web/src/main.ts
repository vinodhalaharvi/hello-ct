import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { HelloService } from "./gen/hello/hello_connect";

const transport = createConnectTransport({
  baseUrl: "http://localhost:8080",
});

const client = createClient(HelloService, transport);

document.getElementById("btn")?.addEventListener("click", async () => {
  const res = await client.sayHello({ name: "World" });
  document.getElementById("result")!.innerText = res.greeting?.message || "No response";
});

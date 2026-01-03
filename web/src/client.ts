import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { HelloService } from "./gen/hello/hello_connect.js";

const transport = createConnectTransport({
  baseUrl: "http://localhost:8080",
});

export const helloClient = createClient(HelloService, transport);

// Test it
async function main() {
  const response = await helloClient.sayHello({ name: "World" });
  console.log(response.greeting?.message);
}

main();

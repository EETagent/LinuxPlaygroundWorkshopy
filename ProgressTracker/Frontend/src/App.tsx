import { Component, createResource, For, Show } from "solid-js";
import Terminal from "./Components/Terminal";
import { createSignal, onCleanup } from "solid-js";
import Tabs from "./Components/HostnameTabs";

interface ChallengeInterface {
  hostname: string;
  command: string;
}
interface DashBoardInterface {
  challenge: string;
  data: ChallengeInterface;
}

const findUniqueHostnames = async (data: DashBoardInterface[]) => [
  "Všechno",
  ...new Set(data.map((e) => e.data.hostname)),
];

const App: Component = () => {
  const websocket: WebSocket = new WebSocket("ws://75.119.149.184:8080/dashboard");

  const [data, setData] = createSignal<DashBoardInterface[]>();

  const [hostnames] = createResource(data, findUniqueHostnames);
  const [selectedHostname, setSelectedHostname] =
    createSignal<string>("Všechno");

  websocket.addEventListener("message", (event) => {
    const json: DashBoardInterface[] = JSON.parse(event.data);
    /*json.forEach(j => {
      j.data.command = j.data.command.replace("Dokončeno: 0.000000 %", "Dokončeno: " + Math.floor(Math.random() * (100  + 1)) + ".123400 %")
    })*/
    setData(json);
  });
  websocket.addEventListener("open", () => {
    console.log("Opened");
  });

  onCleanup(() => {
    websocket.close();
  });

  return (
    <div
      className="flex flex-col min-w-screen min-h-screen"
      style="background: repeating-linear-gradient(
      -45deg,
      black,
      black 30px,
      #444 31px,
      #444 32px
    );"
    >
      <Show when={data() != undefined}>
        <header>
          <Show when={hostnames.loading !== true}>
            <Tabs
              tabs={hostnames()}
              active={selectedHostname()}
              selectedCallback={(hostname) => setSelectedHostname(hostname)}
            />
          </Show>
        </header>
        <main>
          <For
            each={
              selectedHostname() !== "Všechno"
                ? data().filter((e) => e.data.hostname === selectedHostname())
                : data()
            }
          >
            {(item: DashBoardInterface) => (
              <div className="my-10">
                <div className="mb-3">
                  <div className="font-sans mb-1 text-5xl font-bold text-white">
                    {item.data.hostname} / {item.challenge}
                  </div>
                  <div className="w-full h-6 bg-gray-200 rounded-full dark:bg-gray-700">
                    <div
                      className="h-6 bg-gray-600 rounded-full dark:bg-gray-300"
                      style={`width: ${parseInt(item.data.command.split("Dokončeno: ")[1])}%`}
                    ></div>
                  </div>
                </div>
                <Terminal
                  username={item.data.hostname}
                  command={item.challenge}
                  output={item.data.command}
                />
              </div>
            )}
          </For>
        </main>
      </Show>
    </div>
  );
};

export default App;

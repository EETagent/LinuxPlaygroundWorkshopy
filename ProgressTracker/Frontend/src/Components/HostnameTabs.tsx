import { Component, createSignal, For, onMount } from "solid-js";

const Tabs: Component<{
  tabs: string[];
  selectedCallback?: Function;
  active?: string;
}> = ({ tabs, selectedCallback, active }) => {
  const [activeTab, setActiveTab] = createSignal<number>();
  onMount(() => {
    setActiveTab(tabs.indexOf(active));
  });
  return (
    <ul className="h-full w-full table m-0 p-0 list-none">
      <For each={tabs}>
        {(tab: string, id) => (
          <li
            className={`font-sans table-cell text-center bg-white inline-block ml-0 p-3 border-0 border-b-2 transition-all cursor-pointer ${
              activeTab() === id()
                ? "border-b-blue-700 text-gray-400"
                : "text-gray-200"
            } `}
            onclick={() => {
              setActiveTab(id);
              selectedCallback(tab);
            }}
          >
            {tab}
          </li>
        )}
      </For>
    </ul>
  );
};

export default Tabs;

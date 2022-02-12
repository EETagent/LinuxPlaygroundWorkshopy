import { Component } from "solid-js";

const Terminal: Component<{
  username: string;
  command: string;
  output: string;
}> = ({ username, command, output }) => {
  return (
    <div className="w-full">
      <div
        className="font-sans px-5 pt-4 shadow-lg text-gray-100 text-sm  subpixel-antialiased
              bg-[#252221] pb-6 rounded-lg leading-normal overflow-hidden"
      >
        <div className="top mb-2 flex">
          <div className="h-3 w-3 bg-[#dd665a] rounded-full" />
          <div className="ml-2 h-3 w-3 bg-[#efbe5c] rounded-full" />
          <div className="ml-2 h-3 w-3 bg-[#7ac656] rounded-full" />
        </div>
        <div className="mt-4 flex flex-col leading-relaxed">
          <div className="flex flex-row">
            <span className="text-[#6384c7]">${username}</span>
            <span>@hackdays</span>
            <span className="ml-1 text-[#78b56c]">{command}</span>
          </div>
          <p className="mb-3" style="white-space: pre-line">
            {output}
          </p>
        </div>
      </div>
    </div>
  );
};

export default Terminal;

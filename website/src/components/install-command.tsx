"use client";

import { useState } from "react";
import { CopyButton } from "@/components/copy-button";

export function InstallCommand({ variant = "hero" }: { variant?: "hero" | "section" }) {
  const [os, setOs] = useState<"unix" | "windows">("unix");

  const unixCommand = "curl -fsSL https://aaxion.codershubinc.com/install.sh | bash";
  const windowsCommand = "irm https://aaxion.codershubinc.com/install.ps1 | iex";

  const isHero = variant === "hero";

  return (
    <div className={`flex flex-col gap-3 w-full max-w-2xl mx-auto ${isHero ? "items-center" : "items-start"}`}>
      {/* Toggle */}
      <div className="flex bg-white/5 border border-white/10 rounded-full p-1 backdrop-blur-md w-fit relative z-50 pointer-events-auto">
        <button
          type="button"
          onClick={() => { console.log('Clicked Unix'); setOs("unix"); }}
          className={`px-4 py-1 text-xs font-semibold rounded-full transition-all cursor-pointer ${
            os === "unix"
              ? "bg-white/10 text-white shadow-sm"
              : "text-gray-500 hover:text-gray-300"
          }`}
        >
          Linux / macOS
        </button>
        <button
          type="button"
          onClick={() => { console.log('Clicked Windows'); setOs("windows"); }}
          className={`px-4 py-1 text-xs font-semibold rounded-full transition-all cursor-pointer ${
            os === "windows"
              ? "bg-white/10 text-white shadow-sm"
              : "text-gray-500 hover:text-gray-300"
          }`}
        >
          Windows
        </button>
      </div>

      {/* Command Box */}
      <div
        className={`bg-black/60 border border-white/10 flex items-center justify-between gap-4 w-full transition-all backdrop-blur-md ${
          os === "windows" ? "justify-center text-gray-500 italic" : "group hover:border-white/20"
        } ${isHero ? "p-2 pl-5 pr-2 rounded-full" : "p-5 rounded-2xl"}`}
      >
        {os === "windows" ? (
          <span className="text-sm py-1.5 md:py-0">Not available for now</span>
        ) : (
          <>
            <code className={`${isHero ? "text-sm" : "text-sm md:text-base"} text-blue-300 font-mono overflow-x-auto whitespace-nowrap`}>
              <span className="text-gray-600 select-none mr-3">$</span>
              {unixCommand}
            </code>
            <CopyButton text={unixCommand} />
          </>
        )}
      </div>
    </div>
  );
}

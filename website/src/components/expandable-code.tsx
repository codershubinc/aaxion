"use client";

import { useState } from "react";
import { CopyButton } from "@/components/copy-button";

export function ExpandableCode({ label, code, variant }: { label: string; code: string; variant: "request" | "response" }) {
  const [open, setOpen] = useState(false);
  const isReq = variant === "request";

  return (
    <div className="w-full">
      <button
        type="button"
        onClick={() => setOpen((prev) => !prev)}
        className={`
          inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg text-[11px] font-medium border transition-all cursor-pointer select-none
          ${isReq
            ? "bg-violet-500/10 border-violet-500/20 text-violet-400 hover:bg-violet-500/20"
            : "bg-cyan-500/10 border-cyan-500/20 text-cyan-400 hover:bg-cyan-500/20"
          }
        `}
      >
        <svg
          suppressHydrationWarning
          className={`w-3 h-3 transition-transform duration-200 ${open ? "rotate-90" : ""}`}
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          strokeWidth={2.5}
        >
          <path strokeLinecap="round" strokeLinejoin="round" d="M8.25 4.5l7.5 7.5-7.5 7.5" />
        </svg>
        {label}
      </button>

      {open && (
        <div
          className={`mt-2 rounded-xl border overflow-hidden
            ${isReq ? "bg-violet-500/5 border-violet-500/10" : "bg-cyan-500/5 border-cyan-500/10"}
          `}
        >
          <div className="flex items-center justify-between px-3 py-1.5 border-b border-white/5">
            <p className={`text-[10px] font-semibold uppercase tracking-widest ${isReq ? "text-violet-500" : "text-cyan-500"}`}>
              {isReq ? "Request Body" : "Response"}
            </p>
            <CopyButton text={code} />
          </div>
          <pre className={`px-4 py-3 text-xs overflow-x-auto font-mono leading-relaxed ${isReq ? "text-violet-200" : "text-cyan-200"}`}>
            <code>{code}</code>
          </pre>
        </div>
      )}
    </div>
  );
}

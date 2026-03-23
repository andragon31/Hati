// @ts-nocheck
/* eslint-disable */
declare var Bun: any;
declare var process: any;

import type { Plugin } from "@opencode-ai/plugin"

const HATI_BIN = process.env.HATI_BIN ?? "hati"

export const Hati: Plugin = async (ctx) => {
  return {
    "experimental.chat.system.transform": async (_input, output) => {
      const instr = `\n\n## Hati Protocol\nYou have access to Hati, the AI Execution Layer.`
      if (output.system.length > 0) {
        output.system[output.system.length - 1] += instr
      } else {
        output.system.push(instr)
      }
    },
  }
}

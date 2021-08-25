import { gqlTip } from "../queries";
import { TTip } from "../types";
import create from "zustand";

type TUseTipStore = {
  tip?: TTip;
  refreshTip: () => void;
}

export const useTip = create<TUseTipStore>((set, get) => ({
  tip: undefined,
  refreshTip: () => gqlTip().then((tip) => set({ tip })),
}));

// Polling blockchain tip. Save interval id to window so hot-reload can clear
const startTipPolling = () => {
  useTip.getState().refreshTip();
  clearInterval(window.__TIP_INTERVAL__)
  window.__TIP_INTERVAL__ = setInterval(() => {
    useTip.getState().refreshTip();
  }, 1000 * 2);
};
startTipPolling();

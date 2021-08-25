import { gqlGetUtxos, gqlMintAsset, gqlWalletCreate, gqlWalletFund } from "../queries";
import create from "zustand";
import { TUtxo } from "components/types";

const LOCAL_STORAGE_KEY_WALLET = "cardanoToolkit:walletAddress";

type TUseWalletStore = {
  isWalletConnected: boolean;
  walletAddress?: string;
  walletBalanceADA?: number;
  walletBalanceAssets: { [assetId: string]: number };
  walletUtxos: TUtxo[];
  connectWallet: (addressString?: string, exists?: boolean) => void;
  disconnectWallet: () => void;
  fundWallet: (quantity?: string) => void;
  mintAsset: (name: string, quantity: string) => void;
  refreshUtxos: () => void;
}

export const useWallet = create<TUseWalletStore>((set, get) => ({
  // DATA
  isWalletConnected: !!localStorage.getItem(LOCAL_STORAGE_KEY_WALLET),
  walletAddress: localStorage.getItem(LOCAL_STORAGE_KEY_WALLET) ?? undefined,
  walletBalanceADA: undefined,
  walletBalanceAssets: {},
  walletUtxos: [],
  // CONNECTING
  connectWallet: (addressString, exists) => {
    // --- If users says it exists already, just set values, don't create
    if (exists && addressString) {
      set({ isWalletConnected: true, walletAddress: addressString });
      localStorage.setItem(LOCAL_STORAGE_KEY_WALLET, addressString);
      return;
    }
    // --- If doesnt exist, try creating
    gqlWalletCreate(addressString)
      .then((walletAddress) => {
        set({ isWalletConnected: true, walletAddress });
        localStorage.setItem(LOCAL_STORAGE_KEY_WALLET, walletAddress);
      })
      .catch(() => {
        set({ isWalletConnected: false, walletAddress: undefined })
        localStorage.removeItem(LOCAL_STORAGE_KEY_WALLET);
      });
  },
  disconnectWallet: () => {
    set({ isWalletConnected: false, walletAddress: undefined })
    localStorage.removeItem(LOCAL_STORAGE_KEY_WALLET);
  },
  // ACTIONS
  fundWallet: (quantity) => {
    if (!get().walletAddress) {
      return console.error("No connected wallet");
    }
    gqlWalletFund(get().walletAddress!, quantity);
  },
  mintAsset: (name, quantity) => {
    if (!get().walletAddress) {
      return console.error("No connected wallet");
    }
    gqlMintAsset(name, quantity, get().walletAddress!);
  },
  refreshUtxos: () => {
    if (!get().walletAddress) {
      set({ walletBalanceADA: undefined, walletBalanceAssets: {} });
      return console.error("No connected wallet");
    }
    gqlGetUtxos({ address: get().walletAddress!, excludeScripts: true }).then((walletUtxos) => {
      set({
        walletBalanceADA: walletUtxos
          .filter((utxo) => utxo.tokens.length === 0)
          .reduce((total, utxo) => total + Number(utxo.value!), 0),
        walletBalanceAssets: walletUtxos
          .map((utxo) => utxo.tokens)
          .flat()
          .reduce((acc, token) => {
            const prevSum = acc[token.asset.assetId];
            return {
              ...acc,
              [token.asset.assetId]: prevSum
                ? Number(prevSum) + Number(token.quantity)
                : Number(token.quantity)
            };
          }, {} as TUseWalletStore["walletBalanceAssets"]),
        walletUtxos,
      })
    });
  }
}));

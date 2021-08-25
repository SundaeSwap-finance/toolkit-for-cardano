export type TAsset = {
  assetId: string;
  assetName: string;
  policyId: string;
  description?: string;
  logo?: string;
  name?: string;
  ticker?: string;
  url?: string;
};

export type TAssetAmount = {
  asset: TAsset;
  amount: number;
};

export type TToken = {
  asset: TAsset;
  quantity: string;
};

export type TUtxo = {
  address: string;
  datumHash?: string;
  index: number;
  tokens: TToken[];
  value: string;
};

export type TTip = {
  block: number;
  epoch: number;
  era: string;
  hash: string;
  slot: number;
};

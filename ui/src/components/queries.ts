import toast from "react-hot-toast";
import { TUtxo } from "./types";

const API_URL = "";

const gql = (
  query: string,
  variables: Record<string, any>,
  key: string,
  debug?: boolean
): Promise<any> => {
  // Debug option for printing to console
  if (debug) {
    console.log(`gql: ${key}`);
    console.log(query);
    console.log(JSON.stringify(variables, null, 2));
  }
  // Make request
  return fetch(API_URL, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ query, variables }),
  })
    .then((res) => res.json())
    .then((result) => {
      // Notification
      // Debug
      if (debug) {
        console.log(JSON.stringify(result, null, 2));
      }
      // Return
      if (result.errors) throw result.errors;
      return result.data[key] ?? result.data;
    })
    .catch((err) => {
      console.error(err);
      throw err?.[0] ?? err;
    });
};


/**
 * Creates a new address and optionally funds it with the specified amount of ADA.
 * @param initialFunds
 */
export const gqlWalletCreate = (initialFunds?: string): Promise<string> => {
  return gql(
    `mutation(
      $initialFunds: String!
    ) {
      walletCreate(initialFunds: $initialFunds)
    }`,
    {
      initialFunds: initialFunds || String(423_654_000),
    },
    "walletCreate",
  );
};

/**
 * Fund the specified address with ADA. Deposits 1,000 ADA by default (1e3 * 1e6)
 * @param address
 * @param quantity
 * @returns
 */
 export const gqlWalletFund = (
  address: string,
  quantity: string = "1000000000"
): Promise<void> => {
  toast(`Funding wallet ${Number(quantity) / 1_000_000} ADA.`);
  return gql(
    `
    mutation(
      $address: String!
      $quantity: String!
    ) {
      walletFund(address: $address, quantity: $quantity) {
        ok
      }
    }
    `,
    {
      address,
      quantity
    },
    "walletFund",
  );
};

/**
 * Mint a new token
 * @param assetName
 * @param quantity
 * @param address
 * @returns
 */
 export const gqlMintAsset = (
  assetName: string,
  quantity: string,
  walletAddress: string
): Promise<void> => {
  toast(`Minting ${quantity} ${assetName}.`);
  return gql(
    `mutation(
      $assetName: String!
      $quantity: String!
      $walletAddress: String!
    ) {
      mint(assetName: $assetName, quantity: $quantity, wallet: $walletAddress) {
        ok
      }
    }`,
    {
      assetName,
      quantity,
      walletAddress,
    },
    "mint"
  );
};

export const gqlGetUtxos = ({
  address,
  excludeScripts,
  excludeTokens,
}: {
  address?: string;
  excludeScripts?: boolean;
  excludeTokens?: boolean;
} = {}): Promise<TUtxo[]> => {
  return gql(
    `
    query (
      $address: String
      $excludeScripts: Boolean
      $excludeTokens: Boolean
    ) {
      utxos(
        address: $address,
        excludeScripts: $excludeScripts,
        excludeTokens: $excludeTokens,
      ) {
        address,
        datumHash,
        index,
        tokens {
          asset {
            assetId
            assetName
            #description
            #logo
            #name
            policyId
            #ticker
            #url
          },
          quantity
        },
        value,
      }
    }
    `,
    {
      address,
      excludeScripts,
      excludeTokens,
    },
    "utxos",
  );
};

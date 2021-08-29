import { TTip, TUtxo } from "./types";

const API_URL = "/graphql";

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
 * @param name
 */
export const gqlWalletCreate = (name?: string): Promise<string> => {
  return gql(
    `mutation(
      $initialFunds: String!
      $name: String!
    ) {
      walletCreate(
        initialFunds: $initialFunds,
        name: $name,
      )
    }`,
    {
      initialFunds: String(423_654_000),
      name,
    },
    "walletCreate",
  );
};

/**
 * Gets all wallets
 * @param name
 */
export const gqlWallets = (): Promise<string[]> => {
  return gql(
    `query {
      wallets(query: "")
    }`,
    {},
    "wallets",
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

/**
 * Get all utxos for a wallet
 * @param address
 * @param excludeScripts
 * @param excludeTokens
 * @returns
 */
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

/**
 * Get blockchain tip data
 */
export const gqlTip = (): Promise<TTip> => {
  return gql(
    `
    query {
      tip {
        block
        epoch
        era
        hash
        slot
      }
    }
    `,
    {},
    "tip"
  );
};

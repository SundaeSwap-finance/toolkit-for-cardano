import React, { useEffect, useState } from "react";
import styled from "styled-components";
import { useWallet } from "./useWallet";
import useInterval from "../hooks/useInterval";
import { Button } from "../styled/Button";
import Input from "../styled/Input";
import Select from "../styled/Select";
import { gqlWallets } from "../queries";

export const Wallet = () => {
  // Get wallet/actions context
  const {
    isWalletConnected,
    walletAddress,
    connectWallet,
    disconnectWallet,
    refreshUtxos,
    walletBalanceADA,
    walletBalanceAssets,
    walletUtxos,
  } = useWallet();
  // Continually refresh utxos (on mount/connection also fetch)
  useEffect(() => {
    if (isWalletConnected) refreshUtxos();
  }, [isWalletConnected]);
  useInterval(() => { refreshUtxos(); }, 4000);
  // Section
  const [walletSection, setWalletSection] = useState<"balance" | "utxos">("balance");
  const [walletTokenFilter, setWalletTokenFilter] = useState("");
  // Save wallet string being typed in, to be created
  const [newWalletAddress, setNewWalletAddress] = useState("");
  // Get wallet options for user (refresh on wallet creations/disconnections)
  const [walletOptions, setWalletOptions] = useState<string[]>([]);
  useEffect(() => { gqlWallets().then(setWalletOptions); }, [isWalletConnected]);

  // Render
  // --- Not connected
  if (!isWalletConnected) {
    return (
      <StyledWallet>
        <div className="wallet__disconnected">
          <div className="create">
            <small>Create a new wallet!</small>
            <Input value={walletAddress} placeholder="-- wallet address --" onChange={(e) => setNewWalletAddress(e.target.value)} />
            <Button size="xs" onClick={() => connectWallet(newWalletAddress, false)}>
              Create New Wallet
            </Button>
          </div>
          <div className="connect">
            <small>Connect to an existing wallet!</small>
            <Select onChange={(e) => connectWallet(e.target.value, true)}>
              <option>---</option>
              {walletOptions.sort().map((wallet) => <option key={wallet} value={wallet}>{wallet}</option>)}
            </Select>
          </div>
        </div>
      </StyledWallet>
    );
  }

  // --- Connected
  return (
    <StyledWallet>
      <div className="wallet__header">
        <header>
          <div>
            <small>Wallet</small>
            {walletAddress}
          </div>
          <Button size="xxs" onClick={disconnectWallet}>Disconnect</Button>
        </header>
        <div>
          <span className={walletSection === "balance" ? "active" : ""} onClick={() => setWalletSection("balance")}>Balances</span>
          <span className={walletSection === "utxos" ? "active" : ""} onClick={() => setWalletSection("utxos")}>UTXOs</span>
        </div>
      </div>
      <div className="wallet__body">
        {walletSection === "balance" && (
          <>
            <div className="wallet__ada">
              <span>ADA Balance:</span>
              <span>₳ {(Number(walletBalanceADA ?? 0) / 1_000_000).toLocaleString()}</span>
            </div>
            <div className="wallet__assets">
              {Object.entries(walletBalanceAssets ?? {}).map(([assetId, amount]) => (
                <div key={assetId} className="wallet__asset">
                  <div>
                    <p>{assetId.split('.')[1]}</p>
                    <small>{assetId.split('.')[0]}</small>
                  </div>
                  <div>
                    <p>{amount}</p>
                  </div>
                </div>
              ))}
            </div>
          </>
        )}
        {walletSection === "utxos" && (
          <>
            <div className="wallet__utxos__filter">
              <Input
                value={walletTokenFilter}
                placeholder="Filter by token name..."
                onChange={(e) => setWalletTokenFilter(e.target.value)}
              />
            </div>
            <div className="wallet__utxos">
              {walletUtxos
                .filter((utxo) => !!walletTokenFilter ? utxo.tokens.some((t) => new RegExp(walletTokenFilter, 'i').test(t.asset.assetName)) : true)
                .sort((a, b) => Number(b.value) - Number(a.value))
                .map((utxo) => (
                  <div key={`${utxo.address}-${utxo.index}`} className="wallet__utxo">
                    <p className="wallet__utxo__values">
                      <span>{utxo.value} ₳</span>
                      {utxo.tokens.map((token) => (
                        <span key={`${token.asset.assetName}${token.quantity}`}>{token.quantity} {token.asset.assetName}</span>
                      ))}
                    </p>
                    <p className="wallet__utxo__address">
                      <span>{utxo.address}</span><span>#{utxo.index}</span>
                    </p>
                  </div>
                ))}
            </div>
          </>
        )}
      </div>
    </StyledWallet>
  );
};

export const StyledWallet = styled.div`
  transition: box-shadow 0.25s ease;
  height: 520px;
  width: 420px;
  max-width: 100%;
  overflow: hidden;
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-liq-card);
  background: var(--main-bg);
  color: var(--text);

  .wallet__disconnected {
    height: 100%;
    width: 100%;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    text-align: center;
    padding: 1em 2em;
    input, select, option {
      text-align: center;
    }
    .create, .connect {
      border-radius: var(--radius-sm);
      border: 4px solid var(--swapbox-bg);
      height: 100%;
      width: 100%;
      margin: 1em 0;
      display: flex;
      align-items: center;
      justify-content: center;
      flex-direction: column;
      small {
        margin-bottom: 6px;
      }
      input, select {
        margin-bottom: 8px;
        width: 80%;
      }
      select {
        margin-bottom: 24px;
      }
      button {
        margin-bottom: 16px;
      }
    }
  }

  .wallet__header {
    & > header {
      display: flex;
      justify-content: space-between;
      padding: 20px 18px;
      & > div {
        display: flex;
        flex-direction: column;
        small {
          font-size: 10px;
          opacity: 0.4;
          margin-top: -4px;
          margin-bottom: 2px;
        }
      }
    }
    & > div {
      display: flex;
      justify-content: flex-start;
      padding: 0 20px;
      span {
        margin-right: 12px;
        padding-bottom: 6px;
        border-bottom: 2px solid var(--text);
        color: var(--text);
        font-size: 14px;
        cursor: pointer;
        &.active {
          color: var(--text-primary);
          border-bottom: 2px solid var(--text-primary);
        }
      }
    }
  }
  .wallet__body {
    height: 100%;
    overflow-y: scroll;
    -ms-overflow-style: none;
    scrollbar-width: none;
    &::-webkit-scrollbar {
      display: none;
    }
  }

  .wallet__ada {
    padding: 18px 18px;
    border-bottom: 1px solid var(--main-bg);
    background: var(--swapbox-bg);
    display: flex;
    justify-content: space-between;
  }

  .wallet__assets {
    height: auto;
    background: var(--swapbox-bg);
  }
  .wallet__asset {
    display: flex;
    justify-content: space-between;
    padding: 18px 18px;
    border-bottom: 1px solid var(--main-bg);
    & > div:first-of-type {
      max-width: 50%;
      padding-top: 4px;
      p {
        margin: 0;
        color: var(--text-primary);
      }
    }
    & > div:last-of-type {
      display: flex;
      align-items: center;
      p {
        margin: 0;
        font-family: Inter, sans-serif;
      }
    }
    small {
      font-size: 8px;
      opacity: 0.3;
    }
  }

  .wallet__utxos {
    height: auto;
    background: var(--swapbox-bg);
  }
  .wallet__utxos__filter {
    width: 100%;
    input {
      width: 100%;
      border-radius: 0;
      font-size: 12px;
      padding: 14px 20px 10px;
    }
  }
  .wallet__utxo {
    display: flex;
    flex-direction: column;
    padding: 6px 18px;
    border-bottom: 1px solid var(--main-bg);
    font-size: 10px;
    letter-spacing: -0.2px;
    &__address {
      margin: 0;
      display: flex;
      justify-content: space-between;
      width: 100%;
      color: var(--text-muted);
    }
    &__values {
      margin: 2px 0 4px;
      span:not(:first-of-type) {
        color: var(--text-primary);
        &::before {
          content: "/";
          margin: 0 4px;
          color: var(--text-silent);
        }
      }
    }
  }
`;

import React, { useEffect, useState } from "react";
import styled from "styled-components";
import { useWallet } from "./useWallet";
import useInterval from "../hooks/useInterval";
import { Button } from "../styled/Button";
import Input from "../styled/Input";

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
  } = useWallet();
  // Save wallet string being selected/entered
  const [walletAddressSelection, setWalletAddressSelection] = useState("");
  // Continually refresh utxos (on mount/connection also fetch)
  useEffect(() => {
    if (isWalletConnected) refreshUtxos();
  }, [isWalletConnected]);
  useInterval(() => { refreshUtxos(); }, 4000);

  // Render
  // --- Not connected
  if (!isWalletConnected) {
    return (
      <StyledWallet>
        <div className="wallet__disconnected">
          <div className="connect">
            <small>Connect to an existing wallet!</small>
            <Input value={walletAddress} placeholder="---" onChange={(e) => setWalletAddressSelection(e.target.value)} />
            <Button size="xs" disabled={!walletAddressSelection} onClick={() => connectWallet(walletAddressSelection, !!walletAddressSelection)}>
              Connect to Wallet
            </Button>
          </div>
          <div className="create">
            <small>Create a new wallet!</small>
            <Button size="xs" onClick={() => connectWallet()}>
              Create New Wallet
            </Button>
          </div>
        </div>
      </StyledWallet>
    );
  }
  // --- Connected
  return (
    <StyledWallet>
      <div className="wallet__address">
        <div>
          <small>Wallet</small>
          {walletAddress}
        </div>
        <Button size="xxs" onClick={disconnectWallet}>Disconnect</Button>
      </div>
      <div className="wallet__ada">
        <span>Balance:</span>
        <span>₳ {(Number(walletBalanceADA) / 1_000_000).toLocaleString()}</span>
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
    </StyledWallet>
  );
};

export const StyledWallet = styled.div`
  transition: box-shadow 0.25s ease;
  height: 480px;
  width: 420px;
  max-width: 100%;
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-liq-card);
  background: var(--main-bg);
  color: var(--text);
  overflow-y: scroll;

  .wallet__disconnected {
    height: 100%;
    width: 100%;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    text-align: center;
    padding: 1em 2em;
    input {
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
        margin-bottom: 12px;
      }
      input {
        margin-bottom: 8px;
        width: 80%;
      }
    }
  }

  .wallet__address {
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
  .wallet__ada {
    padding: 18px 18px;
    border-top: 1px solid var(--main-bg);
    border-bottom: 1px solid var(--main-bg);
    background: var(--swapbox-bg);
    display: flex;
    justify-content: space-between;
  }
  .wallet__assets {
    height: 100%;
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
`;

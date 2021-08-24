import { Button } from "../styled/Button";
import React, { useState } from "react";
import styled from "styled-components";
import { useWallet } from "./useWallet";
import Input from "../styled/Input";

export const WalletActions = () => {
  const { fundWallet, isWalletConnected, mintAsset } = useWallet();
  // ADA Faucet Helpers
  const [fundQuantity, setFundQuantity] = useState<string>("");
  const fundHandler = () => {
    fundWallet(fundQuantity);
    setFundQuantity("");
  }
  // Asset Minting Helpers
  const [mintAssetQuantity, setMintAssetQuantity] = useState<string>("");
  const [mintAssetName, setMintAssetName] = useState<string>("");
  const mintAssetHandler = () => {
    mintAsset(mintAssetName, mintAssetQuantity);
    setMintAssetQuantity("");
    setMintAssetName("");
  }

  // Render
  // --- Not connected
  if (!isWalletConnected) {
    return (
      <StyledWalletActions>
        <h2>Connect a wallet to begin!</h2>
      </StyledWalletActions>
    )
  }
  // --- Connected
  return (
    <StyledWalletActions>
      <StyledWalletAction>
        <div className="wallet__action__header">
          <small>ADA Faucet</small>
        </div>
        <div className="wallet__action__body">
          <Input type="text" placeholder="0.0" value={fundQuantity} onChange={(e) => setFundQuantity(e.target.value)} />
          <Button size="xs" onClick={fundHandler}>Fund</Button>
        </div>
      </StyledWalletAction>
      <StyledWalletAction>
        <div className="wallet__action__header">
          <small>Mint Asset Name/Amount</small>
        </div>
        <div className="wallet__action__body">
          <Input type="text" placeholder="---" value={mintAssetName} onChange={(e) => setMintAssetName(e.target.value)} />
          <Input type="text" placeholder="0.0" value={mintAssetQuantity} onChange={(e) => setMintAssetQuantity(e.target.value)} />
          <Button size="xs" onClick={mintAssetHandler}>Mint</Button>
        </div>
      </StyledWalletAction>
    </StyledWalletActions>
  );
};


export const StyledWalletAction = styled.div`
  background: var(--main-bg);
  box-shadow: var(--shadow-connectWallet);
  border-radius: var(--radius-lg);
  min-height: 40px;
  width: 320px;
  min-width: 320px;
  max-width: 100%;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  text-align: center;
  ${Button}, ${Input}, small {
    width: 100%;
    text-align: center;
  }
  ${Input}, small {
    margin-bottom: 8px;
  }
  ${Input} {
    background: var(--swapbox-bg);
    &::placeholder {
      color: var(--text);
    }
  }
  .wallet__action__header {
    background: var(--main-bg);
    padding: 0.5em 1em;
  }
  .wallet__action__body {
    background: var(--main-bg);
    padding: 0 1em 1em 1em;
  }
  small {
    color: var(--text);
  }
`;

export const StyledWalletActions = styled.div`
  display: flex;
  flex-direction: column;
  & > h2 {
    color: var(--text);
  }
  ${StyledWalletAction} {
    margin: 12px 0;
  }
`;
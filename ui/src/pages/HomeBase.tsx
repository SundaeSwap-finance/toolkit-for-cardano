import React from "react";
import styled from "styled-components";
import { Wallet } from "../components/wallet/Wallet";
import { WalletActions } from "../components/wallet/WalletActions";

export const HomeBase = () => {
  return (
    <StyledHomeBase>
      <section>
        <Wallet />
      </section>
      <section>
        <WalletActions />
      </section>
    </StyledHomeBase>
  );
}

const StyledHomeBase = styled.main`
  display: flex;
  height: 100%;
  width: 100%;
  max-width: 960px;
  margin: 0 auto;
  padding: 1em;
  section {
    height: 100%;
    width: 50%;
    display: flex;
    justify-content: center;
    align-items: center;
    flex-direction: column;
    & > h2 {
      margin-top: -54px;
      color: var(--text);
      opacity: 0.6;
    }
  }
  @media(max-width: 45em) {
    flex-direction: column;
    section {
      width: 100%;
      margin: 4em 0;
    }
  }
`;

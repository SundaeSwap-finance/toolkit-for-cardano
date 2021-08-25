import { A } from "../styled/A";
import React from "react";
import styled from "styled-components";

export const Footer = () => {
  return (
    <StyledFooter>
      <p>
        <A href="https://sundaeswap.finance" target="_blank">Scoops by SundaeSwap</A>
      </p>
      <p>
        <A href="https://github.com/SundaeSwap-finance/cardano-toolkit" target="_blank">Github Repo</A>
      </p>
    </StyledFooter>
  );
};

export const StyledFooter = styled.footer`
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 4px 20px;
  p {
    font-size: 12px;
  }
  @media(max-width: 45em) {
    position: initial;
  }
`;

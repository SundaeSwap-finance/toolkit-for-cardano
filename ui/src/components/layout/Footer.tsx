import React from "react";
import styled from "styled-components";

export const Footer = () => {
  // TODO: Include epoch/blocks?
  return (
    <StyledFooter>
      <p>Cardano Toolkit Made by Your Friends at SundaeSwap</p>
    </StyledFooter>
  );
};

export const StyledFooter = styled.footer`
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  text-align: center;
  color: white;
  opacity: 0.5;
`;

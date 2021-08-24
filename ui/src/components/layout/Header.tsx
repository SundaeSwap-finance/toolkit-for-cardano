import React from "react";
import styled from "styled-components";
import iceCreamCup from "../../assets/ice-cream-cup.png";

export const Header = () => {
  return (
    <StyledHeader>
      <div>{/* <h1>Cardano Toolkit</h1> */}</div>
    </StyledHeader>
  )
}

const StyledHeader = styled.header`
  position: absolute;
  top: 20px;
  left: 0;
  right: 0;
  text-align: center;

  background-image: url(${iceCreamCup});
  background-position: center;
  background-size: contain;
  background-repeat: no-repeat;

  div {
    height: 40px;
    h1 {
      font-size: 18px;
      font-weight: 700;
      font-family: "MilliardExtraBold";
      color: var(--text-secondary);
      opacity: 0.5;
      margin: 0;
      padding: 1em;
    }
  }
  @media(max-width: 45em) {
    position: initial;
    margin-top: 24px;
  }
`;

import React from "react";
import styled from "styled-components";
import sundaeLeft from "../../assets/sundae_left.png";
import sundaeRight from "../../assets/sundae_right.png";

export const Header = () => {
  return (
    <StyledHeader>
      <div className="sundae-asset--left" />
      <div>{/* <h1>Cardano Toolkit</h1> */}</div>
      <div className="sundae-asset--right" />
    </StyledHeader>
  )
}

const StyledHeader = styled.header`
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  text-align: center;
  display: flex;
  justify-content: space-between;
  height: 72px;

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

  .sundae-asset--left, .sundae-asset--right {
    background-position: center;
    background-size: contain;
    background-repeat: no-repeat;
    height: 54px;
    width: 54px;
    margin: 20px;
  }
  .sundae-asset--left {
    background-image: url(${sundaeLeft});
    transform: rotate(-30deg);
  }
  .sundae-asset--right {
    background-image: url(${sundaeRight});
    transform: rotate(30deg);
  }
`;

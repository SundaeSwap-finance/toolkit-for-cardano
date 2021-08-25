import { Tip } from "../tip/Tip";
import React from "react";
import styled from "styled-components";
import sundaeLeft from "../../assets/sundae_left.png";
import sundaeRight from "../../assets/sundae_right.png";

export const Header = () => {
  return (
    <StyledHeader>
      <div className="sundae-asset--left" />
      <div className="header__tip">
        <Tip />
      </div>
      <div className="sundae-asset--right" />
    </StyledHeader>
  )
}

const StyledHeader = styled.header`
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  display: flex;
  justify-content: space-between;
  height: 72px;

  .header__tip {
    height: 100%;
    display: flex;
    align-items: center;
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

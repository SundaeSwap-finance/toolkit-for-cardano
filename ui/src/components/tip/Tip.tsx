import React, { useEffect } from "react";
import styled from "styled-components";
import { useTip } from "./useTip";

export const Tip = () => {
  const { tip } = useTip();

  if (!tip) return null;
  return (
    <StyledTip>
      <div>
        <small>Epoch:</small>
        <span className="counter">{tip?.epoch}</span>
      </div>
      <div>
        <small>Block:</small>
        <span className="counter">{Number(tip?.block).toLocaleString()}</span>
      </div>
    </StyledTip>
  );
}

const StyledTip = styled.div`
  display: flex;
  overflow: hidden;
  border-radius: var(--radius-md);
  background: var(--gradient-card-header);
  box-shadow: var(--shadow-overview);
  & > div {
    display: flex;
    align-items: center;
    small {
      font-size: 10px;
      font-weight: 700;
      color: var(--text-muted);
    }
    .counter {
      margin-left: 8px;
      font-size: 12px;
      font-weight: 700;
      color: var(--secondary);
    }
    &:first-of-type {
      padding: 12px 6px 12px 24px;
    }
    &:last-of-type {
      padding: 12px 24px 12px 6px;
    }
  }
`;

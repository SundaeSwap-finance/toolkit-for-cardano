import styled, { css } from "styled-components";

interface IButtonProps {
  disabled?: boolean;
  size?: string;
}

export const Button = styled.button<IButtonProps>`
  position: relative;
  outline: none;
  cursor: pointer;
  appearance: none;
  text-align: center;
  color: var(--text);
  font-weight: 500;
  min-width: 170px;
  transition: background-color 0.25s ease, color 0.25s ease;
  border: none;
  border-radius: var(--radius-md);
  overflow: hidden;

  // Colors
  background-color: var(--primary);
  color: var(--button-col-primary);
  &:hover {
    background-color: var(--button-hover-primary);
    color: var(--primary);
  }

  ${({ size }) => {
    switch (size) {
      case "xxs": {
        return css`
          padding: 8px;
          font-size: 10px;
          line-height: 0; // 11px;
          min-width: 90px;
        `;
      }
      case "xs": {
        return css`
          padding: 12px 20px;
          font-size: 14px;
          line-height: 0; // 14px;

          @media (min-width: 1680px) {
            font-size: 16px;
            line-height: 0; // 15px;
          }
        `;
      }
      case "sm": {
        return css`
          padding: 16px;
          font-size: 14px;
          line-height: 0; // 14px;
          height: 50px;

          @media (min-width: 1680px) {
            font-size: 16px;
          }
        `;
      }
      case "md": {
        return css`
          padding: 24px;
          font-size: 16px;
          line-height: 0; // 15px;
          height: 50px;

          @media (min-width: 1680px) {
            font-size: 18px;
            line-height: 0; // 18px;
          }
        `;
      }
    }
  }}

  ${({ disabled }) =>
    disabled &&
    css`
      cursor: not-allowed;
    `};
`;
import styled from "styled-components";

export const Input = styled.input`
  outline: none;
  border: none;
  background: var(--primary-light);
  box-sizing: border-box;
  padding: 8px 16px;
  border-radius: 8px;
  font-size: 15px;
  font-weight: 700;
  color: var(--text);
  ::placeholder {
    font-weight: 500;
  }
`;

export default Input;

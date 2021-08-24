import { createGlobalStyle } from 'styled-components';

export const GlobalCSSReset = createGlobalStyle`
  * {
    box-sizing: border-box;
  }
  html,
  body,
  #root {
    width: 100%;
    height: 100%;
    margin: 0;
    padding: 0;
    font-family: sans-serif;
  }
`;

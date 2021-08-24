import React from "react";
import { render } from "react-dom";
import styled from "styled-components";
import { Header } from "./components/layout/Header";
import { GlobalCSSVariables } from "./components/styled/GlobalCSSTheme";
import { GlobalCSSReset } from "./components/styled/GlobalCSSReset";
import { HomeBase } from "./pages/HomeBase";
import { Footer } from "./components/layout/Footer";
import { Toaster } from "react-hot-toast";

const CardanoToolkitUI = () => {
  return (
    <>
      <GlobalCSSReset />
      <GlobalCSSVariables />
      <StyledCardanoToolkitUI>
        <Header />
        <HomeBase />
        <Footer />
        <Toaster />
      </StyledCardanoToolkitUI>
    </>
  );
};

const StyledCardanoToolkitUI = styled.div`
  position: fixed;
  overflow: scroll;
  height: 100%;
  width: 100%;
  background: var(--gradient-sundae);
`;

// Render
render(<CardanoToolkitUI />, document.getElementById("root"));

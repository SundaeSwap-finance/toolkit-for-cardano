import { createGlobalStyle } from "styled-components";

export const GlobalCSSVariables = createGlobalStyle`
  /** FONTS */
  @font-face {
    font-family: "MilliardBook";
    font-style: normal;
    font-weight: 400;
    src: local("Milliard"),
      url("https://fonts.cdnfonts.com/s/28732/Rene Bieder  Milliard Book.woff") format("woff");
  }
  @font-face {
    font-family: "MilliardLight";
    font-style: normal;
    font-weight: 300;
    src: local("Milliard"),
      url("https://fonts.cdnfonts.com/s/28732/Rene Bieder  Milliard Light.woff") format("woff");
  }
  @font-face {
    font-family: "MilliardMedium";
    font-style: normal;
    font-weight: 500;
    src: local("Milliard"),
      url("https://fonts.cdnfonts.com/s/28732/Rene Bieder  Milliard Medium.woff") format("woff");
  }
  @font-face {
    font-family: "MilliardSemiBold";
    font-style: normal;
    font-weight: 600;
    src: local("Milliard"),
      url("https://fonts.cdnfonts.com/s/28732/Rene Bieder  Milliard SemiBold.woff") format("woff");
  }
  @font-face {
    font-family: "MilliardBold";
    font-style: normal;
    font-weight: 700;
    src: local("Milliard"),
      url("https://fonts.cdnfonts.com/s/28732/Rene Bieder  Milliard Bold.woff") format("woff");
  }
  @font-face {
    font-family: "MilliardExtraBold";
    font-style: normal;
    font-weight: 800;
    src: local("Milliard"),
      url("https://fonts.cdnfonts.com/s/28732/Rene Bieder  Milliard ExtraBold.woff") format("woff");
  }
  @font-face {
    font-family: "MilliardHeavy";
    font-style: normal;
    font-weight: 850;
    src: local("Milliard"),
      url("https://fonts.cdnfonts.com/s/28732/Rene Bieder  Milliard Heavy.woff") format("woff");
  }
  @font-face {
    font-family: "MilliardBlack";
    font-style: normal;
    font-weight: 900;
    src: local("Milliard"),
      url("https://fonts.cdnfonts.com/s/28732/Rene Bieder  Milliard Black.woff") format("woff");
  }
  @font-face {
    font-family: "Inter";
    font-style: normal;
    font-weight: 400;
    font-display: swap;
    src: url(https://fonts.gstatic.com/s/inter/v3/UcC73FwrK3iLTeHuS_fvQtMwCp50KnMa1ZL7W0Q5nw.woff2)
      format("woff2");
    unicode-range: U+0000-00FF, U+0131, U+0152-0153, U+02BB-02BC, U+02C6, U+02DA, U+02DC,
      U+2000-206F, U+2074, U+20AC, U+2122, U+2191, U+2193, U+2212, U+2215, U+FEFF, U+FFFD;
  }
  /* latin */
  @font-face {
    font-family: "Inter";
    font-style: normal;
    font-weight: 500;
    font-display: swap;
    src: url(https://fonts.gstatic.com/s/inter/v3/UcC73FwrK3iLTeHuS_fvQtMwCp50KnMa1ZL7W0Q5nw.woff2)
      format("woff2");
    unicode-range: U+0000-00FF, U+0131, U+0152-0153, U+02BB-02BC, U+02C6, U+02DA, U+02DC,
      U+2000-206F, U+2074, U+20AC, U+2122, U+2191, U+2193, U+2212, U+2215, U+FEFF, U+FFFD;
  }

  * {
    font-family: "MilliardMedium"
  }

  body {
    /** RADII */
    --radius-xs: 5px;
    --radius-sm: 10px;
    --radius-md: 15px;
    --radius-lg: 20px;
    --radius-xl: 25px;
    --radius-circle: 50%;

    /** MAIN COLORS */
    --primary: #ff9479;
    --primary-light: #fff3f0;
    --primary-semilight: #ffd1c6;
    --primary-dark: #f58d88;
    --primary-alt: #ffb99f;
    --secondary: #fcf4f2;
    --secondary-light: #fcf2ea;
    --secondary-dark: #fae4d3;
    --disabled: #d8d8d8;
    --inactive: #f4f4f4;

    /** BUTTONS */
    --button-col-primary: #fff;
    --button-col-primary-light: #000;
    --button-col-primary-dark: #fff;
    --button-col-secondary: #f78f8a;
    --button-col-secondary-light: #f78f8a;
    --button-col-secondary-dark: #000;
    --button-col-disabled: #8d8d8d;
    --button-col-inactive: #000;
    --button-col-disconnect: #ff886a;

    /** BUTTON HOVER */
    --button-hover-primary: #fcf2ea;
    --button-hover-primary-light: #ff9479;
    --button-hover-primary-semilight: #ffb99f;
    --button-hover-primary-dark: #dc7e7a;
    --button-hover-secondary: #f58d88;
    --button-hover-color-secondary: #fcf4f2;
    --button-hover-secondary-light: #f58d88;
    --button-hover-secondary-dark: #f7d6bc;
    --button-hover-color-secondary-dark: #191b1f;
    --button-hover-color-inactive: #fff;

    /** BORDERS */
    --border-primary: #ff9479;
    --border-primary-light: #ffd1c6;
    --border-primary-dark: #f78f8a;
    --border-secondary: #fcf4f2;
    --border-secondary-light: #fff;
    --border-secondary-dark: #fae4d3;
    --border-alternative: #ffb376;
    --border-alternative-muted: #faddd5;
    --border-muted: #adadad;
    --border-silent: #d5d5d5;
    --border: #c3c3c3;
    --border-swap-routes: rgba(185, 185, 185, 0.3);

    /** TEXT COLORS */
    --text: #000;
    --text-muted: #7b7878;
    --text-silent: #d8d8d8;
    --text-unobtrusive: #535353;
    --text-primary: #f58d88;
    --text-dark: #f78f8a;
    --text-secondary: #ff9479;
    --text-no-pos: #8b8b8b;
    --text-header-btn: #f78f8a;
    --danger: #d90000;
    --success: #09a827;

    /** MISC COLORS */
    --contrast: #000;
    --complementary: #fff;

    /** ICON COLORS */
    --header-icon: #5a5a5a;
    --swapbox-separator: #6e6e6e;

    /** BACKGROUNDS */
    --main-bg: #fff;
    --swapbox-bg: #f4f4f4;
    --select-bg: #fff;
    --card-bottom-bg: #fcf2ea;
    --header-button-bg: #fff3f0;
    --active-navigation: #fff3f0;
    --connected-wallet-bg: #fff3f0;
    --confirm-summary: #fff3f0;

    /** SHADOWS */
    --shadow-assetSelect: 0px 3px 60px rgba(0, 0, 0, 0.16);
    --shadow-card: 0px 0.5px 20px rgba(0, 0, 0, 0.07);
    --shadow-cardHover: "drop-shadow(0px 0px 20px rgba(0, 226, 255, 0.6)";
    --shadow-cardCTA: 0px 3px 10px rgba(0, 0, 0, 0.16);
    --shadow-connectWallet: 0px 5px 50px rgba(0, 0, 0, 0.1);
    --shadow-headerButtons: 0px 5px 50px rgba(0, 0, 0, 0.06);
    --shadow-liqHeader: 0px 0.2px 60px rgba(0, 0, 0, 0.28);
    --shadow-modal: 0px 3px 30px rgba(0, 0, 0, 0.16);
    --shadow-navigation: 0px 5px 99px rgba(0, 0, 0, 0.09);
    --shadow-pdp: 0px 5px 30px rgba(0, 0, 0, 0.1);
    --shadow-overview: 0px 3px 30px rgba(0, 0, 0, 0.16);
    --shadow-liq-card: 0px 0px 20px rgba(49, 0, 255, 0.3);
    --shadow-liq-card: 0px 0.5px 20px rgba(0, 0, 0, 0.07);

    /** GRADIENTS */
    --gradient-sundae: linear-gradient(to right, rgb(248, 226, 160) 0%, rgb(247, 205, 167) 21.67%, rgb(247, 207, 166) 39.41%, rgb(246, 192, 171) 53.69%, rgb(245, 178, 175) 100%);
    --gradient-silky: linear-gradient(rgba(250, 249, 250, 0) 0%, rgba(250, 249, 250, 0.26) 11.14%, rgba(250, 249, 250, 0.45) 25.87%, rgba(250, 249, 250, 0.58) 31.44%, rgba(250, 249, 250, 0.64) 43.85%, rgba(250, 249, 250, 0.79) 57.84%, rgba(250, 249, 250, 0.94) 71.67%, rgb(250, 249, 250) 100%);
    --gradient-card-header: linear-gradient(to right, rgb(248, 228, 167) 0%, rgb(244, 181, 177) 100%);

    /** BAR GRAPH GRADIENTS */
    --bar-1: #f5b2af;
    --bar-2: #f6c0ab;
    --bar-3: #f7cfa6;
    --bar-4: #f7cda7;
    --bar-5: #f8e2a0;

    /** LINE CHART COLORS */
    --line-primary: rgb(247, 143, 138)
  }
`;
import { default as hotToast } from "react-hot-toast";

const toastConfig = {
  style: {
    border: '1px solid #713200',
    padding: '16px',
    color: '#713200',
    fontSize: '8px',
  }
}

export const toast = (message: string): void => {
  hotToast.success(message, toastConfig);
};

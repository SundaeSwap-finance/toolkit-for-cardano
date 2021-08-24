import { useEffect, useRef } from "react";

type Callback = () => void;

export function useInterval(callback: Callback, delay: number): void {
  const savedCallback = useRef<Callback>();

  useEffect(() => {
    savedCallback.current = callback;
  }, [callback]);

  useEffect(() => {
    function tick(): void {
      if (!savedCallback.current) return;
      savedCallback.current();
    }

    const intervalId = setInterval(tick, delay);

    return () => clearInterval(intervalId);
  }, [delay]);
}

export default useInterval;

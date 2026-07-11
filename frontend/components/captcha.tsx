"use client";

import { createElement, useEffect, useRef, useState } from "react";
import Script from "next/script";
import { CAP_WIDGET_ENDPOINT } from "../lib/cap";

interface CaptchaProps {
  // Çözülünce token, sıfırlanınca/süresi dolunca null döner.
  onVerify: (token: string | null) => void;
}

export default function Captcha({ onVerify }: CaptchaProps) {
  const widgetRef = useRef<HTMLElement | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const el = widgetRef.current;
    if (!el) return;

    const handleSolve = (e: Event) => {
      const token = (e as CustomEvent<{ token?: string }>).detail?.token ?? null;
      onVerify(token);
      setError(null);
    };
    const handleError = (e: Event) => {
      const detail = (e as CustomEvent<{ message?: string }>).detail;
      if (detail?.message) console.warn("[cap-widget] hata:", detail.message);
      onVerify(null);
      setError("Doğrulama başarısız oldu. Lütfen tekrar deneyin.");
    };
    const handleReset = () => onVerify(null);
    const handleExpire = () => {
      onVerify(null);
      setError("Doğrulama süresi doldu. Lütfen tekrar deneyin.");
    };

    el.addEventListener("solve", handleSolve);
    el.addEventListener("error", handleError);
    el.addEventListener("reset", handleReset);
    el.addEventListener("expire", handleExpire);
    return () => {
      el.removeEventListener("solve", handleSolve);
      el.removeEventListener("error", handleError);
      el.removeEventListener("reset", handleReset);
      el.removeEventListener("expire", handleExpire);
    };
  }, [onVerify]);

  return (
    <div className="mb-4">
      <Script
        src="https://cdn.jsdelivr.net/npm/@cap.js/widget@0.1.56"
        strategy="afterInteractive"
      />
      <label className="block text-sm font-medium mb-1 text-left">
        Güvenlik Doğrulaması
      </label>
      <p className="text-left opacity-50 text-xs mb-2">
        Spam&apos;i önlemek için lütfen doğrulamayı tamamlayın.
      </p>
      {/* cap-widget bir web component; JSX IntrinsicElements'e bağımlı olmamak için createElement ile render edilir. */}
      {createElement("cap-widget", {
        ref: widgetRef,
        "data-cap-api-endpoint": CAP_WIDGET_ENDPOINT,
      })}
      {error && (
        <div className="mt-2 p-2 bg-red-100 border border-red-400 text-red-700 rounded text-sm">
          {error}
        </div>
      )}
    </div>
  );
}
